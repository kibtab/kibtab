// Package domain holds the pure Go models of Kibtab.
//
// The models name no database engine and no spreadsheet client.
// Release v0.2.0 fills the package.
package domain

// SyncPayload names a table and the delta that a client sends.
type SyncPayload struct {
	Table string
	Delta []CellDelta
}
