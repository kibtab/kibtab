// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package. Release v0.7.0 adds the version lock
// and the version bump.
package ports

import (
	"context"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// TableRegistry lists the tables and reads the metadata of one table.
// It owns the version of each table.
type TableRegistry interface {
	// List returns the metadata for every known table.
	List(ctx context.Context) ([]domain.TableMetadata, error)

	// Metadata returns the metadata for one named table. A table with no
	// stored version row sits at version 1.
	Metadata(ctx context.Context, table string) (domain.TableMetadata, error)

	// LockVersion returns the stored version of one table. It stores
	// version 1 when the table has no version row. It holds the version
	// row until the transaction of the caller ends, so a second write
	// reads the new version after the first write commits.
	LockVersion(ctx context.Context, table string) (int64, error)

	// BumpVersion advances the stored version of one table by one. It
	// returns the new version.
	BumpVersion(ctx context.Context, table string) (int64, error)
}
