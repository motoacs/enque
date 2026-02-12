package encoder

import "fmt"

// TokenizeCustomOptions splits a custom options string into tokens,
// supporting double and single quoted strings.
func TokenizeCustomOptions(s string) ([]string, error) {
	var tokens []string
	var current []byte
	inDouble := false
	inSingle := false
	escaped := false

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if escaped {
			current = append(current, ch)
			escaped = false
			continue
		}

		if ch == '\\' && (inDouble || inSingle) {
			if i+1 < len(s) {
				next := s[i+1]
				if (inDouble && (next == '"' || next == '\\')) ||
					(inSingle && (next == '\'' || next == '\\')) {
					escaped = true
					continue
				}
			}
			current = append(current, ch)
			continue
		}

		if ch == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}

		if ch == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}

		if (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r') && !inDouble && !inSingle {
			if len(current) > 0 {
				tokens = append(tokens, string(current))
				current = current[:0]
			}
			continue
		}

		current = append(current, ch)
	}

	if inDouble {
		return nil, fmt.Errorf("unclosed double quote")
	}
	if inSingle {
		return nil, fmt.Errorf("unclosed single quote")
	}

	if len(current) > 0 {
		tokens = append(tokens, string(current))
	}

	return tokens, nil
}
