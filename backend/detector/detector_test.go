package detector

import (
	"testing"
)

func TestParseVersionString(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		want    string
		wantErr bool
	}{
		{
			name:   "standard header",
			output: "NVEncC (x64) 8.05 (r2994) by rigaya, Feb 10 2025 16:01:12 (VC 1942/Win)\n",
			want:   "8.05",
		},
		{
			name:   "version keyword",
			output: "version 8.10\n",
			want:   "8.10",
		},
		{
			name:   "version with three parts",
			output: "NVEncC64 version 8.10.1\n",
			want:   "8.10.1",
		},
		{
			name:    "no version",
			output:  "some random output\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersionString(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVersionString() err=%v, wantErr=%v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseVersionString()=%q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseMajorVersion(t *testing.T) {
	tests := []struct {
		version string
		want    int
	}{
		{"8.05", 8},
		{"8.10.1", 8},
		{"7.99", 7},
		{"10.0", 10},
	}

	for _, tt := range tests {
		got, err := parseMajorVersion(tt.version)
		if err != nil {
			t.Errorf("parseMajorVersion(%q) error: %v", tt.version, err)
			continue
		}
		if got != tt.want {
			t.Errorf("parseMajorVersion(%q)=%d, want %d", tt.version, got, tt.want)
		}
	}
}

func TestFindExecutable_NotFound(t *testing.T) {
	path := findExecutable("", []string{"nonexistent_tool_12345"})
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
}
