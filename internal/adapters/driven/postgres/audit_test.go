package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/kibtab/kibtab/internal/core/domain"
)

func TestAuditRecordExecError(t *testing.T) {
	pool := &fakePool{execErr: errDummy}
	audit := NewAuditRepository(pool, Dialect{})

	err := audit.Record(context.Background(), domain.AuditRow{
		Table: "users", Row: "1", Field: "name",
		User: "rita", OldValue: "bob", NewValue: "alice",
	})
	if err == nil {
		t.Fatal("Record with exec error: got nil, want error")
	}
}

func TestAuditReadQueryError(t *testing.T) {
	pool := &fakePool{queryErr: errDummy}
	audit := NewAuditRepository(pool, Dialect{})

	_, err := audit.Read(context.Background(), "users", "1", "name")
	if err == nil {
		t.Fatal("Read with query error: got nil, want error")
	}
}

func TestAuditReadScanError(t *testing.T) {
	rows := &fakeRows{next: true, scanErr: errDummy}
	pool := &fakePool{rows: rows}
	audit := NewAuditRepository(pool, Dialect{})

	_, err := audit.Read(context.Background(), "users", "1", "name")
	if err == nil {
		t.Fatal("Read with scan error: got nil, want error")
	}
}

func TestAuditReadRowsErr(t *testing.T) {
	rows := &fakeRows{next: false, err: errDummy}
	pool := &fakePool{rows: rows}
	audit := NewAuditRepository(pool, Dialect{})

	_, err := audit.Read(context.Background(), "users", "1", "name")
	if err == nil {
		t.Fatal("Read with rows.Err: got nil, want error")
	}
}

// TestIntegrationAuditRoundTrip writes one audit row and reads it back.
func TestIntegrationAuditRoundTrip(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)
	if err := Migrate(context.Background(), pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	audit := NewAuditRepository(pool, Dialect{})
	entry := domain.AuditRow{
		Table:    "items",
		Row:      "1",
		Field:    "name",
		User:     "rita",
		OldValue: "bob",
		NewValue: "alice",
	}
	if err := audit.Record(context.Background(), entry); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := audit.Read(context.Background(), "items", "1", "name")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("audit rows: got %d, want 1", len(entries))
	}
	got := entries[0]
	if got.User != "rita" || got.OldValue != "bob" || got.NewValue != "alice" {
		t.Fatalf("audit row: got %+v, want the change of rita", got)
	}
	if got.Table != "items" || got.Row != "1" || got.Field != "name" {
		t.Fatalf("audit cell: got %+v, want items/1/name", got)
	}
	if got.ChangedAt.IsZero() {
		t.Fatal("audit time: got the zero time, want the time of the write")
	}
	if time.Since(got.ChangedAt) > time.Minute {
		t.Fatalf("audit time: got %v, want a recent time", got.ChangedAt)
	}
}
