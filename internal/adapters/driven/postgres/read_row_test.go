package postgres

import "testing"

func TestReadRowColumnQueryError(t *testing.T) {
	pool := &fakePool{queryErr: errDummy}
	repo := NewRowRepository(pool, Dialect{})

	_, err := repo.ReadRow(ctx, "users", "r1")
	if err == nil {
		t.Fatal("ReadRow with a column query error: got nil, want error")
	}
}

func TestReadRowWithoutFields(t *testing.T) {
	pool := &fakePool{rows: &fakeRows{}}
	repo := NewRowRepository(pool, Dialect{})

	cells, err := repo.ReadRow(ctx, "users", "r1")
	if err != nil {
		t.Fatalf("ReadRow without fields: unexpected error: %v", err)
	}
	if cells != nil {
		t.Fatalf("ReadRow without fields: got %v, want nil", cells)
	}
}

func TestReadRowQueryError(t *testing.T) {
	pool := &multiFakePool{
		fieldRows:      fieldMetadataRows(),
		selectQueryErr: errDummy,
	}
	repo := NewRowRepository(pool, Dialect{})

	_, err := repo.ReadRow(ctx, "users", "r1")
	if err == nil {
		t.Fatal("ReadRow with a query error: got nil, want error")
	}
}

func TestReadRowMissingRow(t *testing.T) {
	pool := &multiFakePool{
		fieldRows:  fieldMetadataRows(),
		selectRows: &fakeRows{},
	}
	repo := NewRowRepository(pool, Dialect{})

	cells, err := repo.ReadRow(ctx, "users", "r1")
	if err != nil {
		t.Fatalf("ReadRow of a missing row: unexpected error: %v", err)
	}
	if cells != nil {
		t.Fatalf("ReadRow of a missing row: got %v, want nil", cells)
	}
}

func TestReadRowRowsErr(t *testing.T) {
	pool := &multiFakePool{
		fieldRows:  fieldMetadataRows(),
		selectRows: &fakeRows{err: errDummy},
	}
	repo := NewRowRepository(pool, Dialect{})

	_, err := repo.ReadRow(ctx, "users", "r1")
	if err == nil {
		t.Fatal("ReadRow with rows.Err: got nil, want error")
	}
}

func TestReadRowScanError(t *testing.T) {
	pool := &multiFakePool{
		fieldRows:  fieldMetadataRows(),
		selectRows: &fakeRows{next: true, scanErr: errDummy},
	}
	repo := NewRowRepository(pool, Dialect{})

	_, err := repo.ReadRow(ctx, "users", "r1")
	if err == nil {
		t.Fatal("ReadRow with a scan error: got nil, want error")
	}
}

func TestReadRowReturnsTheCells(t *testing.T) {
	pool := &multiFakePool{
		fieldRows:  fieldMetadataRows(),
		selectRows: &fakeRows{next: true, scanVals: []any{"alice"}},
	}
	repo := NewRowRepository(pool, Dialect{})

	cells, err := repo.ReadRow(ctx, "users", "r1")
	if err != nil {
		t.Fatalf("ReadRow: %v", err)
	}
	if len(cells) != 1 || cells[0].Field != "name" || cells[0].Value != "alice" {
		t.Fatalf("ReadRow: got %+v, want the name cell of alice", cells)
	}
}
