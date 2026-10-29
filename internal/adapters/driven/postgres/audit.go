package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kibtab/kibtab/internal/core/domain"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// AuditRepository reads and writes the audit rows of one cell.
type AuditRepository struct {
	pool    poolIface
	dialect ports.Dialect
}

// Compile-time proof that AuditRepository satisfies both audit ports.
var (
	_ ports.AuditWriter = (*AuditRepository)(nil)
	_ ports.AuditReader = (*AuditRepository)(nil)
)

// NewAuditRepository builds an AuditRepository from a pool and a dialect.
func NewAuditRepository(pool poolIface, dialect ports.Dialect) *AuditRepository {
	return &AuditRepository{pool: pool, dialect: dialect}
}

// Record stores one audit row. The engine stamps the time of the change,
// so each row of one transaction holds the time of the commit.
func (a *AuditRepository) Record(ctx context.Context,
	entry domain.AuditRow) error {
	const query = `INSERT INTO _kibtab_meta.audit_log
(table_name, row_key, field_name, changed_by, old_value, new_value, changed_at)
VALUES ($1, $2, $3, $4, $5, $6, now())`
	ex := execIn(ctx, a.pool)
	_, err := ex.Exec(ctx, query, entry.Table, entry.Row, entry.Field,
		entry.User, entry.OldValue, entry.NewValue)
	if err != nil {
		return fmt.Errorf("write audit row for %q %q %q: %w",
			entry.Table, entry.Row, entry.Field, err)
	}
	return nil
}

// Read returns the audit rows for one cell in write order. The page of
// the query keeps the answer at a fixed size.
func (a *AuditRepository) Read(ctx context.Context, table string, row string,
	field string) ([]domain.AuditRow, error) {
	const body = `SELECT changed_by, old_value, new_value, changed_at
FROM _kibtab_meta.audit_log
WHERE table_name = $1 AND row_key = $2 AND field_name = $3
ORDER BY id ASC`
	clause, limitArgs := a.dialect.Page(pageSize, 4)
	query := body + clause
	args := []any{table, row, field}
	args = append(args, limitArgs...)

	rows, err := execIn(ctx, a.pool).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query audit rows for %q %q: %w",
			table, row, err)
	}
	defer rows.Close()

	var entries []domain.AuditRow
	for rows.Next() {
		var user, oldValue, newValue pgtype.Text
		var changedAt time.Time
		if err := rows.Scan(&user, &oldValue, &newValue, &changedAt); err != nil {
			return nil, fmt.Errorf("scan audit row: %w", err)
		}
		entries = append(entries, domain.AuditRow{
			Table:     table,
			Row:       row,
			Field:     field,
			User:      user.String,
			OldValue:  oldValue.String,
			NewValue:  newValue.String,
			ChangedAt: changedAt,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("audit rows error: %w", err)
	}
	return entries, nil
}
