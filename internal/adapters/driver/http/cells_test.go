package http

import (
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// postCells sends one delta to the cells route.
func postCells(t *testing.T, registry fakeRegistry, syncService *fakeSync, body string) *httptest.ResponseRecorder {
	t.Helper()

	handler := NewHandler(registry, fakeRows{}, &fakeAudit{}, syncService, discardLogger())
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(
		nethttp.MethodPost, "/v1/tables/users/cells", strings.NewReader(body)))
	return recorder
}

func TestAddCells(t *testing.T) {
	syncService := &fakeSync{result: domain.SyncResult{
		RowResults: []domain.RowResult{
			{Row: "1", Version: 3},
			{Row: "2", Version: 3},
		},
	}}
	body := `{"delta":[{"row":"1","field":"name","value":"alice","version":3},` +
		`{"row":"2","field":"name","value":"bob","version":3}]}`

	recorder := postCells(t, usersTable(), syncService, body)

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	want := `{"rows":[{"key":"1","version":3},{"key":"2","version":3}]}`
	if got := recorder.Body.String(); got != want {
		t.Fatalf("body: got %q, want %q", got, want)
	}
	if syncService.payload.Table != "users" {
		t.Fatalf("table: got %q, want %q", syncService.payload.Table, "users")
	}
	if len(syncService.payload.Delta) != 2 {
		t.Fatalf("delta: got %+v", syncService.payload.Delta)
	}
	cell := syncService.payload.Delta[1]
	if cell.Row != "2" || cell.Field != "name" ||
		cell.Value != "bob" || cell.Version != 3 {
		t.Fatalf("delta: got %+v", cell)
	}
}

func TestAddCellsInvalidJSON(t *testing.T) {
	recorder := postCells(t, usersTable(), &fakeSync{}, "{")

	if recorder.Code != nethttp.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusBadRequest)
	}
	checkError(t, recorder, codeInvalidRequest, "the JSON body is invalid")
}

func TestAddCellsBodyTooLarge(t *testing.T) {
	body := `{"delta":[{"row":"` + strings.Repeat("a", maxBodyBytes) + `"}]}`

	recorder := postCells(t, usersTable(), &fakeSync{}, body)

	if recorder.Code != nethttp.StatusRequestEntityTooLarge {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusRequestEntityTooLarge)
	}
	checkError(t, recorder, codePayloadTooLarge, "the JSON body is too large")
}

func TestAddCellsConflict(t *testing.T) {
	conflict := domain.ConflictError{Table: "users", Row: "1", Want: 3, Got: 2}
	syncService := &fakeSync{result: domain.SyncResult{
		RowResults: []domain.RowResult{{Row: "1", Version: 3, Error: conflict}},
	}}
	body := `{"delta":[{"row":"1","field":"name","value":"alice","version":3}]}`

	recorder := postCells(t, usersTable(), syncService, body)

	if recorder.Code != nethttp.StatusConflict {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusConflict)
	}
	checkError(t, recorder, codeConflict, conflict.Error())
}
