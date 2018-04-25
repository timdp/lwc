package lwcutil

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
)

var uiReady bool
var stdout io.Writer
var stdoutBuffer *bytes.Buffer
var fatal func(...interface{})
var LastError interface{}

func initUi() {
	if !uiReady {
		if flag.Lookup("test.v") != nil {
			stdoutBuffer = new(bytes.Buffer)
			stdout = bufio.NewWriter(stdoutBuffer)
			fatal = func(err ...interface{}) {
				LastError = err[0]
			}
		} else {
			stdout = os.Stdout
			fatal = log.Fatal
		}
		uiReady = true
	}
}

func GetStdout() io.Writer {
	initUi()
	return stdout
}

func Fatal(err interface{}) {
	initUi()
	fatal(err)
}

func FlushStdoutBuffer() []byte {
	stdout.(*bufio.Writer).Flush()
	b := stdoutBuffer.Bytes()
	stdoutBuffer.Reset()
	return b
}
