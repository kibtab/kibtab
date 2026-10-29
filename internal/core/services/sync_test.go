package services

import (
	"errors"
	"sync"
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// oneRowPayload returns a valid payload with one delta at one version.
func oneRowPayload(version int64) domain.SyncPayload {
	return domain.SyncPayload{
		Table: "t1",
		Delta: []domain.CellDelta{{
			Row:   "r1",
			Field: "name",
			Value: "alice", Version: version,
		}},
	}
}

// newService builds a SyncService with one fake for each port.
func newService(registry *fakeRegistry, repository *fakeRepository,
	audit *fakeAudit, dialect *fakeDialect,
	runner *fakeRunner) *SyncService {
	return NewSyncService(registry, repository, audit, dialect, runner)
}

func TestSyncWritesEachRowGroupInOneTransaction(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 3}}
	repository := &fakeRepository{}
	audit := &fakeAudit{}
	dialect := &fakeDialect{coerced: map[string]string{"alice": "ALICE"}}
	runner := &fakeRunner{}
	service := newService(registry, repository, audit, dialect, runner)

	payload := domain.SyncPayload{
		Table: "t1",
		Delta: []domain.CellDelta{
			{Row: "r2", Field: "name", Value: "alice", Version: 3},
			{Row: "r1", Field: "name", Value: "alice", Version: 3},
			{Row: "r1", Field: "email", Value: "alice@example.org", Version: 3},
		},
	}

	result, err := service.Sync(payload)
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}
	if len(result.RowResults) != 2 {
		t.Fatalf("Sync() row results: got %d, want 2", len(result.RowResults))
	}
	if result.RowResults[0].Row != "r1" || result.RowResults[1].Row != "r2" {
		t.Fatalf("Sync() row order: got %q then %q, want r1 then r2",
			result.RowResults[0].Row, result.RowResults[1].Row)
	}
	for _, rowResult := range result.RowResults {
		if rowResult.Error != nil {
			t.Fatalf("Sync() row %q error: %v", rowResult.Row, rowResult.Error)
		}
	}
	// Each written group bumps the table version by one, so the row
	// results hold versions 4 and 5.
	if result.RowResults[0].Version != 4 {
		t.Fatalf("Sync() first row version: got %d, want 4",
			result.RowResults[0].Version)
	}
	if result.RowResults[1].Version != 5 {
		t.Fatalf("Sync() second row version: got %d, want 5",
			result.RowResults[1].Version)
	}
	if runner.runs != 2 || runner.commits != 2 || runner.rollbacks != 0 {
		t.Fatalf("runner: got %d runs, %d commits, %d rollbacks, want 2, 2, 0",
			runner.runs, runner.commits, runner.rollbacks)
	}
	if registry.bumps != 2 {
		t.Fatalf("registry bumps: got %d, want 2", registry.bumps)
	}
	if len(repository.writes) != 2 {
		t.Fatalf("repository writes: got %d, want 2", len(repository.writes))
	}
	first := repository.writes[0]
	if first.row != "r1" || len(first.values) != 2 {
		t.Fatalf("first write: got row %q with %d values, want r1 with 2",
			first.row, len(first.values))
	}
	if first.values[0].Value != "ALICE" {
		t.Fatalf("first write value: got %q, want %q", first.values[0].Value, "ALICE")
	}
}

func TestSyncWritesOneAuditRowPerCellChange(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 1}}
	repository := &fakeRepository{
		oldRow: map[string][]domain.CellValue{
			"r1": {{Field: "name", Value: "bob"}},
		},
	}
	audit := &fakeAudit{}
	runner := &fakeRunner{}
	service := newService(registry, repository, audit, &fakeDialect{}, runner)

	payload := domain.SyncPayload{
		Table: "t1",
		User:  "rita",
		Delta: []domain.CellDelta{
			{Row: "r1", Field: "name", Value: "alice", Version: 1},
			{Row: "r1", Field: "email", Value: "alice@example.org", Version: 1},
		},
	}

	result, err := service.Sync(payload)
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}
	if result.RowResults[0].Error != nil {
		t.Fatalf("Sync() row error: %v", result.RowResults[0].Error)
	}
	if len(audit.entries) != 2 {
		t.Fatalf("audit rows: got %d, want 2", len(audit.entries))
	}
	name := audit.entries[0]
	if name.Table != "t1" || name.Row != "r1" || name.Field != "name" {
		t.Fatalf("audit row: got %+v, want the cell t1/r1/name", name)
	}
	if name.User != "rita" {
		t.Fatalf("audit user: got %q, want %q", name.User, "rita")
	}
	if name.OldValue != "bob" || name.NewValue != "alice" {
		t.Fatalf("audit values: got old %q new %q, want old bob new alice",
			name.OldValue, name.NewValue)
	}
	email := audit.entries[1]
	if email.Field != "email" || email.OldValue != "" {
		t.Fatalf("audit row: got %+v, want the cell t1/r1/email with no old value",
			email)
	}
}

func TestSyncConflictReturnsError(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 4}}
	repository := &fakeRepository{}
	audit := &fakeAudit{}
	runner := &fakeRunner{}
	service := newService(registry, repository, audit, &fakeDialect{}, runner)

	result, err := service.Sync(oneRowPayload(3))
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}
	if len(result.RowResults) != 1 {
		t.Fatalf("Sync() row results: got %d, want 1", len(result.RowResults))
	}
	rowResult := result.RowResults[0]
	if rowResult.Error == nil {
		t.Fatal("Sync() conflict: got no error, want one")
	}
	var conflict domain.ConflictError
	if !errors.As(rowResult.Error, &conflict) {
		t.Fatalf("Sync() error type: got %T, want ConflictError", rowResult.Error)
	}
	if conflict.Want != 3 || conflict.Got != 4 {
		t.Fatalf("conflict: got want %d, got %d, want want 3, got 4",
			conflict.Want, conflict.Got)
	}
	want := "table t1 row r1: version conflict, got 4, want 3"
	if conflict.Error() != want {
		t.Fatalf("conflict message: got %q, want %q", conflict.Error(), want)
	}
	if len(repository.writes) != 0 {
		t.Fatalf("repository writes: got %d, want 0", len(repository.writes))
	}
	if len(audit.entries) != 0 {
		t.Fatalf("audit rows: got %d, want 0", len(audit.entries))
	}
	if registry.bumps != 0 {
		t.Fatalf("registry bumps: got %d, want 0", registry.bumps)
	}
	if runner.rollbacks != 1 {
		t.Fatalf("runner rollbacks: got %d, want 1", runner.rollbacks)
	}
}

func TestSyncRollbackKeepsOtherGroups(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 1}}
	repository := &fakeRepository{failRow: "r1"}
	audit := &fakeAudit{}
	runner := &fakeRunner{}
	service := newService(registry, repository, audit, &fakeDialect{}, runner)

	payload := domain.SyncPayload{
		Table: "t1",
		Delta: []domain.CellDelta{
			{Row: "r1", Field: "name", Value: "alice", Version: 1},
			{Row: "r2", Field: "name", Value: "bob", Version: 1},
		},
	}

	result, err := service.Sync(payload)
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}
	if len(result.RowResults) != 2 {
		t.Fatalf("Sync() row results: got %d, want 2", len(result.RowResults))
	}
	if result.RowResults[0].Error == nil {
		t.Fatal("Sync() failed group: got no error, want one")
	}
	if result.RowResults[1].Error != nil {
		t.Fatalf("Sync() later group error: %v", result.RowResults[1].Error)
	}
	if runner.runs != 2 || runner.commits != 1 || runner.rollbacks != 1 {
		t.Fatalf("runner: got %d runs, %d commits, %d rollbacks, want 2, 1, 1",
			runner.runs, runner.commits, runner.rollbacks)
	}
	if len(repository.writes) != 2 {
		t.Fatalf("repository writes: got %d, want 2", len(repository.writes))
	}
	// The failed group keeps the expected version for the next group, so
	// the second group still matches the stored version.
	if registry.bumps != 1 {
		t.Fatalf("registry bumps: got %d, want 1", registry.bumps)
	}
}

func TestSyncCoercionErrorRollsBack(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 1}}
	repository := &fakeRepository{}
	audit := &fakeAudit{}
	dialect := &fakeDialect{err: errors.New("bad value")}
	runner := &fakeRunner{}
	service := newService(registry, repository, audit, dialect, runner)

	result, err := service.Sync(oneRowPayload(1))
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}
	if result.RowResults[0].Error == nil {
		t.Fatal("Sync() coercion: got no error, want one")
	}
	if len(repository.writes) != 0 {
		t.Fatalf("repository writes: got %d, want 0", len(repository.writes))
	}
	if runner.rollbacks != 1 {
		t.Fatalf("runner rollbacks: got %d, want 1", runner.rollbacks)
	}
}

func TestSyncRegistryErrorRollsBack(t *testing.T) {
	registry := &fakeRegistry{lockErr: errors.New("no metadata")}
	repository := &fakeRepository{}
	audit := &fakeAudit{}
	runner := &fakeRunner{}
	service := newService(registry, repository, audit, &fakeDialect{}, runner)

	result, err := service.Sync(oneRowPayload(1))
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}
	if result.RowResults[0].Error == nil {
		t.Fatal("Sync() registry: got no error, want one")
	}
	if len(repository.writes) != 0 {
		t.Fatalf("repository writes: got %d, want 0", len(repository.writes))
	}
	if runner.rollbacks != 1 {
		t.Fatalf("runner rollbacks: got %d, want 1", runner.rollbacks)
	}
}

func TestSyncReadRowErrorRollsBack(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 1}}
	repository := &fakeRepository{readErr: errors.New("no old value")}
	audit := &fakeAudit{}
	runner := &fakeRunner{}
	service := newService(registry, repository, audit, &fakeDialect{}, runner)

	result, err := service.Sync(oneRowPayload(1))
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}
	if result.RowResults[0].Error == nil {
		t.Fatal("Sync() read row: got no error, want one")
	}
	if len(repository.writes) != 0 {
		t.Fatalf("repository writes: got %d, want 0", len(repository.writes))
	}
	if runner.rollbacks != 1 {
		t.Fatalf("runner rollbacks: got %d, want 1", runner.rollbacks)
	}
}

func TestSyncAuditErrorRollsBack(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 1}}
	repository := &fakeRepository{}
	audit := &fakeAudit{err: errors.New("audit failed")}
	runner := &fakeRunner{}
	service := newService(registry, repository, audit, &fakeDialect{}, runner)

	result, err := service.Sync(oneRowPayload(1))
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}
	if result.RowResults[0].Error == nil {
		t.Fatal("Sync() audit: got no error, want one")
	}
	if registry.bumps != 0 {
		t.Fatalf("registry bumps: got %d, want 0", registry.bumps)
	}
	if runner.rollbacks != 1 {
		t.Fatalf("runner rollbacks: got %d, want 1", runner.rollbacks)
	}
}

func TestSyncBumpErrorRollsBack(t *testing.T) {
	registry := &fakeRegistry{
		meta:    domain.TableMetadata{Version: 1},
		bumpErr: errors.New("bump failed"),
	}
	repository := &fakeRepository{}
	audit := &fakeAudit{}
	runner := &fakeRunner{}
	service := newService(registry, repository, audit, &fakeDialect{}, runner)

	result, err := service.Sync(oneRowPayload(1))
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}
	if result.RowResults[0].Error == nil {
		t.Fatal("Sync() bump: got no error, want one")
	}
	if runner.rollbacks != 1 {
		t.Fatalf("runner rollbacks: got %d, want 1", runner.rollbacks)
	}
}

func TestSyncValidationError(t *testing.T) {
	service := newService(&fakeRegistry{}, &fakeRepository{}, &fakeAudit{},
		&fakeDialect{}, &fakeRunner{})

	result, err := service.Sync(domain.SyncPayload{})
	if err == nil {
		t.Fatal("Sync() invalid payload: got no error, want one")
	}
	if result.RowResults != nil {
		t.Fatalf("Sync() invalid payload results: got %v, want nil", result.RowResults)
	}
}

func TestSyncRaceOneWriteWins(t *testing.T) {
	engine := &raceEngine{version: 1}
	audit := &fakeAudit{}
	service := NewSyncService(engine, engine, audit, engine, engine)

	const writers = 2
	results := make([]domain.SyncResult, writers)
	errs := make([]error, writers)

	var wg sync.WaitGroup
	wg.Add(writers)
	for i := 0; i < writers; i++ {
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = service.Sync(oneRowPayload(1))
		}(i)
	}
	wg.Wait()

	wins := 0
	conflicts := 0
	for i := 0; i < writers; i++ {
		if errs[i] != nil {
			t.Fatalf("writer %d: Sync() error: %v", i, errs[i])
		}
		for _, rowResult := range results[i].RowResults {
			switch {
			case rowResult.Error == nil:
				wins++
			case isConflict(rowResult.Error):
				conflicts++
			default:
				t.Fatalf("writer %d: unexpected error: %v", i, rowResult.Error)
			}
		}
	}
	if wins != 1 {
		t.Fatalf("race: got %d winning writes, want 1", wins)
	}
	if conflicts != 1 {
		t.Fatalf("race: got %d conflicting writes, want 1", conflicts)
	}
}

// isConflict reports whether err is a domain.ConflictError.
func isConflict(err error) bool {
	var conflict domain.ConflictError
	return errors.As(err, &conflict)
}
