package postgres

import (
	"context"
	"fmt"
	"strconv"
)

// EnvKeyAuditRetentionDays is the environment variable that holds the
// days that an audit row stays in the database.
const EnvKeyAuditRetentionDays = "KIBTAB_AUDIT_RETENTION_DAYS"

// defaultAuditRetentionDays is the retention when the environment holds
// no value for EnvKeyAuditRetentionDays.
const defaultAuditRetentionDays = 90

// parseRetentionDays reads the retention days from the environment. It
// returns the default when the variable is empty. A value that is no
// number or a value below zero returns an error.
func parseRetentionDays(getenv func(string) string) (int, error) {
	raw := getenv(EnvKeyAuditRetentionDays)
	if raw == "" {
		return defaultAuditRetentionDays, nil
	}
	days, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", EnvKeyAuditRetentionDays, err)
	}
	if days < 0 {
		return 0, fmt.Errorf("%s must not be negative, got %d",
			EnvKeyAuditRetentionDays, days)
	}
	return days, nil
}

// PruneAudit removes the audit rows that are older than the retention
// days of the engine. It returns the number of removed rows.
func (e *Engine) PruneAudit(ctx context.Context) (int64, error) {
	const query = `DELETE FROM _kibtab_meta.audit_log
WHERE changed_at < now() - make_interval(days => $1)`
	tag, err := e.pool.Exec(ctx, query, e.retentionDays)
	if err != nil {
		return 0, fmt.Errorf("prune audit rows: %w", err)
	}
	return tag.RowsAffected(), nil
}
