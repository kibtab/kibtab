package services

import "testing"

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name string
		a    int64
		b    int64
		want int64
	}{
		{name: "both positive, a larger", a: 5, b: 3, want: 5},
		{name: "both positive, b larger", a: 3, b: 5, want: 5},
		{name: "both positive, equal", a: 4, b: 4, want: 4},
		{name: "one zero", a: 0, b: 2, want: 0},
		{name: "one negative", a: -1, b: 2, want: 0},
		{name: "both negative", a: -1, b: -2, want: 0},
		{name: "both zero", a: 0, b: 0, want: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := CompareVersions(test.a, test.b); got != test.want {
				t.Fatalf("CompareVersions(%d, %d) = %d, want %d", test.a, test.b, got, test.want)
			}
		})
	}
}
