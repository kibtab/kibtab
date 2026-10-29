package postgres

import (
	"context"
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

func TestExecInUsesThePoolWithoutATransaction(t *testing.T) {
	pool := &fakePool{}

	if got := execIn(context.Background(), pool); got != execer(pool) {
		t.Fatalf("execIn without a transaction: got %T, want the pool", got)
	}
}

func TestExecInUsesTheTransaction(t *testing.T) {
	pool := &fakePool{}
	tx := &fakeTx{}

	ctx := withTx(context.Background(), tx)
	if got := execIn(ctx, pool); got != execer(tx) {
		t.Fatalf("execIn with a transaction: got %T, want the transaction", got)
	}
}

func TestWriteRunsOnTheTransaction(t *testing.T) {
	pool := &fakePool{}
	tx := &fakeTx{}
	repo := NewRowRepository(pool, Dialect{})
	ctx := withTx(context.Background(), tx)

	err := repo.Write(ctx, "users", "r1", []domain.CellValue{
		{Field: "name", Value: "alice"},
	})
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if tx.execCalls != 1 {
		t.Fatalf("transaction writes: got %d, want 1", tx.execCalls)
	}
	if pool.execCalls != 0 {
		t.Fatalf("pool writes: got %d, want 0", pool.execCalls)
	}
}
