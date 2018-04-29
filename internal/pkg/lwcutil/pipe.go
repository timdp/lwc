package lwcutil

import (
	"io"
	"log"
)

type Pipe struct {
	R *io.PipeReader
	W *io.PipeWriter
}

func NewPipe() Pipe {
	r, w := io.Pipe()
	return Pipe{r, w}
}

func GetPipeReaders(pipes []Pipe) []*io.PipeReader {
	readers := make([]*io.PipeReader, len(pipes))
	for i, p := range pipes {
		readers[i] = p.R
	}
	return readers
}

func GetPipeWriters(pipes []Pipe) []*io.PipeWriter {
	writers := make([]*io.PipeWriter, len(pipes))
	for i, p := range pipes {
		writers[i] = p.W
	}
	return writers
}

func MultiPipe(reader io.Reader, pws []*io.PipeWriter) {
	writers := make([]io.Writer, len(pws))
	for i, pw := range pws {
		defer pw.Close()
		writers[i] = io.Writer(pw)
	}
	writer := io.MultiWriter(writers...)
	if _, err := io.Copy(writer, reader); err != nil {
		log.Fatal(err)
	}
}
