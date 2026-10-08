// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package.
package ports

// AuditWriter records a change to a cell.
type AuditWriter interface {
	// Write records one audit row.
	Write(table string, row string, field string, oldValue string, newValue string) error
}
