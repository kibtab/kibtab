package http

// ErrorResponse is the JSON body that Kibtab sends for a failed request.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
