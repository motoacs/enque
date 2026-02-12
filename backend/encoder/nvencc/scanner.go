package nvencc

import (
	"bufio"
	"bytes"
)

// NewCRLFScanner returns a bufio.Scanner that splits on both \r and \n.
// NVEncC outputs progress updates with \r (carriage return) overwriting the same line.
func NewCRLFScanner(r *bytes.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split(ScanCRLF)
	return s
}

// ScanCRLF is a split function for bufio.Scanner that splits on \r or \n.
func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the earliest \r or \n
	crIdx := bytes.IndexByte(data, '\r')
	lfIdx := bytes.IndexByte(data, '\n')

	switch {
	case crIdx >= 0 && lfIdx >= 0:
		if crIdx < lfIdx {
			// \r comes first
			if crIdx+1 == lfIdx {
				// \r\n pair: consume both, return content before \r
				return lfIdx + 1, data[:crIdx], nil
			}
			// standalone \r
			return crIdx + 1, data[:crIdx], nil
		}
		// \n comes first
		return lfIdx + 1, data[:lfIdx], nil

	case crIdx >= 0:
		// only \r found
		return crIdx + 1, data[:crIdx], nil

	case lfIdx >= 0:
		// only \n found
		return lfIdx + 1, data[:lfIdx], nil
	}

	// No delimiter found
	if atEOF {
		return len(data), data, nil
	}
	// Request more data
	return 0, nil, nil
}
