// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package. Release v0.7.0 records one audit row
// with its user.
package ports

import (
	"context"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// AuditWriter records a change to a cell.
type AuditWriter interface {
	// Record stores one audit row. It runs on the transaction that the
	// context holds.
	Record(ctx context.Context, entry domain.AuditRow) error
}
