package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// poolIface is the minimal subset of pgxpool.Pool that the adapter needs.
// It allows tests to inject a fake pool for error-path testing.
type poolIface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Close()
}

// compile-time proof that *pgxpool.Pool satisfies poolIface.
var _ poolIface = (*pgxpool.Pool)(nil)
