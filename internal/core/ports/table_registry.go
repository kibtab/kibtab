// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package.
package ports

import "github.com/kibtab/kibtab/internal/core/domain"

// TableRegistry lists the tables and the metadata for each table.
type TableRegistry interface {
	// List returns the metadata for every known table.
	List() ([]domain.TableMetadata, error)

	// Metadata returns the metadata for one named table.
	Metadata(table string) (domain.TableMetadata, error)
}
