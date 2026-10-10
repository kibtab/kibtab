package contracttest

import (
	"errors"
	"strings"
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

const testKey = "k"

// fakeRepo is an in-memory RowRepository. It satisfies the suite so that
// the suite logic is covered without a database. It emits every field of
// the table for each row, as a database adapter does.
type fakeRepo struct {
	key   string
	order []string
	rows  map[string]map[string]string
}

func newFakeRepo(key string) *fakeRepo {
	return &fakeRepo{key: key, rows: map[string]map[string]string{}}
}

func (f *fakeRepo) Write(table string, row string, values []domain.CellValue) error {
	for _, v := range values {
		if v.Field == "" {
			return errors.New("field name is empty")
		}
		if strings.ContainsRune(v.Value, 0) {
			return errors.New("value contains a NUL byte")
		}
	}
	if len(values) == 0 {
		return nil
	}
	id := table + "/" + row
	if _, seen := f.rows[id]; !seen {
		f.order = append(f.order, id)
		f.rows[id] = map[string]string{}
	}
	for _, v := range values {
		f.rows[id][v.Field] = v.Value
	}
	return nil
}

func (f *fakeRepo) Read(table string) ([]domain.CellValue, error) {
	var cells []domain.CellValue
	for _, id := range f.order {
		if !strings.HasPrefix(id, table+"/") {
			continue
		}
		cells = append(cells, domain.CellValue{
			Field: f.key,
			Value: strings.TrimPrefix(id, table+"/"),
		})
		for _, field := range fields {
			cells = append(cells, domain.CellValue{
				Field: field,
				Value: f.rows[id][field],
			})
		}
	}
	return cells, nil
}

// testFixture builds a fixture over a fresh in-memory repository.
func testFixture(repo *fakeRepo) Fixture {
	return Fixture{
		Repo:  repo,
		Table: "users",
		Key:   testKey,
		Reset: func() {
			repo.order = nil
			repo.rows = map[string]map[string]string{}
		},
	}
}

// TestRowRepositoryPassesWithACompliantAdapter checks that the suite
// reports no failure when an adapter meets the port.
func TestRowRepositoryPassesWithACompliantAdapter(t *testing.T) {
	errs := RowRepository(testFixture(newFakeRepo(testKey)))
	for _, err := range errs {
		t.Errorf("suite reported a failure for a compliant adapter: %v", err)
	}
}

// TestRowRepositoryReportsAnEmptyRead checks that the suite reports a
// missing row. It also reaches the branch for an empty cell list.
func TestRowRepositoryReportsAnEmptyRead(t *testing.T) {
	repo := &emptyRepo{}
	f := Fixture{
		Repo:  repo,
		Table: "users",
		Key:   testKey,
		Reset: func() {},
	}
	errs := RowRepository(f)
	if len(errs) == 0 {
		t.Fatal("suite reported no failure for an adapter that reads nothing")
	}
}

// emptyRepo reads no cells and writes nothing.
type emptyRepo struct{}

func (emptyRepo) Write(string, string, []domain.CellValue) error { return nil }

func (emptyRepo) Read(string) ([]domain.CellValue, error) { return nil, nil }

// TestHelpersPassOnAGoodValue checks that each helper returns no error.
func TestHelpersPassOnAGoodValue(t *testing.T) {
	if err := check("l", "a", "a"); err != nil {
		t.Errorf("check: got %v, want nil", err)
	}
	if err := checkMin("l", 2, 2); err != nil {
		t.Errorf("checkMin: got %v, want nil", err)
	}
	if err := checkBool("l", true); err != nil {
		t.Errorf("checkBool: got %v, want nil", err)
	}
	if err := checkError("l", errors.New("boom")); err != nil {
		t.Errorf("checkError: got %v, want nil", err)
	}
	if err := checkNilCells("l", nil); err != nil {
		t.Errorf("checkNilCells: got %v, want nil", err)
	}
	if err := wrap("l", nil); err != nil {
		t.Errorf("wrap: got %v, want nil", err)
	}
}

// TestHelpersReportABadValue checks the message that each helper returns.
func TestHelpersReportABadValue(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "check", err: check("l", "a", "b"), want: `l: got "a", want "b"`},
		{name: "checkMin", err: checkMin("l", 1, 2), want: "l: got 1, want at least 2"},
		{name: "checkBool", err: checkBool("l", false), want: "l: got false, want true"},
		{name: "checkError", err: checkError("l", nil), want: "l: got no error, want one"},
		{name: "wrap", err: wrap("l", errors.New("boom")), want: "l: boom"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.err == nil {
				t.Fatalf("%s: got no error, want one", test.name)
			}
			if test.want != "" && test.err.Error() != test.want {
				t.Fatalf("%s: got %q, want %q", test.name, test.err, test.want)
			}
		})
	}
}

// TestDropNilKeepsOnlyFailures checks that dropNil drops an empty error.
func TestDropNilKeepsOnlyFailures(t *testing.T) {
	got := dropNil([]error{nil, errors.New("boom"), nil})
	if len(got) != 1 {
		t.Fatalf("dropNil returned %d errors, want 1", len(got))
	}
	if got[0].Error() != "boom" {
		t.Fatalf("dropNil returned %q, want %q", got[0], "boom")
	}
}

// TestCheckNilCellsReportsCells checks the message for a cell list.
func TestCheckNilCellsReportsCells(t *testing.T) {
	err := checkNilCells("l", []domain.CellValue{{Field: "f", Value: "v"}})
	if err == nil {
		t.Fatal("checkNilCells: got no error, want one")
	}
	want := "l: got [{f v}], want nil"
	if err.Error() != want {
		t.Fatalf("checkNilCells: got %q, want %q", err, want)
	}
}
