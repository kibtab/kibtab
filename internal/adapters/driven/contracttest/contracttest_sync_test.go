// Package contracttest_test verifies the contract test helpers work as expected.
package contracttest

import (
	"errors"
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

func TestSyncServiceSuite(t *testing.T) {
	recorder := &SyncRecorder{}

	fixture := SyncFixture{
		Service: recorder,
		Send: func(payload domain.SyncPayload) (domain.SyncResult, error) {
			result, err := recorder.Sync(payload)
			// Check for row errors in the result, like a real client would.
			for _, row := range result.RowResults {
				if row.Error != nil {
					return result, row.Error
				}
			}
			return result, err
		},
		IsConflict: func(err error) bool {
			var conflict domain.ConflictError
			return errors.As(err, &conflict)
		},
	}

	errs := SyncService(fixture)
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}