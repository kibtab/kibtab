// Package domain holds the pure Go models of Kibtab.
//
// The models name no database engine and no spreadsheet client.
// Release v0.2.0 fills the package.
package domain

// TableMetadata names one table and its current version.
type TableMetadata struct {
	Name    string
	Fields  []string
	Version int64
}
