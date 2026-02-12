package encoder

import (
	"strings"

	"github.com/motoacs/enque/backend/model"
)

func TokenizeCustomOptions(input string) ([]string, error) {
	var tokens []string
	var buf strings.Builder
	inSingle := false
	inDouble := false
	escaped := false
	flush := func() {
		if buf.Len() > 0 {
			tokens = append(tokens, buf.String())
			buf.Reset()
		}
	}

	for _, r := range input {
		if escaped {
			buf.WriteRune(r)
			escaped = false
			continue
		}
		switch r {
		case '\\':
			escaped = true
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			} else {
				buf.WriteRune(r)
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			} else {
				buf.WriteRune(r)
			}
		case ' ', '\t', '\n', '\r':
			if inSingle || inDouble {
				buf.WriteRune(r)
			} else {
				flush()
			}
		default:
			buf.WriteRune(r)
		}
	}
	if escaped || inSingle || inDouble {
		return nil, &model.EnqueError{Code: model.ErrValidation, Message: "custom_options has unclosed quote or dangling escape"}
	}
	flush()
	return tokens, nil
}
