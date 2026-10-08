// Package services holds the use cases of the core.
//
// A service calls the ports only. Release v0.2.0 fills the package.
package services

import "github.com/kibtab/kibtab/internal/core/domain"

// MergeByRow groups deltas by row and keeps the latest version per row.
func MergeByRow(deltas []domain.CellDelta) map[string]int64 {
	latest := make(map[string]int64)
	for _, delta := range deltas {
		v := delta.Version
		if current, ok := latest[delta.Row]; !ok || v > current {
			latest[delta.Row] = v
		}
	}
	return latest
}
