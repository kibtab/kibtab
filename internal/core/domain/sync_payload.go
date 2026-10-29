// Package domain holds the pure Go models of Kibtab.
//
// The models name no database engine and no spreadsheet client.
// Release v0.2.0 fills the package. Release v0.7.0 adds the user of a
// sync.
package domain

// SyncPayload names a table and the delta that a client sends.
// The user names the person that the write serves. The transport reads
// it from the request, and it reaches each audit row of the sync.
type SyncPayload struct {
	Table string
	User  string
	Delta []CellDelta
}
