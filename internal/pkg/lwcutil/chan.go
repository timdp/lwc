package lwcutil

import (
	"bufio"
	"io"
)

// EndOfFiles marks the end of the list of files
const EndOfFiles string = ""

// ValidateFileName checks whether the given string is a valid file name
func ValidateFileName(name string) {
	if len(name) == 0 {
		Fatal("invalid zero-length file name")
	}
}

// NewFilesChanFromSlice creates a files chan from a string slice
func NewFilesChanFromSlice(values []string) *chan string {
	c := make(chan string)
	go func() {
		for _, value := range values {
			ValidateFileName(value)
			c <- value
		}
		c <- EndOfFiles
	}()
	return &c
}

// NewFilesChanFromReader creates a files chan from a reader, one name per line
func NewFilesChanFromReader(reader io.Reader, separator byte) *chan string {
	c := make(chan string)
	scanner := bufio.NewScanner(reader)
	scanner.Split(ScanBytes(0, false))
	go func() {
		for scanner.Scan() {
			name := scanner.Text()
			ValidateFileName(name)
			c <- name
		}
		c <- EndOfFiles
	}()
	return &c
}
