package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// fakeRegistry is a TableRegistry for the transport tests.
type fakeRegistry struct {
	tables      []domain.TableMetadata
	err         error
	meta        domain.TableMetadata
	metadataErr error
}

func (f fakeRegistry) List() ([]domain.TableMetadata, error) {
	return f.tables, f.err
}

func (f fakeRegistry) Metadata(string) (domain.TableMetadata, error) {
	return f.meta, f.metadataErr
}

func TestListTables(t *testing.T) {
	registry := fakeRegistry{tables: []domain.TableMetadata{
		{Name: "users", Fields: []string{"id", "name"}, Version: 3},
		{Name: "orders", Fields: []string{"id"}, Version: 1},
	}}
	handler := NewHandler(registry, fakeRows{}, discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables", nil))

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}

	var got tableList
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode the body: %v", err)
	}
	if len(got.Tables) != 2 {
		t.Fatalf("tables: got %d, want 2", len(got.Tables))
	}
	if got.Tables[0].Name != "users" || got.Tables[0].Version != 3 {
		t.Fatalf("first table: got %+v", got.Tables[0])
	}
	if len(got.Tables[0].Fields) != 2 {
		t.Fatalf("first table fields: got %v", got.Tables[0].Fields)
	}
}

func TestListTablesEmpty(t *testing.T) {
	handler := NewHandler(fakeRegistry{}, fakeRows{}, discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables", nil))

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	if got, want := recorder.Body.String(), `{"tables":[]}`; got != want {
		t.Fatalf("body: got %q, want %q", got, want)
	}
}

func TestListTablesError(t *testing.T) {
	handler := NewHandler(fakeRegistry{err: errors.New("boom")}, fakeRows{},
		discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables", nil))

	if recorder.Code != nethttp.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusInternalServerError)
	}

	var got ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode the body: %v", err)
	}
	if got.Code != codeInternal {
		t.Fatalf("code: got %q, want %q", got.Code, codeInternal)
	}
	if got.Message != "the table list failed" {
		t.Fatalf("message: got %q", got.Message)
	}
}

func TestListTablesWriteError(t *testing.T) {
	h := handlers{
		tables: fakeRegistry{tables: []domain.TableMetadata{{Name: "users"}}},
		log:    discardLogger(),
	}
	writer := &failingWriter{header: nethttp.Header{}}

	h.listTables(writer, httptest.NewRequest(nethttp.MethodGet, "/v1/tables", nil))
}
