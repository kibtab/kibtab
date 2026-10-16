package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// fakeRows is a RowRepository for the transport tests.
type fakeRows struct {
	cells []domain.CellValue
	err   error
}

func (f fakeRows) Read(string) ([]domain.CellValue, error) {
	return f.cells, f.err
}

func (f fakeRows) Write(string, string, []domain.CellValue) error {
	return f.err
}

// fakeSync is a SyncService for the transport tests. It records the
// payload that the transport sends.
type fakeSync struct {
	result  domain.SyncResult
	err     error
	payload domain.SyncPayload
}

func (f *fakeSync) Sync(payload domain.SyncPayload) (domain.SyncResult, error) {
	f.payload = payload
	return f.result, f.err
}

func TestListRows(t *testing.T) {
	registry := fakeRegistry{meta: domain.TableMetadata{
		Name:   "users",
		Fields: []string{"name"},
	}}
	rows := fakeRows{cells: []domain.CellValue{
		{Field: "k", Value: "1"},
		{Field: "name", Value: "alice"},
		{Field: "k", Value: "2"},
		{Field: "name", Value: "bob"},
	}}
	handler := NewHandler(registry, rows, &fakeSync{}, discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables/users/rows", nil))

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}

	var got rowList
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode the body: %v", err)
	}
	if len(got.Rows) != 2 {
		t.Fatalf("rows: got %d, want 2", len(got.Rows))
	}
	if got.Rows[0].Key != "1" || got.Rows[1].Key != "2" {
		t.Fatalf("keys: got %q and %q", got.Rows[0].Key, got.Rows[1].Key)
	}
	if len(got.Rows[0].Cells) != 1 || got.Rows[0].Cells[0].Value != "alice" {
		t.Fatalf("first row cells: got %+v", got.Rows[0].Cells)
	}
}

func TestListRowsEmpty(t *testing.T) {
	handler := NewHandler(fakeRegistry{}, fakeRows{}, &fakeSync{}, discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables/users/rows", nil))

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	if got, want := recorder.Body.String(), `{"rows":[]}`; got != want {
		t.Fatalf("body: got %q, want %q", got, want)
	}
}

func TestListRowsMetadataError(t *testing.T) {
	handler := NewHandler(fakeRegistry{metadataErr: errors.New("boom")},
		fakeRows{}, &fakeSync{}, discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables/users/rows", nil))

	if recorder.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusInternalServerError)
	}
	checkError(t, recorder, codeInternal, "the table metadata failed")
}

func TestListRowsReadError(t *testing.T) {
	handler := NewHandler(fakeRegistry{}, fakeRows{err: errors.New("boom")},
		&fakeSync{}, discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables/users/rows", nil))

	if recorder.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusInternalServerError)
	}
	checkError(t, recorder, codeInternal, "the row list failed")
}

func TestListRowsWriteError(t *testing.T) {
	h := handlers{
		tables: fakeRegistry{},
		rows:   fakeRows{},
		log:    discardLogger(),
	}
	writer := &failingWriter{header: nethttp.Header{}}

	h.listRows(writer,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables/users/rows", nil))
}

func TestGroupRowsFieldBeforeKey(t *testing.T) {
	cells := []domain.CellValue{{Field: "name", Value: "alice"}}

	rows := groupRows(cells, []string{"name"})
	if len(rows) != 0 {
		t.Fatalf("rows: got %d, want 0", len(rows))
	}
}

func TestGroupRowsNoFields(t *testing.T) {
	cells := []domain.CellValue{{Field: "k", Value: "1"}}

	rows := groupRows(cells, nil)
	if len(rows) != 1 || rows[0].Key != "1" {
		t.Fatalf("rows: got %+v", rows)
	}
	if len(rows[0].Cells) != 0 {
		t.Fatalf("cells: got %+v", rows[0].Cells)
	}
}

// usersTable returns a registry that holds one table with one field.
func usersTable() fakeRegistry {
	return fakeRegistry{meta: domain.TableMetadata{
		Name:   "users",
		Fields: []string{"name"},
	}}
}

// postRow sends one row to the add route.
func postRow(t *testing.T, registry fakeRegistry, syncService *fakeSync, body string) *httptest.ResponseRecorder {
	t.Helper()

	handler := NewHandler(registry, fakeRows{}, syncService, discardLogger())
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(
		nethttp.MethodPost, "/v1/tables/users/rows", strings.NewReader(body)))
	return recorder
}

func TestAddRow(t *testing.T) {
	syncService := &fakeSync{result: domain.SyncResult{
		RowResults: []domain.RowResult{{Row: "1", Version: 3}},
	}}
	body := `{"version":3,"key":"1","cells":[{"field":"name","value":"alice"}]}`

	recorder := postRow(t, usersTable(), syncService, body)

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	if got, want := recorder.Body.String(), `{"rows":[{"key":"1","version":3}]}`; got != want {
		t.Fatalf("body: got %q, want %q", got, want)
	}
	if syncService.payload.Table != "users" {
		t.Fatalf("table: got %q, want %q", syncService.payload.Table, "users")
	}
	if len(syncService.payload.Delta) != 1 {
		t.Fatalf("delta: got %+v", syncService.payload.Delta)
	}
	delta := syncService.payload.Delta[0]
	if delta.Row != "1" || delta.Field != "name" ||
		delta.Value != "alice" || delta.Version != 3 {
		t.Fatalf("delta: got %+v", delta)
	}
}

func TestAddRowEmptyKey(t *testing.T) {
	recorder := postRow(t, fakeRegistry{}, &fakeSync{},
		`{"version":3,"cells":[{"field":"name","value":"alice"}]}`)

	if recorder.Code != nethttp.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusBadRequest)
	}
	checkError(t, recorder, codeInvalidRequest, "the row key is empty")
}

func TestAddRowInvalidJSON(t *testing.T) {
	recorder := postRow(t, fakeRegistry{}, &fakeSync{}, "{")

	if recorder.Code != nethttp.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusBadRequest)
	}
	checkError(t, recorder, codeInvalidRequest, "the JSON body is invalid")
}

func TestAddRowMetadataError(t *testing.T) {
	recorder := postRow(t, fakeRegistry{metadataErr: errors.New("boom")},
		&fakeSync{}, `{"version":3,"key":"1"}`)

	if recorder.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusInternalServerError)
	}
	checkError(t, recorder, codeInternal, "the table metadata failed")
}

func TestAddRowUnknownField(t *testing.T) {
	recorder := postRow(t, usersTable(), &fakeSync{},
		`{"version":3,"key":"1","cells":[{"field":"nope","value":"alice"}]}`)

	if recorder.Code != nethttp.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusBadRequest)
	}
	checkError(t, recorder, codeInvalidRequest,
		`the field "nope" is not in the table`)
}

func TestAddRowSyncError(t *testing.T) {
	invalid := domain.ValidationError{
		Row: "1", Field: "name", Value: "alice", Reason: "the value is bad",
	}

	recorder := postRow(t, usersTable(), &fakeSync{err: invalid},
		`{"version":3,"key":"1"}`)

	if recorder.Code != nethttp.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusBadRequest)
	}
	checkError(t, recorder, codeInvalidRequest, invalid.Error())
}

func TestAddRowSyncInternalError(t *testing.T) {
	recorder := postRow(t, usersTable(), &fakeSync{err: errors.New("boom")},
		`{"version":3,"key":"1"}`)

	if recorder.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusInternalServerError)
	}
	checkError(t, recorder, codeInternal, "the row write failed")
}

func TestAddRowConflict(t *testing.T) {
	conflict := domain.ConflictError{Table: "users", Row: "1", Want: 3, Got: 2}
	syncService := &fakeSync{result: domain.SyncResult{
		RowResults: []domain.RowResult{{Row: "1", Version: 3, Error: conflict}},
	}}

	recorder := postRow(t, usersTable(), syncService, `{"version":3,"key":"1"}`)

	if recorder.Code != nethttp.StatusConflict {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusConflict)
	}
	checkError(t, recorder, codeConflict, conflict.Error())
}

func TestAddRowRowValidationError(t *testing.T) {
	invalid := domain.ValidationError{
		Row: "1", Field: "name", Value: "alice", Reason: "the value is bad",
	}
	syncService := &fakeSync{result: domain.SyncResult{
		RowResults: []domain.RowResult{{Row: "1", Version: 3, Error: invalid}},
	}}

	recorder := postRow(t, usersTable(), syncService, `{"version":3,"key":"1"}`)

	if recorder.Code != nethttp.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusBadRequest)
	}
	checkError(t, recorder, codeInvalidRequest, invalid.Error())
}

func TestAddRowRowInternalError(t *testing.T) {
	syncService := &fakeSync{result: domain.SyncResult{
		RowResults: []domain.RowResult{{Row: "1", Version: 3, Error: errors.New("boom")}},
	}}

	recorder := postRow(t, usersTable(), syncService, `{"version":3,"key":"1"}`)

	if recorder.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusInternalServerError)
	}
	checkError(t, recorder, codeInternal, "the row write failed")
}

// checkError asserts the code and the message of an ErrorResponse body.
func checkError(t *testing.T, recorder *httptest.ResponseRecorder, code, message string) {
	t.Helper()

	var got ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode the body: %v", err)
	}
	if got.Code != code {
		t.Fatalf("code: got %q, want %q", got.Code, code)
	}
	if got.Message != message {
		t.Fatalf("message: got %q, want %q", got.Message, message)
	}
}
