// Package domain holds the pure Go models of Kibtab.
//
// The models name no database engine and no spreadsheet client.
// Release v0.7.0 adds the audit row.
package domain

import "time"

// AuditRow names one audit row. It records a change to one cell.
// A write fills the table, the row, the field, the user, and the two
// values. A read fills the time of the change as well.
type AuditRow struct {
	Table     string
	Row       string
	Field     string
	User      string
	OldValue  string
	NewValue  string
	ChangedAt time.Time
}
