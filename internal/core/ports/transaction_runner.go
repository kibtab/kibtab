// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package.
package ports

// TransactionRunner runs a group of writes as one transaction.
type TransactionRunner interface {
	// Run runs fn inside one transaction. It commits on success and rolls
	// back on error.
	Run(fn func() error) error
}
