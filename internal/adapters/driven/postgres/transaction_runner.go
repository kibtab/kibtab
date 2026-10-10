package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// txBeginner is the minimal interface for starting a transaction.
type txBeginner interface {
	BeginTx(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error)
}

// TransactionRunner opens a transaction and runs a function inside it.
type TransactionRunner struct {
	pool txBeginner
}

// Compile-time proof that TransactionRunner satisfies the TransactionRunner port.
var _ ports.TransactionRunner = (*TransactionRunner)(nil)

// NewTransactionRunner builds a TransactionRunner from a pool.
func NewTransactionRunner(pool txBeginner) *TransactionRunner {
	return &TransactionRunner{pool: pool}
}

// Run runs fn inside one transaction. It commits on success and rolls back
// on error.
func (t *TransactionRunner) Run(fn func() error) error {
	ctx := context.Background()
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}
	tx, err := t.pool.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := fn(); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback: %v (after error: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
