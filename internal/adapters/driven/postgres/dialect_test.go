package postgres

import "testing"

func TestDialectQuote(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple name", input: "my_table", want: "\"my_table\""},
		{name: "name with spaces", input: "my table", want: "\"my table\""},
		{name: "name with dot", input: "schema.table", want: "\"schema\".\"table\""},
		{name: "name with quote", input: "ta\"ble", want: "\"ta\"\"ble\""},
		{name: "empty name", input: "", want: "\"\""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Dialect{}.Quote(test.input)
			if got != test.want {
				t.Fatalf("Quote(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestDialectMapValue(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		value     string
		want      string
		wantError bool
	}{
		{name: "simple text", field: "name", value: "hello", want: "hello"},
		{name: "empty value", field: "name", value: "", want: ""},
		{name: "value with spaces", field: "desc", value: "hello world", want: "hello world"},
		{name: "empty field", field: "", value: "v", wantError: true},
		{name: "value with NUL", field: "name", value: "ab\x00cd", wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Dialect{}.MapValue(test.field, test.value)
			if test.wantError {
				if err == nil {
					t.Fatalf("MapValue(%q, %q): got %q, want error",
						test.field, test.value, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("MapValue(%q, %q): unexpected error: %v",
					test.field, test.value, err)
			}
			if got != test.want {
				t.Fatalf("MapValue(%q, %q) = %q, want %q",
					test.field, test.value, got, test.want)
			}
		})
	}
}
