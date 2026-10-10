package ports

import (
	"errors"
	"testing"
	"time"

	"github.com/kibtab/kibtab/internal/core/domain"
)

func TestFakesCompile(t *testing.T) {
	_, _, _, _, _, _, _ = newFakeDialect(), newFakeRowRepository(),
		newFakeTableRegistry(), newFakeAuditWriter(), newFakeTransactionRunner(),
		newFakeClock(), newFakeSyncService()
}

func newFakeDialect() Dialect {
	return &fakeDialect{}
}

type fakeDialect struct{}

func (f *fakeDialect) Quote(name string) string {
	return "\"" + name + "\""
}

func (f *fakeDialect) MapValue(field string, value string) (string, error) {
	if field == "" {
		return "", errors.New("field is empty")
	}
	return value, nil
}

func (f *fakeDialect) Page(limit int, first int) (string, []any) {
	return " LIMIT ?", []any{limit}
}

func newFakeRowRepository() RowRepository {
	return &fakeRowRepository{}
}

type fakeRowRepository struct{}

func (f *fakeRowRepository) Read(table string) ([]domain.CellValue, error) {
	if table == "" {
		return nil, errFakeRead{table: table}
	}
	return []domain.CellValue{{
		Field: "f1",
		Value: "v1",
	}}, nil
}

type errFakeRead struct{ table string }

func (e errFakeRead) Error() string {
	return "row repository read failed for " + e.table
}

func (f *fakeRowRepository) Write(table string, row string, values []domain.CellValue) error {
	if table == "" {
		return errFakeWrite{table: table}
	}
	return nil
}

type errFakeWrite struct{ table string }

func (e errFakeWrite) Error() string {
	return "row repository write failed for " + e.table
}

func newFakeTableRegistry() TableRegistry {
	return &fakeTableRegistry{}
}

type fakeTableRegistry struct{}

func (f *fakeTableRegistry) List() ([]domain.TableMetadata, error) {
	return []domain.TableMetadata{
		{
			Name:    "t1",
			Fields:  []string{"f1"},
			Version: 1,
		},
	}, nil
}

func (f *fakeTableRegistry) Metadata(table string) (domain.TableMetadata, error) {
	if table == "" {
		return domain.TableMetadata{}, errors.New("table registry metadata failed for " + table)
	}
	return domain.TableMetadata{
		Name:    table,
		Fields:  []string{"f1"},
		Version: 1,
	}, nil
}

func newFakeAuditWriter() AuditWriter {
	return &fakeAuditWriter{}
}

type fakeAuditWriter struct{}

func (f *fakeAuditWriter) Write(table string, row string, field string, oldValue string, newValue string) error {
	if table == "" {
		return errFakeAudit{table: table}
	}
	return nil
}

type errFakeAudit struct{ table string }

func (e errFakeAudit) Error() string {
	return "audit writer failed for " + e.table
}

func newFakeTransactionRunner() TransactionRunner {
	return &fakeTransactionRunner{}
}

type fakeTransactionRunner struct{}

func (f *fakeTransactionRunner) Run(fn func() error) error {
	return fn()
}

func newFakeClock() Clock {
	return &fakeClock{}
}

type fakeClock struct{}

func (f *fakeClock) Now() time.Time { return time.Time{} }

func newFakeSyncService() SyncService {
	return &fakeSyncService{}
}

type fakeSyncService struct{}

func (f *fakeSyncService) Sync(payload domain.SyncPayload) (domain.SyncResult, error) {
	if payload.Table == "" {
		return domain.SyncResult{}, errors.New("sync payload table is empty")
	}
	return domain.SyncResult{
		RowResults: []domain.RowResult{{
			Row:     payload.Delta[0].Row,
			Version: payload.Delta[0].Version,
		}},
	}, nil
}
