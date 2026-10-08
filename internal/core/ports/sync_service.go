// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package.
package ports

import "github.com/kibtab/kibtab/internal/core/domain"

// SyncService accepts a delta and returns a result.
type SyncService interface {
	// Sync accepts one sync payload and returns the result.
	Sync(payload domain.SyncPayload) (domain.SyncResult, error)
}
