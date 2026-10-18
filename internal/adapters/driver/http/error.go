package http

import (
	"errors"
	nethttp "net/http"
)

// codeInternal is the code of a failure inside the instance.
const codeInternal = "internal_error"

// codeInvalidRequest is the code of a request that Kibtab rejects.
const codeInvalidRequest = "invalid_request"

// codeConflict is the code of a write that used a stale version.
const codeConflict = "conflict"

// codePayloadTooLarge is the code of a body above the request size limit.
const codePayloadTooLarge = "payload_too_large"

// ErrorResponse is the JSON body that Kibtab sends for a failed request.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeError writes an ErrorResponse with the given status and code.
func (h handlers) writeError(w nethttp.ResponseWriter, status int, code, message string) {
	h.writeJSON(w, status, ErrorResponse{Code: code, Message: message})
}

// writeDecodeError writes the ErrorResponse for a failed body read. A body
// above the request size limit gets a 413. Any other failure gets a 400.
func (h handlers) writeDecodeError(w nethttp.ResponseWriter, err error) {
	var tooLarge *nethttp.MaxBytesError
	if errors.As(err, &tooLarge) {
		h.writeError(w, nethttp.StatusRequestEntityTooLarge, codePayloadTooLarge,
			"the JSON body is too large")
		return
	}
	h.writeError(w, nethttp.StatusBadRequest, codeInvalidRequest,
		"the JSON body is invalid")
}
