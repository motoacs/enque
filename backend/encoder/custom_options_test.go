package encoder

import (
	"testing"
)

func TestTokenizeCustomOptions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "simple args",
			input: "--gop-len 300 --lookahead 16",
			want:  []string{"--gop-len", "300", "--lookahead", "16"},
		},
		{
			name:  "double quoted",
			input: `--metadata "title=My Video"`,
			want:  []string{"--metadata", "title=My Video"},
		},
		{
			name:  "single quoted",
			input: `--metadata 'title=My Video'`,
			want:  []string{"--metadata", "title=My Video"},
		},
		{
			name:  "escaped quotes in double",
			input: `--opt "value with \"quotes\""`,
			want:  []string{"--opt", `value with "quotes"`},
		},
		{
			name:  "escaped quotes in single",
			input: `--opt 'value with \'quotes\''`,
			want:  []string{"--opt", `value with 'quotes'`},
		},
		{
			name:  "multiple spaces",
			input: "  --a   --b  ",
			want:  []string{"--a", "--b"},
		},
		{
			name:  "tabs and newlines",
			input: "--a\t--b\n--c",
			want:  []string{"--a", "--b", "--c"},
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "only whitespace",
			input: "   \t\n ",
			want:  nil,
		},
		{
			name:    "unclosed double quote",
			input:   `--opt "unclosed`,
			wantErr: true,
		},
		{
			name:    "unclosed single quote",
			input:   `--opt 'unclosed`,
			wantErr: true,
		},
		{
			name:  "vpp filter",
			input: "--vpp-nlmeans sigma=0.005 --vpp-unsharp radius=3:weight=0.5",
			want:  []string{"--vpp-nlmeans", "sigma=0.005", "--vpp-unsharp", "radius=3:weight=0.5"},
		},
		{
			name:  "mixed quoting styles",
			input: `--a "hello world" --b 'foo bar' --c plain`,
			want:  []string{"--a", "hello world", "--b", "foo bar", "--c", "plain"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TokenizeCustomOptions(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenizeCustomOptions() error=%v, wantErr=%v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("token[%d]=%q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
