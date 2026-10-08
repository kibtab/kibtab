package services

import (
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

func TestValidatePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload domain.SyncPayload
		wantErr bool
		want    string
	}{
		{
			name: "empty table",
			payload: domain.SyncPayload{
				Table: "",
				Delta: []domain.CellDelta{{
					Row:    "r1",
					Field:  "f1",
					Value:  "v1",
					Version: 1,
				}},
			},
			wantErr: true,
			want:    "cell : value \"\": table name is empty",
		},
		{
			name: "empty delta",
			payload: domain.SyncPayload{
				Table: "t1",
				Delta: []domain.CellDelta{},
			},
			wantErr: true,
			want:    "cell : value \"\": delta is empty",
		},
		{
			name: "invalid row",
			payload: domain.SyncPayload{
				Table: "t1",
				Delta: []domain.CellDelta{{
					Row:    "",
					Field:  "f1",
					Value:  "v1",
					Version: 1,
				}},
			},
			wantErr: true,
			want:    "cell :f1 value \"v1\": sync payload delta 0 has an invalid row",
		},
		{
			name: "invalid field",
			payload: domain.SyncPayload{
				Table: "t1",
				Delta: []domain.CellDelta{{
					Row:    "r1",
					Field:  "",
					Value:  "v1",
					Version: 1,
				}},
			},
			wantErr: true,
			want:    "cell r1: value \"v1\": sync payload delta 0 has an invalid field",
		},
		{
			name: "invalid value",
			payload: domain.SyncPayload{
				Table: "t1",
				Delta: []domain.CellDelta{{
					Row:    "r1",
					Field:  "f1",
					Value:  "",
					Version: 1,
				}},
			},
			wantErr: true,
			want:    "cell r1:f1 value \"\": sync payload delta 0 has an invalid value",
		},
		{
			name: "invalid version",
			payload: domain.SyncPayload{
				Table: "t1",
				Delta: []domain.CellDelta{{
					Row:    "r1",
					Field:  "f1",
					Value:  "v1",
					Version: 0,
				}},
			},
			wantErr: true,
			want:    "cell r1:f1 value \"v1\": sync payload delta 0 has an invalid version",
		},
		{
			name: "invalid negative version",
			payload: domain.SyncPayload{
				Table: "t1",
				Delta: []domain.CellDelta{{
					Row:    "r1",
					Field:  "f1",
					Value:  "v1",
					Version: -1,
				}},
			},
			wantErr: true,
			want:    "cell r1:f1 value \"v1\": sync payload delta 0 has an invalid version",
		},
		{
			name: "valid payload",
			payload: domain.SyncPayload{
				Table: "t1",
				Delta: []domain.CellDelta{{
					Row:    "r1",
					Field:  "f1",
					Value:  "v1",
					Version: 1,
				}},
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidatePayload(test.payload)
			if test.wantErr {
				if err == nil {
					t.Fatalf("got no error, want one")
				}
			if !contains(err.Error(), test.want) {
				t.Fatalf("error: got %q, want a message like %q", err.Error(), test.want)
			}
				return
			}
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
