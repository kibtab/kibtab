package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
	"github.com/kibtab/kibtab/internal/core/services"
)

// rowError returns the first row error of a sync result.
func rowError(result domain.SyncResult) error {
	for _, row := range result.RowResults {
		if row.Error != nil {
			return row.Error
		}
	}
	return nil
}

// TestIntegrationSyncWritesTheAuditRow writes one cell, it finds the
// audit row, and it rejects a write with a stale version.
func TestIntegrationSyncWritesTheAuditRow(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	createTestTable(t, pool, "items", []string{"name"})

	dialect := Dialect{}
	registry := NewTableRegistry(pool, []string{"items"})
	repository := NewRowRepository(pool, dialect)
	audit := NewAuditRepository(pool, dialect)
	runner := NewTransactionRunner(pool)
	service := services.NewSyncService(registry, repository, audit, dialect,
		runner)

	first := domain.SyncPayload{
		Table: "items",
		User:  "rita",
		Delta: []domain.CellDelta{{
			Row: "1", Field: "name", Value: "alice", Version: firstVersion,
		}},
	}
	result, err := service.Sync(first)
	if err != nil {
		t.Fatalf("first Sync: %v", err)
	}
	if err := rowError(result); err != nil {
		t.Fatalf("first Sync row error: %v", err)
	}
	if got := result.RowResults[0].Version; got != firstVersion+1 {
		t.Fatalf("first Sync version: got %d, want %d", got,
			firstVersion+1)
	}

	entries, err := audit.Read(ctx, "items", "1", "name")
	if err != nil {
		t.Fatalf("Read the audit rows: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("audit rows after the first write: got %d, want 1",
			len(entries))
	}
	if entries[0].User != "rita" || entries[0].OldValue != "" ||
		entries[0].NewValue != "alice" {
		t.Fatalf("first audit row: got %+v, want rita alice", entries[0])
	}

	second := domain.SyncPayload{
		Table: "items",
		User:  "rita",
		Delta: []domain.CellDelta{{
			Row: "1", Field: "name", Value: "bob",
			Version: firstVersion + 1,
		}},
	}
	result, err = service.Sync(second)
	if err != nil {
		t.Fatalf("second Sync: %v", err)
	}
	if err := rowError(result); err != nil {
		t.Fatalf("second Sync row error: %v", err)
	}

	entries, err = audit.Read(ctx, "items", "1", "name")
	if err != nil {
		t.Fatalf("Read the audit rows: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("audit rows after the second write: got %d, want 2",
			len(entries))
	}
	if entries[1].OldValue != "alice" || entries[1].NewValue != "bob" {
		t.Fatalf("second audit row: got %+v, want the change to bob",
			entries[1])
	}

	stale := domain.SyncPayload{
		Table: "items",
		User:  "rita",
		Delta: []domain.CellDelta{{
			Row: "1", Field: "name", Value: "carol", Version: firstVersion,
		}},
	}
	result, err = service.Sync(stale)
	if err != nil {
		t.Fatalf("stale Sync: %v", err)
	}
	var conflict domain.ConflictError
	if !errors.As(rowError(result), &conflict) {
		t.Fatalf("stale Sync: got %v, want a conflict", rowError(result))
	}

	cells, err := repository.ReadRow(ctx, "items", "1")
	if err != nil {
		t.Fatalf("ReadRow: %v", err)
	}
	if len(cells) != 1 || cells[0].Value != "bob" {
		t.Fatalf("stored row after the stale write: got %+v, want bob",
			cells)
	}
	entries, err = audit.Read(ctx, "items", "1", "name")
	if err != nil {
		t.Fatalf("Read the audit rows: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("audit rows after the stale write: got %d, want 2",
			len(entries))
	}
}
