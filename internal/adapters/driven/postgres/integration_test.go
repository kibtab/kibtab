package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kibtab/kibtab/internal/adapters/driven/contracttest"
	"github.com/kibtab/kibtab/internal/core/domain"
)

// TestIntegration runs the contract suite against a real PostgreSQL instance.
// It skips when no database is reachable.
func TestIntegration(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	tableName := "users"
	fields := []string{"name", "email"}
	createTestTable(t, pool, tableName, fields)

	dialect := Dialect{}
	repo := NewRowRepository(pool, dialect)

	fixture := contracttest.Fixture{
		Repo:  repo,
		Table: tableName,
		Key:   pkColumn,
		Reset: func() { truncateTable(t, pool, tableName) },
	}
	for _, err := range contracttest.RowRepository(fixture) {
		t.Error(err)
	}
}

// TestIntegrationTableRegistry tests the TableRegistry against a real database.
func TestIntegrationTableRegistry(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	createTestTable(t, pool, "items", []string{"title", "price"})
	if _, err := pool.Exec(ctx,
		"INSERT INTO _kibtab_meta.table_versions (table_name, version) VALUES ('items', 5)"); err != nil {
		t.Fatalf("insert version: %v", err)
	}

	reg := NewTableRegistry(pool, []string{"items"})

	t.Run("list returns metadata", func(t *testing.T) {
		metas, err := reg.List()
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(metas) != 1 {
			t.Fatalf("List: got %d tables, want 1", len(metas))
		}
		if metas[0].Name != "items" {
			t.Fatalf("Name: got %q, want %q", metas[0].Name, "items")
		}
		if metas[0].Version != 5 {
			t.Fatalf("Version: got %d, want 5", metas[0].Version)
		}
		if len(metas[0].Fields) != 2 {
			t.Fatalf("Fields: got %d, want 2", len(metas[0].Fields))
		}
	})

	t.Run("metadata returns zero version when not set", func(t *testing.T) {
		createTestTable(t, pool, "empty_version", []string{"f1"})
		reg2 := NewTableRegistry(pool, []string{"empty_version"})
		meta, err := reg2.Metadata("empty_version")
		if err != nil {
			t.Fatalf("Metadata: %v", err)
		}
		if meta.Version != 0 {
			t.Fatalf("Version: got %d, want 0", meta.Version)
		}
	})
}

// TestIntegrationMigrationsIdempotent verifies that Migrate is idempotent.
func TestIntegrationMigrationsIdempotent(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("first Migrate: %v", err)
	}
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("second Migrate: %v", err)
	}
}

// TestIntegrationTransactionRunner tests the TransactionRunner.
func TestIntegrationTransactionRunner(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	createTestTable(t, pool, "txn_test", []string{"name"})

	runner := NewTransactionRunner(pool)

	t.Run("commit on success", func(t *testing.T) {
		err := runner.Run(func() error {
			return nil
		})
		if err != nil {
			t.Fatalf("Run success: %v", err)
		}
	})

	t.Run("rollback on error", func(t *testing.T) {
		err := runner.Run(func() error {
			return fmt.Errorf("simulated failure")
		})
		if err == nil {
			t.Fatal("Run with failing fn: got nil, want error")
		}
	})

	t.Run("engine wires all ports", func(t *testing.T) {
		truncateTable(t, pool, "txn_test")
		if _, err := pool.Exec(ctx,
			"INSERT INTO _kibtab_meta.table_versions (table_name, version) VALUES ('txn_test', 3)"); err != nil {
			t.Fatalf("insert version: %v", err)
		}

		engine := &Engine{
			pool:    pool,
			dialect: Dialect{},
			tables:  []string{"txn_test"},
		}
		engine.registry = NewTableRegistry(pool, []string{"txn_test"})
		engine.repository = NewRowRepository(pool, Dialect{})
		engine.runner = NewTransactionRunner(pool)

		reg := engine.TableRegistry()
		meta, err := reg.Metadata("txn_test")
		if err != nil {
			t.Fatalf("Metadata: %v", err)
		}
		if meta.Name != "txn_test" {
			t.Fatalf("Name: got %q, want %q", meta.Name, "txn_test")
		}
		if meta.Version != 3 {
			t.Fatalf("Version: got %d, want 3", meta.Version)
		}

		if err := engine.RowRepository().Write("txn_test", "e1",
			[]domain.CellValue{{Field: "name", Value: "test"}}); err != nil {
			t.Fatalf("Write: %v", err)
		}
	})
}

// available returns true when a PostgreSQL database is reachable.
func available(t *testing.T) bool {
	t.Helper()
	ctx := context.Background()
	url := testEnvURL(t, "postgres")
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return false
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return false
	}
	defer pool.Close()
	return pool.Ping(ctx) == nil
}
