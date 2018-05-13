package lwcutil

import (
	"log"
	"os"
)

// OpenFile opens a file by name, or stdin if the name is a hyphen
func OpenFile(name string) *os.File {
	if name == "-" {
		return os.Stdin
	}
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
