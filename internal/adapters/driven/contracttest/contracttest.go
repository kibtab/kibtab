// Package contracttest holds the shared contract suites of Kibtab.
//
// A suite holds an adapter to a port. Each engine adapter runs the same
// suite. A second adapter must not copy a suite. Release v0.3.0 adds the
// first suite.
//
// A suite returns its failures. It does not call the testing package.
// The caller reports each failure with its own test.
package contracttest

import (
	"fmt"

	"github.com/kibtab/kibtab/internal/core/domain"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// Fixture holds what one engine adapter gives to the RowRepository suite.
type Fixture struct {
	// Repo is the adapter under test.
	Repo ports.RowRepository

	// Table is the table that the suite writes.
	Table string

	// Key names the row key column.
	Key string

	// Reset empties the table. The suite calls it before each case.
	Reset func()
}

// fields holds the column names that the suite writes. The caller must
// create the table with these fields beside the row key column.
var fields = []string{"name", "email"}

// RowRepository runs the shared contract suite for the RowRepository port.
// Each engine adapter calls it with its own fixture. It returns one error
// for each failed case.
func RowRepository(f Fixture) []error {
	all := []error{}
	all = append(all, writeAndReadOneRow(f)...)
	all = append(all, writeAndReadTwoRows(f)...)
	all = append(all, writeUpdatesExistingRow(f)...)
	all = append(all, writeEmptyValues(f)...)
	all = append(all, writeNulByte(f)...)
	all = append(all, writeEmptyField(f)...)
	all = append(all, readMissingTable(f)...)
	return dropNil(all)
}

// dropNil returns the errors that hold a failure.
func dropNil(errs []error) []error {
	out := []error{}
	for _, err := range errs {
		if err != nil {
			out = append(out, err)
		}
	}
	return out
}

// writeAndReadOneRow checks that one row survives a write and a read.
func writeAndReadOneRow(f Fixture) []error {
	const label = "write and read one row"
	f.Reset()
	errs := []error{}
	row := "r1"
	values := []domain.CellValue{
		{Field: "name", Value: "alice"},
		{Field: "email", Value: "alice@example.org"},
	}
	errs = append(errs, wrap(label+": Write", f.Repo.Write(f.Table, row, values)))
	cells, err := f.Repo.Read(f.Table)
	errs = append(errs, wrap(label+": Read", err))
	errs = append(errs, checkMin(label+": Read cell count", len(cells), 1))
	if len(cells) > 0 {
		errs = append(errs,
			check(label+": first cell field", cells[0].Field, f.Key),
			check(label+": row key value", cells[0].Value, row),
		)
	}
	foundName := false
	foundEmail := false
	for _, c := range cells {
		if c.Field == "name" && c.Value == "alice" {
			foundName = true
		}
		if c.Field == "email" && c.Value == "alice@example.org" {
			foundEmail = true
		}
	}
	errs = append(errs, checkBool(label+": name cell", foundName))
	errs = append(errs, checkBool(label+": email cell", foundEmail))
	return errs
}

// writeAndReadTwoRows checks that the key of each row survives a read.
func writeAndReadTwoRows(f Fixture) []error {
	const label = "write and read two rows"
	f.Reset()
	errs := []error{}
	errs = append(errs, wrap(label+": Write r1", f.Repo.Write(f.Table, "r1",
		[]domain.CellValue{{Field: "name", Value: "alice"}})))
	errs = append(errs, wrap(label+": Write r2", f.Repo.Write(f.Table, "r2",
		[]domain.CellValue{{Field: "name", Value: "bob"}})))
	cells, err := f.Repo.Read(f.Table)
	errs = append(errs, wrap(label+": Read", err))
	errs = append(errs, checkMin(label+": Read cell count", len(cells), 4))
	keys := map[string]bool{}
	for i := 0; i < len(cells); i += 1 + len(fields) {
		if cells[i].Field == f.Key {
			keys[cells[i].Value] = true
		}
	}
	errs = append(errs, checkBool(label+": row r1", keys["r1"]))
	errs = append(errs, checkBool(label+": row r2", keys["r2"]))
	return errs
}

// writeUpdatesExistingRow checks that a second write replaces the value.
func writeUpdatesExistingRow(f Fixture) []error {
	const label = "write updates existing row"
	f.Reset()
	errs := []error{}
	errs = append(errs, wrap(label+": Write v1", f.Repo.Write(f.Table, "r1",
		[]domain.CellValue{{Field: "name", Value: "alice"}})))
	errs = append(errs, wrap(label+": Write v2", f.Repo.Write(f.Table, "r1",
		[]domain.CellValue{{Field: "name", Value: "alice2"}})))
	cells, err := f.Repo.Read(f.Table)
	errs = append(errs, wrap(label+": Read", err))
	found := false
	for _, c := range cells {
		if c.Field == "name" && c.Value == "alice2" {
			found = true
		}
	}
	errs = append(errs, checkBool(label+": updated value", found))
	return errs
}

// writeEmptyValues checks that a write with no value changes nothing.
func writeEmptyValues(f Fixture) []error {
	const label = "write with empty values is a no-op"
	f.Reset()
	return []error{wrap(label, f.Repo.Write(f.Table, "r1", nil))}
}

// writeNulByte checks that a NUL byte in a value returns an error.
func writeNulByte(f Fixture) []error {
	const label = "write with NUL byte in value returns error"
	f.Reset()
	return []error{checkError(label, f.Repo.Write(f.Table, "r1",
		[]domain.CellValue{{Field: "name", Value: "bad\x00value"}}))}
}

// writeEmptyField checks that an empty field name returns an error.
func writeEmptyField(f Fixture) []error {
	const label = "write with empty field name returns error"
	f.Reset()
	return []error{checkError(label, f.Repo.Write(f.Table, "r1",
		[]domain.CellValue{{Field: "", Value: "v"}}))}
}

// readMissingTable checks that a read of a missing table returns no cells.
func readMissingTable(f Fixture) []error {
	const label = "read from non-existent table returns nil"
	cells, err := f.Repo.Read("does_not_exist_999")
	return []error{
		wrap(label+": Read", err),
		checkNilCells(label, cells),
	}
}

// check returns an error when got differs from want.
func check(label string, got string, want string) error {
	if got == want {
		return nil
	}
	return fmt.Errorf("%s: got %q, want %q", label, got, want)
}

// checkMin returns an error when got sits below want.
func checkMin(label string, got int, want int) error {
	if got >= want {
		return nil
	}
	return fmt.Errorf("%s: got %d, want at least %d", label, got, want)
}

// checkBool returns an error when ok is false.
func checkBool(label string, ok bool) error {
	if ok {
		return nil
	}
	return fmt.Errorf("%s: got false, want true", label)
}

// checkError returns an error when err is nil.
func checkError(label string, err error) error {
	if err != nil {
		return nil
	}
	return fmt.Errorf("%s: got no error, want one", label)
}

// checkNilCells returns an error when cells is not nil.
func checkNilCells(label string, cells []domain.CellValue) error {
	if cells == nil {
		return nil
	}
	return fmt.Errorf("%s: got %v, want nil", label, cells)
}

// wrap returns an error with a label when err is not nil.
func wrap(label string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", label, err)
}
