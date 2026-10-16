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
			Row:     "r1",
			Field:   "name",
			Value:   "alice",
			Version: version,
		}},
	}
}

func TestSyncWritesEachRowGroupInOneTransaction(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 3}}
	repository := &fakeRepository{}
	dialect := &fakeDialect{coerced: map[string]string{"alice": "ALICE"}}
	runner := &fakeRunner{}
	service := NewSyncService(registry, repository, dialect, runner)

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
		if rowResult.Version != 3 {
			t.Fatalf("Sync() row %q version: got %d, want 3",
				rowResult.Row, rowResult.Version)
		}
	}
	if runner.runs != 2 || runner.commits != 2 || runner.rollbacks != 0 {
		t.Fatalf("runner: got %d runs, %d commits, %d rollbacks, want 2, 2, 0",
			runner.runs, runner.commits, runner.rollbacks)
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

func TestSyncConflictReturnsError(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 4}}
	repository := &fakeRepository{}
	runner := &fakeRunner{}
	service := NewSyncService(registry, repository, &fakeDialect{}, runner)

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
	if runner.rollbacks != 1 {
		t.Fatalf("runner rollbacks: got %d, want 1", runner.rollbacks)
	}
}

func TestSyncRollbackKeepsOtherGroups(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 1}}
	repository := &fakeRepository{failRow: "r1"}
	runner := &fakeRunner{}
	service := NewSyncService(registry, repository, &fakeDialect{}, runner)

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
}

func TestSyncCoercionErrorRollsBack(t *testing.T) {
	registry := &fakeRegistry{meta: domain.TableMetadata{Version: 1}}
	repository := &fakeRepository{}
	dialect := &fakeDialect{err: errors.New("bad value")}
	runner := &fakeRunner{}
	service := NewSyncService(registry, repository, dialect, runner)

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
	registry := &fakeRegistry{err: errors.New("no metadata")}
	repository := &fakeRepository{}
	runner := &fakeRunner{}
	service := NewSyncService(registry, repository, &fakeDialect{}, runner)

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

func TestSyncValidationError(t *testing.T) {
	service := NewSyncService(&fakeRegistry{}, &fakeRepository{}, &fakeDialect{}, &fakeRunner{})

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
	service := NewSyncService(engine, engine, engine, engine)

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
