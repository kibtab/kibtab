// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.7.0 adds the AuditReader port.
package ports

import (
	"context"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// AuditReader reads the audit rows of one cell.
type AuditReader interface {
	// Read returns the audit rows for one cell in write order.
	Read(ctx context.Context, table string, row string,
		field string) ([]domain.AuditRow, error)
}
