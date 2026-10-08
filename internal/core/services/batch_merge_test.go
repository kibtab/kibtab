package services

import (
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

func TestMergeByRow(t *testing.T) {
	tests := []struct {
		name   string
		deltas []domain.CellDelta
		want   map[string]int64
	}{
		{
			name:   "empty",
			deltas: []domain.CellDelta{},
			want:   map[string]int64{},
		},
		{
			name: "single row",
			deltas: []domain.CellDelta{{
				Row:    "r1",
				Field:  "f1",
				Value:  "v1",
				Version: 1,
			}},
			want: map[string]int64{"r1": 1},
		},
		{
			name: "multiple rows",
			deltas: []domain.CellDelta{{
				Row:    "r1",
				Field:  "f1",
				Value:  "v1",
				Version: 1,
			}, {
				Row:    "r2",
				Field:  "f1",
				Value:  "v1",
				Version: 2,
			}},
			want: map[string]int64{"r1": 1, "r2": 2},
		},
		{
			name: "latest version wins",
			deltas: []domain.CellDelta{{
				Row:    "r1",
				Field:  "f1",
				Value:  "v1",
				Version: 1,
			}, {
				Row:    "r1",
				Field:  "f2",
				Value:  "v2",
				Version: 3,
			}},
			want: map[string]int64{"r1": 3},
		},
		{
			name: "equal versions keep first in map logic but latest wins on tie",
			deltas: []domain.CellDelta{{
				Row:    "r1",
				Field:  "f1",
				Value:  "v1",
				Version: 2,
			}, {
				Row:    "r1",
				Field:  "f2",
				Value:  "v2",
				Version: 2,
			}},
			want: map[string]int64{"r1": 2},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := MergeByRow(test.deltas)
			if len(got) != len(test.want) {
				t.Fatalf("MergeByRow() = %v, want %v", got, test.want)
			}
			for row, wantVersion := range test.want {
				if gotVersion, ok := got[row]; !ok {
					t.Fatalf("MergeByRow() missing row %q, want %v", row, test.want)
				} else if gotVersion != wantVersion {
					t.Fatalf("MergeByRow() row %q = %d, want %d", row, gotVersion, wantVersion)
				}
			}
		})
	}
}
