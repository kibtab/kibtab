package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// EnvKeyTables is the environment variable that holds the comma-separated
// table list.
const EnvKeyTables = "KIBTAB_TABLES"

// EnvKeyMaxConns is the environment variable that holds the pool size.
const EnvKeyMaxConns = "KIBTAB_DB_MAX_CONNS"

// EnvKeyMinConns is the environment variable that holds the pool minimum.
const EnvKeyMinConns = "KIBTAB_DB_MIN_CONNS"

// Engine holds the PostgreSQL connection pool and the table list.
type Engine struct {
	pool       *pgxpool.Pool
	dialect    Dialect
	tables     []string
	registry   *TableRegistry
	repository *RowRepository
	runner     *TransactionRunner
}

// Compile-time proof that Engine methods satisfy the port interfaces.
var (
	_ ports.Dialect           = Dialect{}
	_ ports.TableRegistry     = (*TableRegistry)(nil)
	_ ports.RowRepository     = (*RowRepository)(nil)
	_ ports.TransactionRunner = (*TransactionRunner)(nil)
)

// NewEngine builds the PostgreSQL engine from a database URL and an
// environment getter. It parses the pool settings, reads the table list,
// and runs the migrations.
func NewEngine(ctx context.Context, getenv func(string) string) (*Engine, error) {
	url := getenv("DATABASE_URL")
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	if max := getenv(EnvKeyMaxConns); max != "" {
		maxConns, err := strconv.Atoi(max)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", EnvKeyMaxConns, err)
		}
		config.MaxConns = int32(maxConns)
	}
	if min := getenv(EnvKeyMinConns); min != "" {
		minConns, err := strconv.Atoi(min)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", EnvKeyMinConns, err)
		}
		config.MinConns = int32(minConns)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	tables := ParseTableList(getenv(EnvKeyTables))

	engine := &Engine{
		pool:    pool,
		dialect: Dialect{},
		tables:  tables,
	}
	engine.registry = NewTableRegistry(pool, tables)
	engine.repository = NewRowRepository(pool, engine.dialect)
	engine.runner = NewTransactionRunner(pool)

	return engine, nil
}

// Close releases the connection pool. The caller must not use the engine
// after close.
func (e *Engine) Close() error {
	e.pool.Close()
	return nil
}

// Migrate runs the metadata schema migrations. The caller should call this
// once at start-up, before any other query.
func (e *Engine) Migrate(ctx context.Context) error {
	return Migrate(ctx, e.pool)
}

// Dialect returns the Dialect port for PostgreSQL.
func (e *Engine) Dialect() ports.Dialect {
	return e.dialect
}

// TableRegistry returns the TableRegistry port.
func (e *Engine) TableRegistry() ports.TableRegistry {
	return e.registry
}

// RowRepository returns the RowRepository port.
func (e *Engine) RowRepository() ports.RowRepository {
	return e.repository
}

// TransactionRunner returns the TransactionRunner port.
func (e *Engine) TransactionRunner() ports.TransactionRunner {
	return e.runner
}

// Pool returns the connection pool. It is used by tests that need direct
// access to the database.
func (e *Engine) Pool() *pgxpool.Pool {
	return e.pool
}

// Tables returns the parsed table list.
func (e *Engine) Tables() []string {
	return append([]string(nil), e.tables...)
}
