package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kibtab/kibtab/internal/core/domain"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// TestEngineAccessors tests the Engine's port accessor methods.
func TestEngineAccessors(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	engine := &Engine{
		pool:    pool,
		dialect: Dialect{},
		tables:  []string{"users"},
	}
	engine.registry = NewTableRegistry(pool, []string{"users"})
	engine.repository = NewRowRepository(pool, Dialect{})
	engine.runner = NewTransactionRunner(pool)

	t.Run("Dialect returns a Dialect", func(t *testing.T) {
		d := engine.Dialect()
		if d.Quote("test") != `"test"` {
			t.Fatalf("Dialect.Quote: got %q, want %q", d.Quote("test"), `"test"`)
		}
	})

	t.Run("TableRegistry returns a working registry", func(t *testing.T) {
		reg := engine.TableRegistry()
		if reg == nil {
			t.Fatal("TableRegistry returned nil")
		}
	})

	t.Run("RowRepository returns a working repository", func(t *testing.T) {
		repo := engine.RowRepository()
		if repo == nil {
			t.Fatal("RowRepository returned nil")
		}
	})

	t.Run("TransactionRunner returns a working runner", func(t *testing.T) {
		tr := engine.TransactionRunner()
		if tr == nil {
			t.Fatal("TransactionRunner returned nil")
		}
	})

	t.Run("Pool returns the connection pool", func(t *testing.T) {
		p := engine.Pool()
		if p == nil {
			t.Fatal("Pool returned nil")
		}
	})

	t.Run("Tables returns a copy", func(t *testing.T) {
		tables := engine.Tables()
		if len(tables) != 1 || tables[0] != "users" {
			t.Fatalf("Tables: got %v, want [users]", tables)
		}
		tables[0] = "hacked"
		if engine.tables[0] != "users" {
			t.Fatal("Tables returned a shared slice")
		}
	})
}

// TestNewEngineValidURL tests creating an engine with a valid DATABASE_URL.
func TestNewEngineValidURL(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	ctx := context.Background()
	getenv := func(key string) string {
		if key == "DATABASE_URL" {
			return testEnvURL(t, "postgres")
		}
		return ""
	}

	engine, err := NewEngine(ctx, getenv)
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	defer engine.Close()

	if err := engine.Migrate(ctx); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
}

// TestNewEngineEmptyURL tests that NewEngine fails when DATABASE_URL is empty.
func TestNewEngineEmptyURL(t *testing.T) {
	ctx := context.Background()
	getenv := func(key string) string { return "" }

	_, err := NewEngine(ctx, getenv)
	if err == nil {
		t.Fatal("NewEngine with empty DATABASE_URL: got nil, want error")
	}
}

// TestNewEngineInvalidURL tests that NewEngine fails on a bad URL.
func TestNewEngineInvalidURL(t *testing.T) {
	ctx := context.Background()
	getenv := func(key string) string {
		if key == "DATABASE_URL" {
			return "not-a-valid-url"
		}
		return ""
	}

	_, err := NewEngine(ctx, getenv)
	if err == nil {
		t.Fatal("NewEngine with invalid URL: got nil, want error")
	}
}

// TestNewEngineNegativeMaxConns tests that NewEngine rejects a negative pool size
// which causes pgxpool.NewWithConfig to fail.
func TestNewEngineNegativeMaxConns(t *testing.T) {
	ctx := context.Background()
	getenv := func(key string) string {
		if key == "DATABASE_URL" {
			return "postgres://localhost/kibtab"
		}
		if key == EnvKeyMaxConns {
			return "-1"
		}
		return ""
	}

	_, err := NewEngine(ctx, getenv)
	if err == nil {
		t.Fatal("NewEngine with negative max conns: got nil, want error")
	}
}

// TestNewEngineInvalidMaxConns tests that NewEngine rejects a bad pool size.
func TestNewEngineInvalidMaxConns(t *testing.T) {
	ctx := context.Background()
	getenv := func(key string) string {
		if key == "DATABASE_URL" {
			return "postgres://localhost/kibtab"
		}
		if key == EnvKeyMaxConns {
			return "not-a-number"
		}
		return ""
	}

	_, err := NewEngine(ctx, getenv)
	if err == nil {
		t.Fatal("NewEngine with bad max conns: got nil, want error")
	}
}

// TestNewEngineInvalidMinConns tests that NewEngine rejects a bad pool minimum.
func TestNewEngineInvalidMinConns(t *testing.T) {
	ctx := context.Background()
	getenv := func(key string) string {
		if key == "DATABASE_URL" {
			return "postgres://localhost/kibtab"
		}
		if key == EnvKeyMinConns {
			return "not-a-number"
		}
		return ""
	}

	_, err := NewEngine(ctx, getenv)
	if err == nil {
		t.Fatal("NewEngine with bad min conns: got nil, want error")
	}
}

// TestNewEngineSetsPoolConfig tests that env vars are applied to the pool.
func TestNewEngineSetsPoolConfig(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	ctx := context.Background()
	getenv := func(key string) string {
		if key == "DATABASE_URL" {
			return testEnvURL(t, "postgres")
		}
		if key == EnvKeyMaxConns {
			return "10"
		}
		if key == EnvKeyMinConns {
			return "2"
		}
		return ""
	}

	engine, err := NewEngine(ctx, getenv)
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	defer engine.Close()

	stats := engine.Pool().Stat()
	if stats.MaxConns() != 10 {
		t.Fatalf("MaxConns: got %d, want 10", stats.MaxConns())
	}
}

// TestDialectImplementsPort is a compile-time check that Dialect satisfies
// the ports.Dialect interface.
func TestDialectImplementsPort(t *testing.T) {
	var _ ports.Dialect = Dialect{}
}

// TestRowRepositoryClosedPoolError tests that Read on a closed pool returns an error.
func TestRowRepositoryClosedPoolError(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	createTestTable(t, pool, "test_table", []string{"name"})

	repo := NewRowRepository(pool, Dialect{})
	pool.Close()

	cells, err := repo.Read("test_table")
	if err == nil {
		t.Fatalf("Read on closed pool: got %v, want error", cells)
	}
}

// TestTableRegistryClosedPoolError tests that Metadata on a closed pool errors.
func TestTableRegistryClosedPoolError(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	reg := NewTableRegistry(pool, []string{"items"})
	pool.Close()

	_, err := reg.Metadata("items")
	if err == nil {
		t.Fatal("Metadata on closed pool: got nil, want error")
	}
}

// TestTableRegistryListError tests that List propagates errors from Metadata.
func TestTableRegistryListError(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	reg := NewTableRegistry(pool, []string{"items"})
	pool.Close()

	_, err := reg.List()
	if err == nil {
		t.Fatal("List on closed pool: got nil, want error")
	}
}

// TestTableRegistryVersionError tests that Metadata returns an error when
// the version query fails due to missing migrations.
func TestTableRegistryVersionError(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	createTestTable(t, pool, "items", []string{"title", "price"})

	reg := NewTableRegistry(pool, []string{"items"})
	_, err := reg.Metadata("items")
	if err == nil {
		t.Fatal("Metadata without migrations: got nil, want error")
	}
}

// TestRowRepositoryWriteToNonExistentTable tests that Write to a missing
// table returns an error.
func TestRowRepositoryWriteToNonExistentTable(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	repo := NewRowRepository(pool, Dialect{})

	err := repo.Write("nonexistent_table_12345", "r1",
		[]domain.CellValue{{Field: "name", Value: "test"}})
	if err == nil {
		t.Fatal("Write to non-existent table: got nil, want error")
	}
}

// TestTransactionRunnerBeginError tests the Run method when the transaction
// cannot begin because the pool is closed.
func TestTransactionRunnerBeginError(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	runner := NewTransactionRunner(pool)
	pool.Close()

	err := runner.Run(func() error {
		return nil
	})
	if err == nil {
		t.Fatal("Run on closed pool: got nil, want error")
	}
}

// TestMigrateAcquireError tests that Migrate returns an error when it cannot
// acquire a connection.
func TestMigrateAcquireError(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("first Migrate: %v", err)
	}

	pool.Close()

	pool2, err := pgxpool.NewWithConfig(ctx, pool.Config().Copy())
	if err != nil {
		t.Fatalf("create second pool: %v", err)
	}
	pool2.Close()

	err = Migrate(ctx, pool2)
	if err == nil {
		t.Fatal("Migrate on closed pool: got nil, want error")
	}
}
