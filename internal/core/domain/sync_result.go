// Package domain holds the pure Go models of Kibtab.
//
// The models name no database engine and no spreadsheet client.
// Release v0.2.0 fills the package.
package domain

// SyncResult names the result that the instance returns after a sync.
type SyncResult struct {
	RowResults []RowResult
}

// RowResult names the result for one row in a sync.
type RowResult struct {
	Row     string
	Version int64
	Error   error
}
