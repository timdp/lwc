package lwcutil

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// LastError is the last error produced during tests
var LastError interface{}

var uiReady bool
var stdout io.Writer
var stdoutBuffer *bytes.Buffer
var fatal func(...interface{})

func initUI() {
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

// GetStdout returns the writer that represents stdout in the environment's UI
func GetStdout() io.Writer {
	initUI()
	return stdout
}

// Fatal logs a fatal error to the environment's UI
func Fatal(err interface{}) {
	initUI()
	fatal(err)
}

// FlushStdoutBuffer flushes stdout for the environment's UI and returns its contents
func FlushStdoutBuffer() []byte {
	stdout.(*bufio.Writer).Flush()
	b := stdoutBuffer.Bytes()
	stdoutBuffer.Reset()
	return b
}

// StdoutIsTTY returns true if stdout is a terminal, false otherwise
func StdoutIsTTY() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
