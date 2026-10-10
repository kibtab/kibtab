package postgres

import (
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// runRowRepositoryContract runs the shared contract suite for the
// RowRepository port. It is called by each engine adapter's test.
func runRowRepositoryContract(t *testing.T, repo *RowRepository, tableName string, fields []string) {
	t.Helper()

	t.Run("write and read one row", func(t *testing.T) {
		truncateTable(t, repo.pool, tableName)
		row := "r1"
		values := []domain.CellValue{
			{Field: "name", Value: "alice"},
			{Field: "email", Value: "alice@example.org"},
		}
		if err := repo.Write(tableName, row, values); err != nil {
			t.Fatalf("Write: %v", err)
		}
		cells, err := repo.Read(tableName)
		if err != nil {
			t.Fatalf("Read: %v", err)
		}
		if len(cells) == 0 {
			t.Fatal("Read returned no cells")
		}
		// First cell should be the primary key.
		if cells[0].Field != pkColumn {
			t.Fatalf("first cell field: got %q, want %q", cells[0].Field, pkColumn)
		}
		if cells[0].Value != row {
			t.Fatalf("primary key value: got %q, want %q", cells[0].Value, row)
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
		if !foundName {
			t.Fatal("name cell not found")
		}
		if !foundEmail {
			t.Fatal("email cell not found")
		}
	})

	t.Run("write and read two rows", func(t *testing.T) {
		truncateTable(t, repo.pool, tableName)
		if err := repo.Write(tableName, "r1", []domain.CellValue{
			{Field: "name", Value: "alice"},
		}); err != nil {
			t.Fatalf("Write r1: %v", err)
		}
		if err := repo.Write(tableName, "r2", []domain.CellValue{
			{Field: "name", Value: "bob"},
		}); err != nil {
			t.Fatalf("Write r2: %v", err)
		}
		cells, err := repo.Read(tableName)
		if err != nil {
			t.Fatalf("Read: %v", err)
		}
		// Should have at least 4 cells: 2 keys + 2 field values.
		if len(cells) < 4 {
			t.Fatalf("Read returned %d cells, want at least 4", len(cells))
		}
		keys := map[string]bool{}
		for i := 0; i < len(cells); i += 1 + len(fields) {
			if cells[i].Field == pkColumn {
				keys[cells[i].Value] = true
			}
		}
		if !keys["r1"] {
			t.Fatal("row r1 not found")
		}
		if !keys["r2"] {
			t.Fatal("row r2 not found")
		}
	})

	t.Run("write updates existing row", func(t *testing.T) {
		truncateTable(t, repo.pool, tableName)
		if err := repo.Write(tableName, "r1", []domain.CellValue{
			{Field: "name", Value: "alice"},
		}); err != nil {
			t.Fatalf("Write v1: %v", err)
		}
		if err := repo.Write(tableName, "r1", []domain.CellValue{
			{Field: "name", Value: "alice2"},
		}); err != nil {
			t.Fatalf("Write v2: %v", err)
		}
		cells, err := repo.Read(tableName)
		if err != nil {
			t.Fatalf("Read: %v", err)
		}
		for _, c := range cells {
			if c.Field == "name" && c.Value == "alice2" {
				return
			}
		}
		t.Fatal("updated value not found")
	})

	t.Run("write with empty values is a no-op", func(t *testing.T) {
		truncateTable(t, repo.pool, tableName)
		err := repo.Write(tableName, "r1", nil)
		if err != nil {
			t.Fatalf("Write with nil values: %v", err)
		}
	})

	t.Run("write with NUL byte in value returns error", func(t *testing.T) {
		truncateTable(t, repo.pool, tableName)
		err := repo.Write(tableName, "r1", []domain.CellValue{
			{Field: "name", Value: "bad\x00value"},
		})
		if err == nil {
			t.Fatal("Write with NUL byte: got no error, want one")
		}
	})

	t.Run("write with empty field name returns error", func(t *testing.T) {
		truncateTable(t, repo.pool, tableName)
		err := repo.Write(tableName, "r1", []domain.CellValue{
			{Field: "", Value: "v"},
		})
		if err == nil {
			t.Fatal("Write with empty field: got no error, want one")
		}
	})

	t.Run("read from non-existent table returns nil", func(t *testing.T) {
		cells, err := repo.Read("does_not_exist_999")
		if err != nil {
			t.Fatalf("Read non-existent table: %v", err)
		}
		if cells != nil {
			t.Fatalf("Read non-existent table: got %v, want nil", cells)
		}
	})
}
