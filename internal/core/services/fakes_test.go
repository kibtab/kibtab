package services

import (
	"context"
	"errors"
	"sync"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// writeCall records one call to fakeRepository.Write.
type writeCall struct {
	table  string
	row    string
	values []domain.CellValue
}

// fakeRegistry is a TableRegistry that returns fixed metadata.
// BumpVersion advances the metadata version and it counts each call.
type fakeRegistry struct {
	meta     domain.TableMetadata
	err      error
	bumpErr  error
	bumps    int
	lockErr  error
	lockHeld int
}

func (f *fakeRegistry) List(context.Context) ([]domain.TableMetadata, error) {
	if f.err != nil {
		return nil, f.err
	}
	return []domain.TableMetadata{f.meta}, nil
}

func (f *fakeRegistry) Metadata(_ context.Context,
	table string) (domain.TableMetadata, error) {
	if f.err != nil {
		return domain.TableMetadata{}, f.err
	}
	meta := f.meta
	meta.Name = table
	return meta, nil
}

func (f *fakeRegistry) LockVersion(context.Context, string) (int64, error) {
	if f.lockErr != nil {
		return 0, f.lockErr
	}
	f.lockHeld++
	return f.meta.Version, nil
}

func (f *fakeRegistry) BumpVersion(context.Context, string) (int64, error) {
	if f.bumpErr != nil {
		return 0, f.bumpErr
	}
	f.bumps++
	f.meta.Version++
	return f.meta.Version, nil
}

// fakeRepository is a RowRepository that records each write.
// It returns an error for the row named by failRow.
// ReadRow answers with the cells of oldRow.
type fakeRepository struct {
	failRow string
	readErr error
	oldRow  map[string][]domain.CellValue
	writes  []writeCall
	reads   int
}

func (f *fakeRepository) Read(context.Context,
	string) ([]domain.CellValue, error) {
	return nil, nil
}

func (f *fakeRepository) ReadRow(_ context.Context,
	_ string, row string) ([]domain.CellValue, error) {
	f.reads++
	if f.readErr != nil {
		return nil, f.readErr
	}
	return f.oldRow[row], nil
}

func (f *fakeRepository) Write(_ context.Context,
	table string, row string, values []domain.CellValue) error {
	f.writes = append(f.writes, writeCall{table: table, row: row, values: values})
	if row == f.failRow {
		return errors.New("write failed for row " + row)
	}
	return nil
}

// fakeAudit is an AuditWriter that records each audit row.
type fakeAudit struct {
	err     error
	entries []domain.AuditRow
}

func (f *fakeAudit) Record(_ context.Context, entry domain.AuditRow) error {
	if f.err != nil {
		return f.err
	}
	f.entries = append(f.entries, entry)
	return nil
}

// fakeDialect is a Dialect that can return a mapping error.
type fakeDialect struct {
	err     error
	coerced map[string]string
}

func (f *fakeDialect) Quote(name string) string { return name }

func (f *fakeDialect) MapValue(field string, value string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if coerced, ok := f.coerced[value]; ok {
		return coerced, nil
	}
	return value, nil
}

func (f *fakeDialect) Page(limit int, first int) (string, []any) {
	return "", []any{limit, first}
}

// fakeRunner is a TransactionRunner that records each begin, commit, and
// rollback. It hands a plain context to the function it runs.
type fakeRunner struct {
	runs      int
	commits   int
	rollbacks int
}

func (f *fakeRunner) Run(fn func(context.Context) error) error {
	f.runs++
	if err := fn(context.Background()); err != nil {
		f.rollbacks++
		return err
	}
	f.commits++
	return nil
}

// raceEngine is one fake that implements the TableRegistry, the
// RowRepository, the AuditWriter, the Dialect, and the
// TransactionRunner ports.
// Run holds the lock for the whole transaction, so one write commits
// before the next transaction reads the version.
// A successful write advances the table version in BumpVersion.
type raceEngine struct {
	mu      sync.Mutex
	version int64
}

func (e *raceEngine) Run(fn func(context.Context) error) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return fn(context.Background())
}

func (e *raceEngine) List(context.Context) ([]domain.TableMetadata, error) {
	return []domain.TableMetadata{{Version: e.version}}, nil
}

func (e *raceEngine) Metadata(_ context.Context,
	table string) (domain.TableMetadata, error) {
	return domain.TableMetadata{Name: table, Version: e.version}, nil
}

func (e *raceEngine) LockVersion(context.Context, string) (int64, error) {
	return e.version, nil
}

func (e *raceEngine) BumpVersion(context.Context, string) (int64, error) {
	e.version++
	return e.version, nil
}

func (e *raceEngine) Read(context.Context,
	string) ([]domain.CellValue, error) {
	return nil, nil
}

func (e *raceEngine) ReadRow(context.Context,
	string, string) ([]domain.CellValue, error) {
	return nil, nil
}

func (e *raceEngine) Write(context.Context,
	string, string, []domain.CellValue) error {
	return nil
}

func (e *raceEngine) Record(_ context.Context,
	entry domain.AuditRow) error {
	return nil
}

func (e *raceEngine) Quote(name string) string { return name }

func (e *raceEngine) MapValue(field string, value string) (string, error) {
	return value, nil
}

func (e *raceEngine) Page(limit int, first int) (string, []any) {
	return "", []any{limit, first}
}
