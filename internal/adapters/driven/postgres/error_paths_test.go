package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kibtab/kibtab/internal/core/domain"
)

// --- RowRepository error paths ---

// multiFakePool returns different results for successive Query calls.
// First Query (tableFields): returns fieldRows.
// Second Query (SELECT in Read): returns selectRows and/or selectQueryErr.
type multiFakePool struct {
	queryCalls     int
	fieldRows      pgx.Rows
	selectRows     pgx.Rows
	selectQueryErr error
}

func (m *multiFakePool) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	m.queryCalls++
	if m.queryCalls == 1 {
		return m.fieldRows, nil
	}
	return m.selectRows, m.selectQueryErr
}

func (m *multiFakePool) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return &fakeRow{}
}

func (m *multiFakePool) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag("INSERT 0 1"), nil
}

func (m *multiFakePool) Ping(_ context.Context) error { return nil }
func (m *multiFakePool) Close()                       {}

// fieldMetadataRows returns one row with column_name "name".
func fieldMetadataRows() *fakeRows {
	return &fakeRows{
		next:     true,
		scanVals: []any{"name"},
	}
}

func TestReadQueryError(t *testing.T) {
	pool := &multiFakePool{
		fieldRows:      fieldMetadataRows(),
		selectQueryErr: errDummy,
	}
	repo := NewRowRepository(pool, Dialect{})
	_, err := repo.Read("users")
	if err == nil {
		t.Fatal("Read with SELECT query error: got nil, want error")
	}
}

func TestReadScanError(t *testing.T) {
	selectRows := &fakeRows{
		next:    true,
		scanErr: errDummy,
	}
	pool := &multiFakePool{
		fieldRows:  fieldMetadataRows(),
		selectRows: selectRows,
	}
	repo := NewRowRepository(pool, Dialect{})
	_, err := repo.Read("users")
	if err == nil {
		t.Fatal("Read with scan error: got nil, want error")
	}
}

func TestReadRowsErr(t *testing.T) {
	selectRows := &fakeRows{
		next:     true,
		scanVals: []any{"test"},
		err:      errDummy,
	}
	pool := &multiFakePool{
		fieldRows:  fieldMetadataRows(),
		selectRows: selectRows,
	}
	repo := NewRowRepository(pool, Dialect{})
	_, err := repo.Read("users")
	if err == nil {
		t.Fatal("Read with rows.Err: got nil, want error")
	}
}

func TestWriteExecError(t *testing.T) {
	pool := &fakePool{
		execErr: errDummy,
	}
	repo := NewRowRepository(pool, Dialect{})
	err := repo.Write("users", "r1", []domain.CellValue{
		{Field: "name", Value: "test"},
	})
	if err == nil {
		t.Fatal("Write with exec error: got nil, want error")
	}
}

func TestWriteNoRowsAffected(t *testing.T) {
	pool := &fakePool{
		execNoRows: true,
	}
	repo := NewRowRepository(pool, Dialect{})
	err := repo.Write("users", "r1", []domain.CellValue{
		{Field: "name", Value: "test"},
	})
	if err == nil {
		t.Fatal("Write with no rows affected: got nil, want error")
	}
}

// --- tableFields error paths (first Query in Read) ---

func TestReadTableFieldsScanError(t *testing.T) {
	rows := &fakeRows{
		next:    true,
		scanErr: errDummy,
	}
	pool := &fakePool{
		rows: rows,
	}
	repo := NewRowRepository(pool, Dialect{})
	_, err := repo.Read("users")
	if err == nil {
		t.Fatal("Read with tableFields scan error: got nil, want error")
	}
}

func TestReadTableFieldsRowsErr(t *testing.T) {
	rows := &fakeRows{
		next:     true,
		scanVals: []any{"name"},
		err:      errDummy,
	}
	pool := &fakePool{
		rows: rows,
	}
	repo := NewRowRepository(pool, Dialect{})
	_, err := repo.Read("users")
	if err == nil {
		t.Fatal("Read with tableFields rows.Err: got nil, want error")
	}
}

// --- TableRegistry.columns error paths ---

func TestColumnsScanError(t *testing.T) {
	rows := &fakeRows{
		next:    true,
		scanErr: errDummy,
	}
	pool := &fakePool{
		rows: rows,
	}
	reg := NewTableRegistry(pool, []string{"users"})
	_, err := reg.Metadata("users")
	if err == nil {
		t.Fatal("Metadata with columns scan error: got nil, want error")
	}
}

func TestColumnsRowsErr(t *testing.T) {
	rows := &fakeRows{
		next:     true,
		scanVals: []any{"name"},
		err:      errDummy,
	}
	pool := &fakePool{
		rows: rows,
	}
	reg := NewTableRegistry(pool, []string{"users"})
	_, err := reg.Metadata("users")
	if err == nil {
		t.Fatal("Metadata with columns rows.Err: got nil, want error")
	}
}

// --- runMigrations error paths ---

// fakeExecer implements migrationExecer for error-path testing.
type fakeExecer struct {
	execCount int
	failAt    int // 1-based call index at which Exec returns errDummy
	execErr   error
}

func (f *fakeExecer) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	f.execCount++
	if f.failAt > 0 && f.execCount == f.failAt {
		return pgconn.NewCommandTag(""), f.execErr
	}
	return pgconn.NewCommandTag("SELECT 1"), nil
}

func TestRunMigrationsAcquireLockError(t *testing.T) {
	execer := &fakeExecer{
		failAt:  1,
		execErr: errDummy,
	}
	err := runMigrations(context.Background(), execer)
	if err == nil {
		t.Fatal("runMigrations with lock error: got nil, want error")
	}
}

func TestRunMigrationsReleaseLockError(t *testing.T) {
	// Number of Exec calls: 1 (lock) + 3 (migration stmts) + 1 (unlock) = 5.
	execer := &fakeExecer{
		failAt:  5,
		execErr: errDummy,
	}
	err := runMigrations(context.Background(), execer)
	if err == nil {
		t.Fatal("runMigrations with unlock error: got nil, want error")
	}
}

func TestRunMigrationsMigrationError(t *testing.T) {
	// Fail on the first migration statement (call 2: after the lock).
	execer := &fakeExecer{
		failAt:  2,
		execErr: errDummy,
	}
	err := runMigrations(context.Background(), execer)
	if err == nil {
		t.Fatal("runMigrations with migration error: got nil, want error")
	}
}

// --- TransactionRunner error paths ---

type fakeTx struct {
	pgx.Tx
	commitErr   error
	rollbackErr error
}

func (f *fakeTx) Commit(_ context.Context) error   { return f.commitErr }
func (f *fakeTx) Rollback(_ context.Context) error { return f.rollbackErr }

type fakeTxBeginner struct {
	tx  pgx.Tx
	err error
}

func (f *fakeTxBeginner) BeginTx(_ context.Context, _ pgx.TxOptions) (pgx.Tx, error) {
	return f.tx, f.err
}

func TestTransactionRunnerBeginTxError(t *testing.T) {
	tr := NewTransactionRunner(&fakeTxBeginner{
		tx:  &fakeTx{},
		err: errDummy,
	})
	err := tr.Run(func() error { return nil })
	if err == nil {
		t.Fatal("Run with begin error: got nil, want error")
	}
}

func TestTransactionRunnerRollbackError(t *testing.T) {
	tr := NewTransactionRunner(&fakeTxBeginner{
		tx: &fakeTx{
			rollbackErr: errDummy,
		},
	})
	err := tr.Run(func() error { return errors.New("fn error") })
	if err == nil {
		t.Fatal("Run with rollback error: got nil, want error")
	}
}

func TestTransactionRunnerCommitError(t *testing.T) {
	tr := NewTransactionRunner(&fakeTxBeginner{
		tx: &fakeTx{
			commitErr: errDummy,
		},
	})
	err := tr.Run(func() error { return nil })
	if err == nil {
		t.Fatal("Run with commit error: got nil, want error")
	}
}
