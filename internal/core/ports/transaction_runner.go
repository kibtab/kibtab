// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package. Release v0.7.0 hands the transaction
// context to the function that it runs.
package ports

import "context"

// TransactionRunner runs a group of writes as one transaction.
type TransactionRunner interface {
	// Run runs fn inside one transaction. It commits on success and it
	// rolls back on error. The context that fn receives holds the
	// transaction. Each port call of fn runs on that transaction.
	Run(fn func(ctx context.Context) error) error
}
