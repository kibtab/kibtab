package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// txKey is the context key that holds the transaction of one write group.
type txKey struct{}

// withTx returns a context that holds tx. The TransactionRunner hands
// that context to the sync service, so every port call of one write
// group runs on the same transaction.
func withTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// execIn returns the connection that the context holds. It returns the
// pool when the context holds no transaction.
func execIn(ctx context.Context, pool poolIface) execer {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}
