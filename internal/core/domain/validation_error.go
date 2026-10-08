// Package domain holds the pure Go models of Kibtab.
//
// The models name no database engine and no spreadsheet client.
// Release v0.2.0 fills the package.
package domain

import "fmt"

// ValidationError reports a rejected cell value.
type ValidationError struct {
	Row    string
	Field  string
	Value  string
	Reason string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf(
		"cell %s:%s value %q: %s",
		e.Row, e.Field, e.Value, e.Reason,
	)
}
