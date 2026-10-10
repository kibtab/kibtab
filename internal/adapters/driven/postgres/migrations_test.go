package postgres

import "testing"

func TestSplitStatements(t *testing.T) {
	statements := splitStatements(migrationStmts)
	if len(statements) == 0 {
		t.Fatal("expected at least one statement, got zero")
	}
	for _, stmt := range statements {
		if stmt == "" {
			t.Fatal("expected non-empty statements only")
		}
	}
}

func TestMigrationStmtsContainsSchema(t *testing.T) {
	if !contains(migrationStmts, "_kibtab_meta") {
		t.Fatal("migration SQL must create schema _kibtab_meta")
	}
	if !contains(migrationStmts, "table_versions") {
		t.Fatal("migration SQL must create table _kibtab_meta.table_versions")
	}
	if !contains(migrationStmts, "audit_log") {
		t.Fatal("migration SQL must create table _kibtab_meta.audit_log")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
