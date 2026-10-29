package ports

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kibtab/kibtab/internal/core/domain"
)

func TestFakesCompile(t *testing.T) {
	_, _, _, _, _, _, _, _ = newFakeDialect(), newFakeRowRepository(),
		newFakeTableRegistry(), newFakeAuditWriter(), newFakeAuditReader(),
		newFakeTransactionRunner(), newFakeClock(), newFakeSyncService()
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

func (f *fakeRowRepository) Read(_ context.Context,
	table string) ([]domain.CellValue, error) {
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

func (f *fakeRowRepository) ReadRow(_ context.Context,
	table string, row string) ([]domain.CellValue, error) {
	if table == "" || row == "" {
		return nil, errFakeRead{table: table}
	}
	return []domain.CellValue{{Field: "f1", Value: "old"}}, nil
}

func (f *fakeRowRepository) Write(_ context.Context,
	table string, row string, values []domain.CellValue) error {
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

func (f *fakeTableRegistry) List(context.Context) ([]domain.TableMetadata, error) {
	return []domain.TableMetadata{
		{
			Name:    "t1",
			Fields:  []string{"f1"},
			Version: 1,
		},
	}, nil
}

func (f *fakeTableRegistry) Metadata(_ context.Context,
	table string) (domain.TableMetadata, error) {
	if table == "" {
		return domain.TableMetadata{}, errors.New("table registry metadata failed for " + table)
	}
	return domain.TableMetadata{
		Name:    table,
		Fields:  []string{"f1"},
		Version: 1,
	}, nil
}

func (f *fakeTableRegistry) LockVersion(_ context.Context,
	table string) (int64, error) {
	if table == "" {
		return 0, errors.New("table registry lock failed for " + table)
	}
	return 1, nil
}

func (f *fakeTableRegistry) BumpVersion(_ context.Context,
	table string) (int64, error) {
	if table == "" {
		return 0, errors.New("table registry bump failed for " + table)
	}
	return 2, nil
}

func newFakeAuditWriter() AuditWriter {
	return &fakeAuditWriter{}
}

type fakeAuditWriter struct{}

func (f *fakeAuditWriter) Record(_ context.Context,
	entry domain.AuditRow) error {
	if entry.Table == "" {
		return errFakeAudit{table: entry.Table}
	}
	return nil
}

type errFakeAudit struct{ table string }

func (e errFakeAudit) Error() string {
	return "audit writer failed for " + e.table
}

func newFakeAuditReader() AuditReader {
	return &fakeAuditReader{}
}

type fakeAuditReader struct{}

func (f *fakeAuditReader) Read(_ context.Context, table string,
	row string, field string) ([]domain.AuditRow, error) {
	if table == "" || row == "" || field == "" {
		return nil, errors.New("audit reader failed")
	}
	return []domain.AuditRow{{Table: table, Row: row, Field: field}}, nil
}

func newFakeTransactionRunner() TransactionRunner {
	return &fakeTransactionRunner{}
}

type fakeTransactionRunner struct{}

func (f *fakeTransactionRunner) Run(fn func(context.Context) error) error {
	return fn(context.Background())
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
