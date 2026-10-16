// Package services holds the use cases of the core.
//
// A service calls the ports only. Release v0.2.0 fills the package.
package services

import (
	"sort"

	"github.com/kibtab/kibtab/internal/core/domain"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// SyncService writes a sync payload to the ports.
// It groups the delta by row. It writes each group in one transaction.
type SyncService struct {
	registry   ports.TableRegistry
	repository ports.RowRepository
	dialect    ports.Dialect
	runner     ports.TransactionRunner
}

// Compile-time proof that SyncService satisfies the SyncService port.
var _ ports.SyncService = (*SyncService)(nil)

// NewSyncService builds a SyncService from the ports that it calls.
func NewSyncService(
	registry ports.TableRegistry,
	repository ports.RowRepository,
	dialect ports.Dialect,
	runner ports.TransactionRunner,
) *SyncService {
	return &SyncService{
		registry:   registry,
		repository: repository,
		dialect:    dialect,
		runner:     runner,
	}
}

// Sync validates the payload. It groups the delta by row. It writes each
// group in one transaction. It returns one row result for each group.
// A failed group rolls back on its own. It stops no other group.
func (s *SyncService) Sync(payload domain.SyncPayload) (domain.SyncResult, error) {
	if err := ValidatePayload(payload); err != nil {
		return domain.SyncResult{}, err
	}

	groups := GroupByRow(payload.Delta)
	versions := MergeByRow(payload.Delta)
	rows := rowKeys(groups)

	result := domain.SyncResult{
		RowResults: make([]domain.RowResult, 0, len(rows)),
	}
	for _, row := range rows {
		result.RowResults = append(result.RowResults,
			s.writeRow(payload.Table, row, versions[row], groups[row]))
	}
	return result, nil
}

// rowKeys returns the row keys in order.
func rowKeys(groups map[string][]domain.CellDelta) []string {
	rows := make([]string, 0, len(groups))
	for row := range groups {
		rows = append(rows, row)
	}
	sort.Strings(rows)
	return rows
}

// writeRow writes one row group in one transaction. It compares the payload
// version to the table version before the write. It returns a conflict when
// the two versions differ.
func (s *SyncService) writeRow(
	table string,
	row string,
	version int64,
	deltas []domain.CellDelta,
) domain.RowResult {
	err := s.runner.Run(func() error {
		meta, err := s.registry.Metadata(table)
		if err != nil {
			return err
		}
		if meta.Version != version {
			return domain.ConflictError{
				Table: table,
				Row:   row,
				Want:  version,
				Got:   meta.Version,
			}
		}
		values, err := s.coerce(deltas)
		if err != nil {
			return err
		}
		return s.repository.Write(table, row, values)
	})
	return domain.RowResult{Row: row, Version: version, Error: err}
}

// coerce maps each delta value to the engine form for its field.
func (s *SyncService) coerce(deltas []domain.CellDelta) ([]domain.CellValue, error) {
	values := make([]domain.CellValue, 0, len(deltas))
	for _, delta := range deltas {
		value, err := s.dialect.MapValue(delta.Field, delta.Value)
		if err != nil {
			return nil, err
		}
		values = append(values, domain.CellValue{
			Field: delta.Field,
			Value: value,
		})
	}
	return values, nil
}
