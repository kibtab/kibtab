package postgres

import (
	"context"
	"testing"
)

func TestLockVersionStoreError(t *testing.T) {
	pool := &fakePool{execErr: errDummy}
	reg := NewTableRegistry(pool, []string{"users"})

	_, err := reg.LockVersion(context.Background(), "users")
	if err == nil {
		t.Fatal("LockVersion with store error: got nil, want error")
	}
}

func TestLockVersionLockError(t *testing.T) {
	pool := &fakePool{row: &fakeRow{err: errDummy}}
	reg := NewTableRegistry(pool, []string{"users"})

	_, err := reg.LockVersion(context.Background(), "users")
	if err == nil {
		t.Fatal("LockVersion with lock error: got nil, want error")
	}
}

func TestBumpVersionError(t *testing.T) {
	pool := &fakePool{row: &fakeRow{err: errDummy}}
	reg := NewTableRegistry(pool, []string{"users"})

	_, err := reg.BumpVersion(context.Background(), "users")
	if err == nil {
		t.Fatal("BumpVersion with query error: got nil, want error")
	}
}

// TestIntegrationVersionLockAndBump stores the first version of a fresh
// table, it bumps the version, and it reads the new version back.
func TestIntegrationVersionLockAndBump(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)
	if err := Migrate(context.Background(), pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	createTestTable(t, pool, "items", []string{"name"})
	reg := NewTableRegistry(pool, []string{"items"})

	locked, err := reg.LockVersion(context.Background(), "items")
	if err != nil {
		t.Fatalf("LockVersion: %v", err)
	}
	if locked != firstVersion {
		t.Fatalf("lock on a fresh table: got %d, want %d", locked,
			firstVersion)
	}

	bumped, err := reg.BumpVersion(context.Background(), "items")
	if err != nil {
		t.Fatalf("BumpVersion: %v", err)
	}
	if bumped != firstVersion+1 {
		t.Fatalf("bump: got %d, want %d", bumped, firstVersion+1)
	}

	meta, err := reg.Metadata(context.Background(), "items")
	if err != nil {
		t.Fatalf("Metadata: %v", err)
	}
	if meta.Version != firstVersion+1 {
		t.Fatalf("stored version: got %d, want %d", meta.Version,
			firstVersion+1)
	}
}

// TestIntegrationBumpVersionWithoutRow fails when the table has no
// version row.
func TestIntegrationBumpVersionWithoutRow(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)
	if err := Migrate(context.Background(), pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	reg := NewTableRegistry(pool, []string{"items"})
	if _, err := reg.BumpVersion(context.Background(), "items"); err == nil {
		t.Fatal("BumpVersion without a version row: got nil, want error")
	}
}
