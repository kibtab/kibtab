// Package ports holds the interfaces of the core.
//
// An adapter implements each port. The core owns the interfaces.
// Release v0.2.0 fills the package. Release v0.7.0 puts the context of the
// transaction in each call.
package ports

import (
	"context"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// RowRepository reads and writes rows in one table.
// A call that runs inside a transaction carries that transaction in its
// context.
type RowRepository interface {
	// Read returns the rows for one table in key order.
	Read(ctx context.Context, table string) ([]domain.CellValue, error)

	// ReadRow returns the cells of one row without the row key column.
	// It returns no cell when the row does not exist.
	ReadRow(ctx context.Context, table string,
		row string) ([]domain.CellValue, error)

	// Write writes one row in the table.
	Write(ctx context.Context, table string, row string,
		values []domain.CellValue) error
}
