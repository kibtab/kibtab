// Package services holds the use cases of the core.
//
// A service calls the ports only. Release v0.2.0 fills the package.
// Release v0.7.0 writes one audit row per cell change and it bumps the
// table version in the same transaction as the data.
package services

import (
	"context"
	"sort"

	"github.com/kibtab/kibtab/internal/core/domain"
	"github.com/kibtab/kibtab/internal/core/ports"
)

// SyncService writes a sync payload to the ports.
// It groups the delta by row. It writes each group in one transaction.
// The transaction holds the version lock, the data write, the audit
// rows, and the version bump.
type SyncService struct {
	registry   ports.TableRegistry
	repository ports.RowRepository
	audit      ports.AuditWriter
	dialect    ports.Dialect
	runner     ports.TransactionRunner
}

// Compile-time proof that SyncService satisfies the SyncService port.
var _ ports.SyncService = (*SyncService)(nil)

// NewSyncService builds a SyncService from the ports that it calls.
func NewSyncService(
	registry ports.TableRegistry,
	repository ports.RowRepository,
	audit ports.AuditWriter,
	dialect ports.Dialect,
	runner ports.TransactionRunner,
) *SyncService {
	return &SyncService{
		registry:   registry,
		repository: repository,
		audit:      audit,
		dialect:    dialect,
		runner:     runner,
	}
}

// Sync validates the payload. It groups the delta by row. It writes each
// group in one transaction. It returns one row result for each group.
// A failed group rolls back on its own. It stops no other group.
//
// The version of the table advances with each written group. The first
// group checks the version of the payload. Each later group checks the
// version that the last written group produced. A group that fails keeps
// the expected version for the next group.
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
	expected := versions[rows[0]]
	for _, row := range rows {
		rowResult := s.writeRow(payload.Table, row, expected,
			payload.User, groups[row])
		result.RowResults = append(result.RowResults, rowResult)
		if rowResult.Error == nil {
			expected = rowResult.Version
		}
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

// writeRow writes one row group in one transaction. It locks the version
// row before it compares the expected version to the stored version. It
// writes the data, one audit row for each cell change, and the new
// version in that one transaction. A stale version returns a conflict
// and it writes nothing. The row result holds the new version.
func (s *SyncService) writeRow(
	table string,
	row string,
	version int64,
	user string,
	deltas []domain.CellDelta,
) domain.RowResult {
	newVersion := version
	err := s.runner.Run(func(ctx context.Context) error {
		stored, err := s.registry.LockVersion(ctx, table)
		if err != nil {
			return err
		}
		if stored != version {
			return domain.ConflictError{
				Table: table,
				Row:   row,
				Want:  version,
				Got:   stored,
			}
		}
		values, err := s.coerce(deltas)
		if err != nil {
			return err
		}
		old, err := s.repository.ReadRow(ctx, table, row)
		if err != nil {
			return err
		}
		if err := s.repository.Write(ctx, table, row, values); err != nil {
			return err
		}
		if err := s.writeAudit(ctx, table, row, user, deltas, values,
			old); err != nil {
			return err
		}
		newVersion, err = s.registry.BumpVersion(ctx, table)
		return err
	})
	if err != nil {
		return domain.RowResult{Row: row, Version: version, Error: err}
	}
	return domain.RowResult{Row: row, Version: newVersion}
}

// writeAudit writes one audit row for each cell change of one group.
func (s *SyncService) writeAudit(
	ctx context.Context,
	table string,
	row string,
	user string,
	deltas []domain.CellDelta,
	values []domain.CellValue,
	old []domain.CellValue,
) error {
	for i, delta := range deltas {
		if err := s.audit.Record(ctx, domain.AuditRow{
			Table:    table,
			Row:      row,
			Field:    delta.Field,
			User:     user,
			OldValue: storedValue(old, delta.Field),
			NewValue: values[i].Value,
		}); err != nil {
			return err
		}
	}
	return nil
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

// storedValue returns the stored value of one field. It returns an empty
// text when the row holds no value for the field.
func storedValue(cells []domain.CellValue, field string) string {
	for _, cell := range cells {
		if cell.Field == field {
			return cell.Value
		}
	}
	return ""
}
