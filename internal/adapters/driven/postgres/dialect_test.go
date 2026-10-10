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

func TestDialectPage(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		first    int
		want     string
		wantArgs []any
	}{
		{name: "first placeholder", limit: 100, first: 1, want: " LIMIT $1", wantArgs: []any{100}},
		{name: "later placeholder", limit: 25, first: 3, want: " LIMIT $3", wantArgs: []any{25}},
		{name: "zero limit", limit: 0, first: 1, want: " LIMIT $1", wantArgs: []any{0}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, args := Dialect{}.Page(test.limit, test.first)
			if got != test.want {
				t.Fatalf("Page(%d, %d) = %q, want %q",
					test.limit, test.first, got, test.want)
			}
			if len(args) != len(test.wantArgs) {
				t.Fatalf("Page(%d, %d) returned %d args, want %d",
					test.limit, test.first, len(args), len(test.wantArgs))
			}
			if args[0] != test.wantArgs[0] {
				t.Fatalf("Page(%d, %d) arg = %v, want %v",
					test.limit, test.first, args[0], test.wantArgs[0])
			}
		})
	}
}
