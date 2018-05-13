package lwcutil

import (
	"bufio"
	"bytes"
)

// ScanBytes creates a SplitFunc that splits on the given byte value
func ScanBytes(b byte, requireEnd bool) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, b); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if atEOF && !requireEnd {
			return len(data), data, nil
		}
		return 0, nil, nil
	}
}

// ScanLines scans by line, accepting \r, \n, or \r\n as the separator
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		j := i
		if i > 0 && data[i-1] == '\r' {
			j = j - 1
		}
		return i + 1, data[0:j], nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
