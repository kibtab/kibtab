// Package domain holds the pure Go models of Kibtab.
//
// The models name no database engine and no spreadsheet client.
// Release v0.2.0 fills the package. Release v0.4.0 adds the conflict error.
package domain

import "fmt"

// ConflictError reports a write that used a stale version.
type ConflictError struct {
	Table string
	Row   string
	Want  int64
	Got   int64
}

func (e ConflictError) Error() string {
	return fmt.Sprintf(
		"table %s row %s: version conflict, got %d, want %d",
		e.Table, e.Row, e.Got, e.Want,
	)
}
