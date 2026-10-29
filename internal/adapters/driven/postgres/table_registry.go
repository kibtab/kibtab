package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kibtab/kibtab/internal/core/domain"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// pkColumn is the name of the primary key column that Kibtab uses in each
// table it manages. The spreadsheet row header maps to this column.
const pkColumn = "k"

// firstVersion is the version of a table that has no version row. The
// registry stores it at the first read of the version.
const firstVersion = 1

// TableRegistry reads the table list from the environment.
type TableRegistry struct {
	pool   poolIface
	tables []string
}

// Compile-time proof that TableRegistry satisfies the TableRegistry port.
var _ ports.TableRegistry = (*TableRegistry)(nil)

// NewTableRegistry builds a TableRegistry from a pool and the table list.
func NewTableRegistry(pool poolIface, tables []string) *TableRegistry {
	return &TableRegistry{
		pool:   pool,
		tables: tables,
	}
}

// ParseTableList splits a comma-separated table list into names.
// It trims spaces and drops empty entries.
func ParseTableList(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		name := strings.TrimSpace(p)
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}

// List returns the metadata for every known table.
func (r *TableRegistry) List(ctx context.Context) ([]domain.TableMetadata, error) {
	metas := make([]domain.TableMetadata, 0, len(r.tables))
	for _, name := range r.tables {
		meta, err := r.Metadata(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("list table %q: %w", name, err)
		}
		metas = append(metas, meta)
	}
	return metas, nil
}

// Metadata returns the metadata for one named table.
func (r *TableRegistry) Metadata(ctx context.Context,
	table string) (domain.TableMetadata, error) {
	fields, err := r.columns(ctx, table)
	if err != nil {
		return domain.TableMetadata{}, fmt.Errorf("read columns for %q: %w", table, err)
	}
	version, err := r.version(ctx, table)
	if err != nil {
		return domain.TableMetadata{}, fmt.Errorf("read version for %q: %w", table, err)
	}
	return domain.TableMetadata{
		Name:    table,
		Fields:  fields,
		Version: version,
	}, nil
}

// LockVersion returns the stored version of one table. It stores the
// first version when the table has no version row. It holds the version
// row until the transaction of the caller ends, so a second write waits
// and it reads the new version after the first write commits.
func (r *TableRegistry) LockVersion(ctx context.Context, table string) (int64, error) {
	ex := execIn(ctx, r.pool)
	const ensure = `INSERT INTO _kibtab_meta.table_versions (table_name, version)
VALUES ($1, $2) ON CONFLICT (table_name) DO NOTHING`
	if _, err := ex.Exec(ctx, ensure, table, firstVersion); err != nil {
		return 0, fmt.Errorf("store version row for %q: %w", table, err)
	}
	var version int64
	const lock = `SELECT version FROM _kibtab_meta.table_versions
WHERE table_name = $1 FOR UPDATE`
	if err := ex.QueryRow(ctx, lock, table).Scan(&version); err != nil {
		return 0, fmt.Errorf("lock version for %q: %w", table, err)
	}
	return version, nil
}

// BumpVersion advances the stored version of one table by one. It
// returns the new version.
func (r *TableRegistry) BumpVersion(ctx context.Context, table string) (int64, error) {
	ex := execIn(ctx, r.pool)
	var version int64
	const bump = `UPDATE _kibtab_meta.table_versions SET version = version + 1
WHERE table_name = $1 RETURNING version`
	if err := ex.QueryRow(ctx, bump, table).Scan(&version); err != nil {
		return 0, fmt.Errorf("bump version for %q: %w", table, err)
	}
	return version, nil
}

// columns returns the field column names for one table, excluding pkColumn.
func (r *TableRegistry) columns(ctx context.Context, table string) ([]string, error) {
	query := `SELECT column_name FROM information_schema.columns
WHERE table_name = $1 AND column_name != $2
ORDER BY ordinal_position`
	rows, err := execIn(ctx, r.pool).Query(ctx, query, table, pkColumn)
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

// version returns the current version for one table from
// _kibtab_meta.table_versions. A table with no version row sits at the
// first version, so a client always reads a version that it can write.
func (r *TableRegistry) version(ctx context.Context, table string) (int64, error) {
	ex := execIn(ctx, r.pool)
	var version int64
	err := ex.QueryRow(ctx,
		"SELECT version FROM _kibtab_meta.table_versions WHERE table_name = $1",
		table,
	).Scan(&version)
	if err == pgx.ErrNoRows {
		return firstVersion, nil
	}
	if err != nil {
		return 0, fmt.Errorf("query version: %w", err)
	}
	return version, nil
}
