package domain

import "testing"

func TestTableMetadataZeroValue(t *testing.T) {
	tm := TableMetadata{}
	if tm.Name != "" {
		t.Fatalf("TableMetadata.Name = %q, want empty", tm.Name)
	}
	if tm.Fields != nil {
		t.Fatalf("TableMetadata.Fields = %v, want nil", tm.Fields)
	}
	if tm.Version != 0 {
		t.Fatalf("TableMetadata.Version = %d, want 0", tm.Version)
	}
}

func TestCellValueZeroValue(t *testing.T) {
	cv := CellValue{}
	if cv.Field != "" {
		t.Fatalf("CellValue.Field = %q, want empty", cv.Field)
	}
	if cv.Value != "" {
		t.Fatalf("CellValue.Value = %q, want empty", cv.Value)
	}
}

func TestCellDeltaZeroValue(t *testing.T) {
	cd := CellDelta{}
	if cd.Row != "" {
		t.Fatalf("CellDelta.Row = %q, want empty", cd.Row)
	}
	if cd.Field != "" {
		t.Fatalf("CellDelta.Field = %q, want empty", cd.Field)
	}
	if cd.Value != "" {
		t.Fatalf("CellDelta.Value = %q, want empty", cd.Value)
	}
	if cd.Version != 0 {
		t.Fatalf("CellDelta.Version = %d, want 0", cd.Version)
	}
}

func TestSyncPayloadZeroValue(t *testing.T) {
	sp := SyncPayload{}
	if sp.Table != "" {
		t.Fatalf("SyncPayload.Table = %q, want empty", sp.Table)
	}
	if sp.Delta != nil {
		t.Fatalf("SyncPayload.Delta = %v, want nil", sp.Delta)
	}
}

func TestSyncResultZeroValue(t *testing.T) {
	sr := SyncResult{}
	if sr.RowResults != nil {
		t.Fatalf("SyncResult.RowResults = %v, want nil", sr.RowResults)
	}
}

func TestValidationErrorFormat(t *testing.T) {
	err := ValidationError{
		Row:    "r1",
		Field:  "f1",
		Value:  "v1",
		Reason: "wrong type",
	}
	want := `cell r1:f1 value "v1": wrong type`
	if got := err.Error(); got != want {
		t.Fatalf("ValidationError.Error() = %q, want %q", got, want)
	}
}
