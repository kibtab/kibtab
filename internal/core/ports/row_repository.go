// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package.
package ports

import "github.com/kibtab/kibtab/internal/core/domain"

// RowRepository reads and writes rows in one table.
type RowRepository interface {
	// Read returns the rows for one table in key order.
	Read(table string) ([]domain.CellValue, error)

	// Write writes one row in the table.
	Write(table string, row string, values []domain.CellValue) error
}
