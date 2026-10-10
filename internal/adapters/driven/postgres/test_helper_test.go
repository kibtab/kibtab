package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// defaultTestDSN is the connection string for a developer run. It uses the
// local socket of a development PostgreSQL.
const defaultTestDSN = "host=/tmp/pgsocket port=5433 user=postgres sslmode=disable"

// testEnvURL returns the connection string for the named database.
// It reads KIBTAB_TEST_DATABASE_URL for the host and the credentials.
// It falls back to defaultTestDSN for a developer run.
func testEnvURL(t *testing.T, dbName string) string {
	t.Helper()
	raw := os.Getenv(EnvKeyTestDatabaseURL)
	if raw == "" {
		raw = defaultTestDSN
	}
	config, err := pgxpool.ParseConfig(raw)
	if err != nil {
		t.Fatalf("parse %s: %v", EnvKeyTestDatabaseURL, err)
	}
	// Each test needs its own database. The name must replace the name in
	// the environment value, not sit beside it. pgx keeps the original
	// string in ConnString, so the helper writes a new one.
	config.ConnConfig.Database = dbName
	return buildDSN(&config.ConnConfig.Config)
}

// buildDSN writes a keyword connection string for one parsed config. It
// turns off TLS when the parsed config holds no TLS setting.
func buildDSN(c *pgconn.Config) string {
	parts := []string{
		"host=" + c.Host,
		fmt.Sprintf("port=%d", c.Port),
		"user=" + c.User,
		"dbname=" + c.Database,
	}
	if c.Password != "" {
		parts = append(parts, "password="+c.Password)
	}
	if c.TLSConfig == nil {
		parts = append(parts, "sslmode=disable")
	}
	return strings.Join(parts, " ")
}

// testPool creates a connection pool for the named database.
func testPool(t *testing.T, dbName string) *pgxpool.Pool {
	t.Helper()
	url := testEnvURL(t, dbName)
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatalf("parse test database URL: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	t.Cleanup(func() {
		pool.Close()
	})
	return pool
}

// testDBName sanitises a test name into a database name.
func testDBName(testName string) string {
	name := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_' {
			return r
		}
		return '_'
	}, strings.ToLower(testName))
	return "kibtab_test_" + name
}

// testSetup creates a fresh database for an integration test.
// It returns the database name and a pool connected to it, plus a cleanup
// function that drops the database.
func testSetup(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	dbName := testDBName(t.Name())

	adminURL := testEnvURL(t, "postgres")
	adminConfig, err := pgxpool.ParseConfig(adminURL)
	if err != nil {
		t.Fatalf("parse admin URL: %v", err)
	}
	adminPool, err := pgxpool.NewWithConfig(context.Background(), adminConfig)
	if err != nil {
		t.Fatalf("create admin pool: %v", err)
	}
	defer adminPool.Close()

	ctx := context.Background()
	if _, err := adminPool.Exec(ctx,
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)); err != nil {
		t.Fatalf("drop database: %v", err)
	}
	if _, err := adminPool.Exec(ctx,
		fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
		t.Fatalf("create database: %v", err)
	}

	testURL := testEnvURL(t, dbName)
	testConfig, err := pgxpool.ParseConfig(testURL)
	if err != nil {
		t.Fatalf("parse test URL: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, testConfig)
	if err != nil {
		t.Fatalf("create test pool: %v", err)
	}

	cleanup := func() {
		pool.Close()
		adminPool2, _ := pgxpool.NewWithConfig(ctx, adminConfig)
		if adminPool2 == nil {
			return
		}
		defer adminPool2.Close()
		_, _ = adminPool2.Exec(ctx,
			fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	}

	return pool, cleanup
}

// createTestTable creates a table with the given name and fields.
// The primary key column is named k.
func createTestTable(t *testing.T, pool poolIface, tableName string, fields []string) {
	t.Helper()
	ctx := context.Background()
	cols := []string{"k text PRIMARY KEY"}
	for _, f := range fields {
		cols = append(cols, fmt.Sprintf("%s text", f))
	}
	stmt := fmt.Sprintf("CREATE TABLE %s (%s)",
		tableName, strings.Join(cols, ", "))
	if _, err := pool.Exec(ctx, stmt); err != nil {
		t.Fatalf("create table %s: %v", tableName, err)
	}
}

// truncateTable removes all rows from a table.
func truncateTable(t *testing.T, pool poolIface, tableName string) {
	t.Helper()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s", tableName)); err != nil {
		t.Fatalf("truncate %s: %v", tableName, err)
	}
}
