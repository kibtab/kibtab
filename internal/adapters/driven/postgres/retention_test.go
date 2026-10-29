package postgres

import (
	"context"
	"testing"
)

func TestParseRetentionDays(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{name: "default", value: "", want: defaultAuditRetentionDays},
		{name: "set", value: "30", want: 30},
		{name: "zero", value: "0", want: 0},
		{name: "negative", value: "-1", wantErr: true},
		{name: "not a number", value: "ninety", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			getenv := func(string) string { return test.value }
			got, err := parseRetentionDays(getenv)
			if test.wantErr {
				if err == nil {
					t.Fatalf("value %q: got %d, want an error",
						test.value, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("value %q: unexpected error: %v", test.value, err)
			}
			if got != test.want {
				t.Fatalf("value %q: got %d, want %d", test.value, got,
					test.want)
			}
		})
	}
}

// TestIntegrationPruneAudit removes one old audit row and it keeps the
// row inside the retention days.
func TestIntegrationPruneAudit(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)
	ctx := context.Background()
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	const insert = `INSERT INTO _kibtab_meta.audit_log
(table_name, row_key, field_name, changed_by, old_value, new_value, changed_at)
VALUES ('items', '1', 'name', 'rita', '', 'alice',
        now() - make_interval(days => $1))`
	if _, err := pool.Exec(ctx, insert, 10); err != nil {
		t.Fatalf("insert the old audit row: %v", err)
	}
	if _, err := pool.Exec(ctx, insert, 1); err != nil {
		t.Fatalf("insert the recent audit row: %v", err)
	}

	engine := &Engine{pool: pool, retentionDays: 5}
	removed, err := engine.PruneAudit(ctx)
	if err != nil {
		t.Fatalf("PruneAudit: %v", err)
	}
	if removed != 1 {
		t.Fatalf("removed rows: got %d, want 1", removed)
	}

	var left int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM _kibtab_meta.audit_log").Scan(&left); err != nil {
		t.Fatalf("count the audit rows: %v", err)
	}
	if left != 1 {
		t.Fatalf("left rows: got %d, want 1", left)
	}
}

func TestPruneAuditError(t *testing.T) {
	if !available(t) {
		t.Skip("no PostgreSQL available")
	}

	pool, cleanup := testSetup(t)
	t.Cleanup(cleanup)
	engine := &Engine{pool: pool, retentionDays: 5}
	pool.Close()

	if _, err := engine.PruneAudit(context.Background()); err == nil {
		t.Fatal("PruneAudit on a closed pool: got nil, want error")
	}
}
