package encoder

import "testing"

func TestTokenizeCustomOptions(t *testing.T) {
	tokens, err := TokenizeCustomOptions(`--gop-len 300 --foo "bar baz" --x='y z'`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 5 {
		t.Fatalf("unexpected token count: %d %#v", len(tokens), tokens)
	}
	if tokens[3] != "bar baz" {
		t.Fatalf("quoted token mismatch: %s", tokens[3])
	}
	if tokens[4] != "--x=y z" {
		t.Fatalf("single quoted token mismatch: %s", tokens[4])
	}
}

func TestTokenizeCustomOptionsInvalidQuote(t *testing.T) {
	if _, err := TokenizeCustomOptions(`--foo "bar`); err == nil {
		t.Fatalf("expected validation error")
	}
}
