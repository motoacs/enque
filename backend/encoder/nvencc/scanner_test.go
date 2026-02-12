package nvencc

import (
	"bufio"
	"bytes"
	"testing"
)

func TestScanCRLF(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			name:   "LF only",
			input:  "line1\nline2\nline3\n",
			expect: []string{"line1", "line2", "line3"},
		},
		{
			name:   "CR only (progress update)",
			input:  "progress1\rprogress2\rprogress3\r",
			expect: []string{"progress1", "progress2", "progress3"},
		},
		{
			name:   "CRLF pairs",
			input:  "line1\r\nline2\r\nline3\r\n",
			expect: []string{"line1", "line2", "line3"},
		},
		{
			name:   "mixed CR and LF",
			input:  "header\ninfo\rprogress1\rprogress2\r\nfinal\n",
			expect: []string{"header", "info", "progress1", "progress2", "final"},
		},
		{
			name:   "no trailing delimiter",
			input:  "last line",
			expect: []string{"last line"},
		},
		{
			name:   "empty lines",
			input:  "a\n\nb\n",
			expect: []string{"a", "", "b"},
		},
		{
			name:   "empty input",
			input:  "",
			expect: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(bytes.NewReader([]byte(tt.input)))
			scanner.Split(ScanCRLF)

			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				t.Fatal(err)
			}

			if len(lines) != len(tt.expect) {
				t.Fatalf("got %d lines %v, want %d lines %v", len(lines), lines, len(tt.expect), tt.expect)
			}
			for i := range lines {
				if lines[i] != tt.expect[i] {
					t.Errorf("line[%d]=%q, want %q", i, lines[i], tt.expect[i])
				}
			}
		})
	}
}
