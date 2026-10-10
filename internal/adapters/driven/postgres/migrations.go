package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// migrationKey is the advisory lock key that guards the migration step.
const migrationKey = 0x6b6962746162

// migrationStmts holds the SQL statements that run at start.
const migrationStmts = `
CREATE SCHEMA IF NOT EXISTS _kibtab_meta;

CREATE TABLE IF NOT EXISTS _kibtab_meta.table_versions (
    table_name text PRIMARY KEY,
    version bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS _kibtab_meta.audit_log (
    id bigserial PRIMARY KEY,
    table_name text NOT NULL,
    row_key text NOT NULL,
    field_name text NOT NULL,
    old_value text,
    new_value text,
    changed_at timestamptz NOT NULL DEFAULT now()
);
`

// splitStatements splits the migration SQL into individual statements.
func splitStatements(sql string) []string {
	parts := strings.Split(sql, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// migrationExecer is the minimal interface for running migration SQL on a
// connection. Both *pgx.Conn and test doubles satisfy it.
type migrationExecer interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// runMigrations runs the lock, migration, and unlock steps on conn.
func runMigrations(ctx context.Context, conn migrationExecer) error {
	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", migrationKey); err != nil {
		return fmt.Errorf("acquire advisory lock: %w", err)
	}

	for _, stmt := range splitStatements(migrationStmts) {
		if _, err := conn.Exec(ctx, stmt); err != nil {
			_, _ = conn.Exec(context.Background(),
				"SELECT pg_advisory_unlock($1)", migrationKey)
			return fmt.Errorf("run migration: %w", err)
		}
	}

	if _, err := conn.Exec(context.Background(),
		"SELECT pg_advisory_unlock($1)", migrationKey); err != nil {
		return fmt.Errorf("release advisory lock: %w", err)
	}
	return nil
}

// Migrate runs the migrations under an advisory lock.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()
	return runMigrations(ctx, conn.Conn())
}
