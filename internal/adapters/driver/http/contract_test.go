package http

import (
	"encoding/json"
	"errors"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// routeContract holds the contract of one route. The success leg checks
// the status, the content type, and the body key. The failure leg checks
// the error envelope.
type routeContract struct {
	name     string
	method   string
	path     string
	body     string
	registry fakeRegistry
	sync     *fakeSync
	bodyKey  string
}

// routeContracts returns one contract for each route of the transport.
func routeContracts() []routeContract {
	write := domain.SyncResult{RowResults: []domain.RowResult{{Row: "1", Version: 3}}}
	return []routeContract{
		{
			name:   "GET /v1/tables",
			method: nethttp.MethodGet,
			path:   "/v1/tables",
			registry: fakeRegistry{tables: []domain.TableMetadata{
				{Name: "users", Fields: []string{"name"}, Version: 3},
			}},
			bodyKey: "tables",
		},
		{
			name:     "GET /v1/tables/{table}/rows",
			method:   nethttp.MethodGet,
			path:     "/v1/tables/users/rows",
			registry: usersTable(),
			bodyKey:  "rows",
		},
		{
			name:   "POST /v1/tables/{table}/rows",
			method: nethttp.MethodPost,
			path:   "/v1/tables/users/rows",
			body: `{"version":3,"key":"1",` +
				`"cells":[{"field":"name","value":"alice"}]}`,
			registry: usersTable(),
			sync:     &fakeSync{result: write},
			bodyKey:  "rows",
		},
		{
			name:   "POST /v1/tables/{table}/cells",
			method: nethttp.MethodPost,
			path:   "/v1/tables/users/cells",
			body: `{"delta":[{"row":"1","field":"name",` +
				`"value":"alice","version":3}]}`,
			registry: usersTable(),
			sync:     &fakeSync{result: write},
			bodyKey:  "rows",
		},
	}
}

// request builds the request of one contract case.
func (c routeContract) request(method string) *nethttp.Request {
	var body io.Reader
	if c.body != "" {
		body = strings.NewReader(c.body)
	}
	request := httptest.NewRequest(method, c.path, body)
	if c.body != "" {
		request.Header.Set("Content-Type", contentTypeJSON)
	}
	return request
}

func TestRouteContractSuccess(t *testing.T) {
	for _, c := range routeContracts() {
		t.Run(c.name, func(t *testing.T) {
			sync := c.sync
			if sync == nil {
				sync = &fakeSync{}
			}
			handler := NewHandler(c.registry, fakeRows{}, sync, discardLogger())
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, c.request(c.method))

			if recorder.Code != nethttp.StatusOK {
				t.Fatalf("status: got %d, want %d", recorder.Code,
					nethttp.StatusOK)
			}
			checkHeader(t, recorder, "Content-Type", contentTypeJSON)
			var body map[string]json.RawMessage
			if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode the body: %v", err)
			}
			if _, ok := body[c.bodyKey]; !ok {
				t.Fatalf("body %q misses the key %q", recorder.Body, c.bodyKey)
			}
		})
	}
}

func TestRouteContractFailure(t *testing.T) {
	for _, c := range routeContracts() {
		t.Run(c.name, func(t *testing.T) {
			registry := fakeRegistry{err: errors.New("boom"),
				metadataErr: errors.New("boom")}
			handler := NewHandler(registry, fakeRows{}, &fakeSync{},
				discardLogger())
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, c.request(c.method))

			if recorder.Code != nethttp.StatusInternalServerError {
				t.Fatalf("status: got %d, want %d", recorder.Code,
					nethttp.StatusInternalServerError)
			}
			checkHeader(t, recorder, "Content-Type", contentTypeJSON)
			var body ErrorResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode the body: %v", err)
			}
			if body.Code != codeInternal || body.Message == "" {
				t.Fatalf("error body: got %+v", body)
			}
		})
	}
}

func TestRouteContractMethod(t *testing.T) {
	for _, c := range routeContracts() {
		t.Run(c.name, func(t *testing.T) {
			handler := NewHandler(c.registry, fakeRows{}, &fakeSync{},
				discardLogger())
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, c.request(nethttp.MethodDelete))

			if recorder.Code != nethttp.StatusMethodNotAllowed {
				t.Fatalf("status: got %d, want %d", recorder.Code,
					nethttp.StatusMethodNotAllowed)
			}
		})
	}
}
