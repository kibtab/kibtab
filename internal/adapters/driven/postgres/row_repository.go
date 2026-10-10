package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kibtab/kibtab/internal/core/domain"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// pageSize is the number of rows that Read returns in one call.
const pageSize = 100

// RowRepository reads and writes rows in one table for PostgreSQL.
type RowRepository struct {
	pool    poolIface
	dialect ports.Dialect
}

// Compile-time proof that RowRepository satisfies the RowRepository port.
var _ ports.RowRepository = (*RowRepository)(nil)

// NewRowRepository builds a RowRepository from a pool and a dialect.
func NewRowRepository(pool poolIface, dialect ports.Dialect) *RowRepository {
	return &RowRepository{pool: pool, dialect: dialect}
}

// Read returns the cell values for one table in key order.
func (r *RowRepository) Read(table string) ([]domain.CellValue, error) {
	ctx := context.Background()
	quoted := r.dialect.Quote(table)
	quotedPK := r.dialect.Quote(pkColumn)

	fields, err := r.tableFields(ctx, table)
	if err != nil {
		return nil, fmt.Errorf("read table %q: %w", table, err)
	}
	if len(fields) == 0 {
		return nil, nil
	}

	columns := make([]string, 0, 1+len(fields))
	columns = append(columns, quotedPK)
	for _, f := range fields {
		columns = append(columns, r.dialect.Quote(f))
	}
	selectCols := strings.Join(columns, ", ")

	query := fmt.Sprintf(
		"SELECT %s FROM %s ORDER BY %s LIMIT $1",
		selectCols, quoted, quotedPK,
	)

	rows, err := r.pool.Query(ctx, query, pageSize)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}
	defer rows.Close()

	var cells []domain.CellValue
	for rows.Next() {
		pkVal := ""
		texts := make([]pgtype.Text, len(fields))
		args := make([]any, 0, 1+len(fields))
		args = append(args, &pkVal)
		for i := range fields {
			args = append(args, &texts[i])
		}
		if err := rows.Scan(args...); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		cells = append(cells, domain.CellValue{
			Field: pkColumn,
			Value: pkVal,
		})
		for i, f := range fields {
			val := ""
			if texts[i].Valid {
				val = texts[i].String
			}
			cells = append(cells, domain.CellValue{
				Field: f,
				Value: val,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return cells, nil
}

// Write writes one row in the table. It is idempotent for the same row key.
func (r *RowRepository) Write(table string, row string, values []domain.CellValue) error {
	ctx := context.Background()
	quoted := r.dialect.Quote(table)
	quotedPK := r.dialect.Quote(pkColumn)

	for _, v := range values {
		if _, err := r.dialect.MapValue(v.Field, v.Value); err != nil {
			return fmt.Errorf("write table %q row %q: %w", table, row, err)
		}
	}

	if len(values) == 0 {
		return nil
	}

	cols := make([]string, 0, 1+len(values))
	cols = append(cols, quotedPK)
	for _, v := range values {
		cols = append(cols, r.dialect.Quote(v.Field))
	}
	colList := strings.Join(cols, ", ")

	params := make([]string, 0, 1+len(values))
	params = append(params, "$1")
	for i := range values {
		params = append(params, fmt.Sprintf("$%d", i+2))
	}
	paramList := strings.Join(params, ", ")

	assignments := make([]string, 0, len(values))
	for i, v := range values {
		assignments = append(assignments, fmt.Sprintf("%s = $%d",
			r.dialect.Quote(v.Field), i+2))
	}
	setClause := strings.Join(assignments, ", ")

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s) DO UPDATE SET %s",
		quoted, colList, paramList, quotedPK, setClause,
	)

	args := make([]any, 0, 1+len(values))
	args = append(args, row)
	for _, v := range values {
		args = append(args, v.Value)
	}

	ct, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("write row: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("write row: no rows affected")
	}
	return nil
}

// tableFields returns the field column names for one table, excluding pkColumn.
func (r *RowRepository) tableFields(ctx context.Context, table string) ([]string, error) {
	query := `SELECT column_name FROM information_schema.columns
WHERE table_name = $1 AND column_name != $2
ORDER BY ordinal_position`
	rows, err := r.pool.Query(ctx, query, table, pkColumn)
	if err != nil {
		return nil, fmt.Errorf("query columns: %w", err)
	}
	defer rows.Close()
	var fields []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return nil, fmt.Errorf("scan column: %w", err)
		}
		fields = append(fields, col)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return fields, nil
}
