package postgres

import "testing"

func TestParseTableList(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{name: "empty", raw: "", want: nil},
		{name: "single", raw: "users", want: []string{"users"}},
		{name: "comma separated", raw: "users,orders", want: []string{"users", "orders"}},
		{name: "spaces", raw: "users, orders , invoices", want: []string{"users", "orders", "invoices"}},
		{name: "trailing comma", raw: "users,orders,", want: []string{"users", "orders"}},
		{name: "leading comma", raw: ",users", want: []string{"users"}},
		{name: "only commas", raw: ",,", want: nil},
		{name: "only spaces", raw: "  ", want: nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ParseTableList(test.raw)
			if len(got) != len(test.want) {
				t.Fatalf("ParseTableList(%q) = %v, want %v", test.raw, got, test.want)
			}
			for i, v := range got {
				if v != test.want[i] {
					t.Fatalf("ParseTableList(%q)[%d] = %q, want %q",
						test.raw, i, v, test.want[i])
				}
			}
		})
	}
}
