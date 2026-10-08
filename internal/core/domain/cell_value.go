// Package domain holds the pure Go models of Kibtab.
//
// The models name no database engine and no spreadsheet client.
// Release v0.2.0 fills the package.
package domain

// CellValue names one stored value in a row.
type CellValue struct {
	Field string
	Value string
}
