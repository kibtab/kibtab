package http

import nethttp "net/http"

// codeInternal is the code of a failure inside the instance.
const codeInternal = "internal_error"

// codeInvalidRequest is the code of a request that Kibtab rejects.
const codeInvalidRequest = "invalid_request"

// ErrorResponse is the JSON body that Kibtab sends for a failed request.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeError writes an ErrorResponse with the given status and code.
func (h handlers) writeError(w nethttp.ResponseWriter, status int, code, message string) {
	h.writeJSON(w, status, ErrorResponse{Code: code, Message: message})
}
