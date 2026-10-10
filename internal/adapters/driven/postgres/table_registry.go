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
func (r *TableRegistry) List() ([]domain.TableMetadata, error) {
	metas := make([]domain.TableMetadata, 0, len(r.tables))
	for _, name := range r.tables {
		meta, err := r.Metadata(name)
		if err != nil {
			return nil, fmt.Errorf("list table %q: %w", name, err)
		}
		metas = append(metas, meta)
	}
	return metas, nil
}

// Metadata returns the metadata for one named table.
func (r *TableRegistry) Metadata(table string) (domain.TableMetadata, error) {
	fields, err := r.columns(context.Background(), table)
	if err != nil {
		return domain.TableMetadata{}, fmt.Errorf("read columns for %q: %w", table, err)
	}
	version, err := r.version(context.Background(), table)
	if err != nil {
		return domain.TableMetadata{}, fmt.Errorf("read version for %q: %w", table, err)
	}
	return domain.TableMetadata{
		Name:    table,
		Fields:  fields,
		Version: version,
	}, nil
}

// columns returns the field column names for one table, excluding pkColumn.
func (r *TableRegistry) columns(ctx context.Context, table string) ([]string, error) {
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

// version returns the current version for one table from
// _kibtab_meta.table_versions.
func (r *TableRegistry) version(ctx context.Context, table string) (int64, error) {
	var version int64
	err := r.pool.QueryRow(ctx,
		"SELECT version FROM _kibtab_meta.table_versions WHERE table_name = $1",
		table,
	).Scan(&version)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("query version: %w", err)
	}
	return version, nil
}
