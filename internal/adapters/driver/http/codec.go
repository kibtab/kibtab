package http

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
)

// contentTypeJSON is the content type of a JSON body.
const contentTypeJSON = "application/json"

// encodeJSON writes value as a JSON response with the given status.
func encodeJSON(w nethttp.ResponseWriter, status int, value any) error {
	body, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("write JSON: %w", err)
	}
	return nil
}

// decodeJSON reads a JSON request body into value.
func decodeJSON(r *nethttp.Request, value any) error {
	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}
	return nil
}
