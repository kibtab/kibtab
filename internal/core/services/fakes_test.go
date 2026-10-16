package services

import (
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
type fakeRegistry struct {
	meta domain.TableMetadata
	err  error
}

func (f *fakeRegistry) List() ([]domain.TableMetadata, error) {
	if f.err != nil {
		return nil, f.err
	}
	return []domain.TableMetadata{f.meta}, nil
}

func (f *fakeRegistry) Metadata(table string) (domain.TableMetadata, error) {
	if f.err != nil {
		return domain.TableMetadata{}, f.err
	}
	meta := f.meta
	meta.Name = table
	return meta, nil
}

// fakeRepository is a RowRepository that records each write.
// It returns an error for the row named by failRow.
type fakeRepository struct {
	failRow string
	writes  []writeCall
}

func (f *fakeRepository) Read(string) ([]domain.CellValue, error) {
	return nil, nil
}

func (f *fakeRepository) Write(table string, row string, values []domain.CellValue) error {
	f.writes = append(f.writes, writeCall{table: table, row: row, values: values})
	if row == f.failRow {
		return errors.New("write failed for row " + row)
	}
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
// rollback.
type fakeRunner struct {
	runs      int
	commits   int
	rollbacks int
}

func (f *fakeRunner) Run(fn func() error) error {
	f.runs++
	if err := fn(); err != nil {
		f.rollbacks++
		return err
	}
	f.commits++
	return nil
}

// raceEngine is one fake that implements the TableRegistry, the
// RowRepository, the Dialect, and the TransactionRunner ports.
// Run holds the lock for the whole transaction, so one write commits
// before the next transaction reads the version.
// A successful write advances the table version.
type raceEngine struct {
	mu      sync.Mutex
	version int64
}

func (e *raceEngine) Run(fn func() error) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return fn()
}

func (e *raceEngine) List() ([]domain.TableMetadata, error) {
	return []domain.TableMetadata{{Version: e.version}}, nil
}

func (e *raceEngine) Metadata(table string) (domain.TableMetadata, error) {
	return domain.TableMetadata{Name: table, Version: e.version}, nil
}

func (e *raceEngine) Read(string) ([]domain.CellValue, error) {
	return nil, nil
}

func (e *raceEngine) Write(string, string, []domain.CellValue) error {
	e.version++
	return nil
}

func (e *raceEngine) Quote(name string) string { return name }

func (e *raceEngine) MapValue(field string, value string) (string, error) {
	return value, nil
}

func (e *raceEngine) Page(limit int, first int) (string, []any) {
	return "", []any{limit, first}
}
