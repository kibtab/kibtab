package contracttest

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// SyncRecorder is a SyncService for the SyncService suite. It records the
// payload of the last call and it answers with the fixture result.
type SyncRecorder struct {
	Payload domain.SyncPayload
	Result  domain.SyncResult
	Err     error
}

// Sync records the payload and it answers with the fixture result.
func (r *SyncRecorder) Sync(payload domain.SyncPayload) (domain.SyncResult, error) {
	r.Payload = payload
	return r.Result, r.Err
}

// reset clears the recorder for one case of the suite.
func (r *SyncRecorder) reset() {
	r.Payload = domain.SyncPayload{}
	r.Result = domain.SyncResult{}
	r.Err = nil
}

// SyncFixture holds what one client gives to the SyncService suite.
type SyncFixture struct {
	// Service is the port under test. The client wires it as the
	// SyncService of its own request path.
	Service *SyncRecorder

	// Send passes one payload through the client under test.
	Send func(domain.SyncPayload) (domain.SyncResult, error)

	// IsConflict reports that err tells the caller of a stale version.
	IsConflict func(error) bool
}

// SyncService runs the shared contract suite for the SyncService port.
// Each client calls it with its own fixture. It returns one error for
// each failed case.
func SyncService(f SyncFixture) []error {
	all := []error{}
	all = append(all, payloadPassesIntact(f)...)
	all = append(all, resultMapsBack(f)...)
	all = append(all, staleWriteReturnsConflict(f)...)
	all = append(all, rejectedValueReturnsReason(f)...)
	all = append(all, portFailureReturnsError(f)...)
	return dropNil(all)
}

// suitePayload returns the payload that each case sends.
func suitePayload() domain.SyncPayload {
	return domain.SyncPayload{
		Table: "users",
		Delta: []domain.CellDelta{
			{Row: "1", Field: "name", Value: "alice", Version: 3},
			{Row: "2", Field: "name", Value: "bob", Version: 3},
		},
	}
}

// payloadPassesIntact checks that the port receives the payload of the
// client with its table, its delta, and the order of the delta.
func payloadPassesIntact(f SyncFixture) []error {
	const label = "payload passes through the client"
	payload := suitePayload()
	f.Service.reset()
	f.Service.Result = domain.SyncResult{RowResults: []domain.RowResult{
		{Row: "1", Version: 3},
		{Row: "2", Version: 3},
	}}

	_, err := f.Send(payload)
	errs := []error{wrap(label+": Send", err)}
	errs = append(errs,
		check(label+": table", f.Service.Payload.Table, payload.Table),
		check(label+": delta", formatDelta(f.Service.Payload.Delta),
			formatDelta(payload.Delta)))
	return errs
}

// formatDelta writes each cell change as one line of text. A caller
// compares the text of the port payload to the text of the payload that
// the client means to send.
func formatDelta(delta []domain.CellDelta) string {
	parts := make([]string, 0, len(delta))
	for _, cell := range delta {
		parts = append(parts, fmt.Sprintf("%s/%s/%s/%d",
			cell.Row, cell.Field, cell.Value, cell.Version))
	}
	return strings.Join(parts, ",")
}

// resultMapsBack checks that the caller receives the row and the version
// of the port answer.
func resultMapsBack(f SyncFixture) []error {
	const label = "result maps back to the caller"
	f.Service.reset()
	f.Service.Result = domain.SyncResult{RowResults: []domain.RowResult{
		{Row: "1", Version: 7},
	}}

	result, err := f.Send(suitePayload())
	errs := []error{wrap(label+": Send", err)}
	errs = append(errs, checkMin(label+": row count", len(result.RowResults), 1))
	if len(result.RowResults) > 0 {
		errs = append(errs,
			check(label+": row key", result.RowResults[0].Row, "1"),
			checkVersion(label+": row version", result.RowResults[0].Version, 7))
	}
	return errs
}

// staleWriteReturnsConflict checks that a stale version reaches the
// caller as a conflict.
func staleWriteReturnsConflict(f SyncFixture) []error {
	const label = "stale write returns a conflict"
	f.Service.reset()
	f.Service.Result = domain.SyncResult{RowResults: []domain.RowResult{
		{Row: "1", Error: domain.ConflictError{
			Table: "users", Row: "1", Want: 3, Got: 2,
		}},
	}}

	_, err := f.Send(suitePayload())
	errs := []error{checkError(label+": Send", err)}
	errs = append(errs,
		checkBool(label+": conflict detected", err != nil && f.IsConflict(err)))
	return errs
}

// rejectedValueReturnsReason checks that the reason of a rejected value
// reaches the caller.
func rejectedValueReturnsReason(f SyncFixture) []error {
	const label = "rejected value returns the reason"
	f.Service.reset()
	invalid := domain.ValidationError{
		Row: "1", Field: "name", Value: "v", Reason: "the value is bad",
	}
	f.Service.Result = domain.SyncResult{RowResults: []domain.RowResult{
		{Row: "1", Error: invalid},
	}}

	_, err := f.Send(suitePayload())
	errs := []error{checkError(label+": Send", err)}
	errs = append(errs, checkBool(label+": reason",
		err != nil && strings.Contains(err.Error(), invalid.Error())))
	return errs
}

// portFailureReturnsError checks that a failure of the port reaches the
// caller as an error.
func portFailureReturnsError(f SyncFixture) []error {
	const label = "port failure returns an error"
	f.Service.reset()
	f.Service.Err = errors.New("the port failed")

	_, err := f.Send(suitePayload())
	return []error{checkError(label, err)}
}

// checkVersion returns an error when got differs from want.
func checkVersion(label string, got int64, want int64) error {
	if got == want {
		return nil
	}
	return fmt.Errorf("%s: got %d, want %d", label, got, want)
}
