package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
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
	return nil
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
	handler := NewHandler(registry, rows, discardLogger())

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
	handler := NewHandler(fakeRegistry{}, fakeRows{}, discardLogger())

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
		fakeRows{}, discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables/users/rows", nil))

	if recorder.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusInternalServerError)
	}

	var got ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode the body: %v", err)
	}
	if got.Message != "the table metadata failed" {
		t.Fatalf("message: got %q", got.Message)
	}
}

func TestListRowsReadError(t *testing.T) {
	handler := NewHandler(fakeRegistry{}, fakeRows{err: errors.New("boom")},
		discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables/users/rows", nil))

	if recorder.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusInternalServerError)
	}

	var got ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode the body: %v", err)
	}
	if got.Message != "the row list failed" {
		t.Fatalf("message: got %q", got.Message)
	}
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
