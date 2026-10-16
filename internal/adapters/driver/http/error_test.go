package http

import (
	"encoding/json"
	"testing"
)

func TestErrorResponseJSON(t *testing.T) {
	got, err := json.Marshal(ErrorResponse{
		Code:    "invalid_request",
		Message: "the body is invalid",
	})
	if err != nil {
		t.Fatalf("marshal ErrorResponse: %v", err)
	}

	want := `{"code":"invalid_request","message":"the body is invalid"}`
	if string(got) != want {
		t.Fatalf("ErrorResponse JSON: got %s, want %s", got, want)
	}
}
