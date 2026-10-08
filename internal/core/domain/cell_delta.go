// Package domain holds the pure Go models of Kibtab.
//
// The models name no database engine and no spreadsheet client.
// Release v0.2.0 fills the package.
package domain

// CellDelta names one cell change in a sync payload.
type CellDelta struct {
	Row    string
	Field  string
	Value  string
	Version int64
}
