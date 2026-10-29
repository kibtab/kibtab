package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// execer runs one statement on a connection. Both the pool and a
// transaction satisfy it, so a port call runs on the transaction when the
// context holds one and on the pool when it does not.
type execer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// poolIface is the minimal subset of pgxpool.Pool that the adapter needs.
// It allows tests to inject a fake pool for error-path testing.
type poolIface interface {
	execer
	Ping(ctx context.Context) error
	Close()
}

// compile-time proof that *pgxpool.Pool satisfies poolIface.
var _ poolIface = (*pgxpool.Pool)(nil)

// compile-time proof that pgx.Tx satisfies execer.
var _ execer = (pgx.Tx)(nil)
