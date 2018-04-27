package lwcutil

import (
	"bufio"
	"io"
)

const END_OF_FILES string = ""

func ValidateFileName(name string) {
	if len(name) == 0 {
		Fatal("invalid zero-length file name")
	}
}

func NewFilesChanFromSlice(values []string) *chan string {
	c := make(chan string)
	go func() {
		for _, value := range values {
			ValidateFileName(value)
			c <- value
		}
		c <- END_OF_FILES
	}()
	return &c
}

func NewFilesChanFromReader(reader io.Reader, separator byte) *chan string {
	c := make(chan string)
	scanner := bufio.NewScanner(reader)
	scanner.Split(SplitOnByte(0, false))
	go func() {
		for scanner.Scan() {
			name := scanner.Text()
			ValidateFileName(name)
			c <- name
		}
		c <- END_OF_FILES
	}()
	return &c
}
