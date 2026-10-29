package http

import (
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// fakeAudit is an AuditReader for the transport tests. It records the
// cell that the transport asks for.
type fakeAudit struct {
	entries []domain.AuditRow
	err     error
	asked   string
}

func (f *fakeAudit) Read(_ context.Context, table string, row string,
	field string) ([]domain.AuditRow, error) {
	f.asked = table + "/" + row + "/" + field
	return f.entries, f.err
}

// auditRouteRequest builds the request of the audit route.
func auditRouteRequest(path string) *nethttp.Request {
	return httptest.NewRequest(nethttp.MethodGet, path, nil)
}

func TestListAudit(t *testing.T) {
	changed := time.Date(2026, time.October, 29, 12, 0, 0, 0, time.UTC)
	reader := &fakeAudit{entries: []domain.AuditRow{{
		Table:     "users",
		Row:       "1",
		Field:     "name",
		User:      "rita",
		OldValue:  "bob",
		NewValue:  "alice",
		ChangedAt: changed,
	}}}
	handler := NewHandler(usersTable(), fakeRows{}, reader, &fakeSync{},
		discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		auditRouteRequest("/v1/tables/users/rows/1/audit?field=name"))

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	if reader.asked != "users/1/name" {
		t.Fatalf("asked: got %q, want %q", reader.asked, "users/1/name")
	}

	var got auditList
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode the body: %v", err)
	}
	if len(got.Audit) != 1 {
		t.Fatalf("audit rows: got %d, want 1", len(got.Audit))
	}
	entry := got.Audit[0]
	if entry.Field != "name" || entry.User != "rita" ||
		entry.OldValue != "bob" || entry.NewValue != "alice" {
		t.Fatalf("entry: got %+v, want the audit row of rita", entry)
	}
	if !entry.ChangedAt.Equal(changed) {
		t.Fatalf("changed_at: got %v, want %v", entry.ChangedAt, changed)
	}
}

func TestListAuditEmptyField(t *testing.T) {
	handler := NewHandler(usersTable(), fakeRows{}, &fakeAudit{},
		&fakeSync{}, discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		auditRouteRequest("/v1/tables/users/rows/1/audit"))

	if recorder.Code != nethttp.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusBadRequest)
	}
	checkError(t, recorder, codeInvalidRequest,
		"the field query parameter is empty")
}

func TestListAuditReadError(t *testing.T) {
	reader := &fakeAudit{err: errors.New("boom")}
	handler := NewHandler(usersTable(), fakeRows{}, reader, &fakeSync{},
		discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		auditRouteRequest("/v1/tables/users/rows/1/audit?field=name"))

	if recorder.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusInternalServerError)
	}
	checkError(t, recorder, codeInternal, "the audit list failed")
}

func TestListAuditEmptyAnswer(t *testing.T) {
	handler := NewHandler(usersTable(), fakeRows{}, &fakeAudit{},
		&fakeSync{}, discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		auditRouteRequest("/v1/tables/users/rows/1/audit?field=name"))

	if got, want := recorder.Body.String(), `{"audit":[]}`; got != want {
		t.Fatalf("body: got %q, want %q", got, want)
	}
}

func TestRequestUser(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{name: "set", header: "rita", want: "rita"},
		{name: "padded", header: "  rita  ", want: "rita"},
		{name: "empty", header: "", want: unknownUser},
		{name: "spaces", header: "   ", want: unknownUser},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(nethttp.MethodPost,
				"/v1/tables/users/cells", nil)
			request.Header.Set(headerUser, test.header)
			if got := requestUser(request); got != test.want {
				t.Fatalf("requestUser: got %q, want %q", got, test.want)
			}
		})
	}
}

func TestAddCellsCarriesTheUserHeader(t *testing.T) {
	syncService := &fakeSync{result: domain.SyncResult{
		RowResults: []domain.RowResult{{Row: "1", Version: 4}},
	}}
	handler := NewHandler(usersTable(), fakeRows{}, &fakeAudit{}, syncService,
		discardLogger())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(nethttp.MethodPost,
		"/v1/tables/users/cells",
		strings.NewReader(`{"delta":[{"row":"1","field":"name",`+
			`"value":"alice","version":3}]}`))
	request.Header.Set(headerUser, "rita")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	if syncService.payload.User != "rita" {
		t.Fatalf("user: got %q, want %q", syncService.payload.User, "rita")
	}
}

func TestAddCellsDefaultsToTheUnknownUser(t *testing.T) {
	syncService := &fakeSync{result: domain.SyncResult{
		RowResults: []domain.RowResult{{Row: "1", Version: 4}},
	}}
	handler := NewHandler(usersTable(), fakeRows{}, &fakeAudit{}, syncService,
		discardLogger())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(nethttp.MethodPost,
		"/v1/tables/users/cells",
		strings.NewReader(`{"delta":[{"row":"1","field":"name",`+
			`"value":"alice","version":3}]}`))
	handler.ServeHTTP(recorder, request)

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	if syncService.payload.User != unknownUser {
		t.Fatalf("user: got %q, want %q", syncService.payload.User,
			unknownUser)
	}
}
