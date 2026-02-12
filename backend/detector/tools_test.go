package detector

import "testing"

func TestIsNVEncVersionSupported(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"NVEncC 8.00", true},
		{"NVEncC 10.2", true},
		{"NVEncC 7.80", false},
		{"unknown", false},
	}
	for _, c := range cases {
		if got := isNVEncVersionSupported(c.in); got != c.want {
			t.Fatalf("version support mismatch for %q: got=%v want=%v", c.in, got, c.want)
		}
	}
}
