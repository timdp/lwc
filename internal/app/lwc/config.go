package lwc

import (
	"fmt"
	"os"
	"time"

	getopt "github.com/pborman/getopt/v2"
)

const DEFAULT_INTERVAL int = 100

type Config struct {
	Lines         bool
	Words         bool
	Chars         bool
	Bytes         bool
	MaxLineLength bool
	Interval      time.Duration
	Help          bool
	Version       bool
	Files         []string
}

func BuildConfig() Config {
	var config Config
	intervalMs := DEFAULT_INTERVAL
	getopt.FlagLong(&config.Lines, "lines", 'l', "print the newline counts")
	getopt.FlagLong(&config.Words, "words", 'w', "print the word counts")
	getopt.FlagLong(&config.Chars, "chars", 'm', "print the character counts")
	getopt.FlagLong(&config.Bytes, "bytes", 'c', "print the byte counts")
	getopt.FlagLong(&config.MaxLineLength, "max-line-length", 'L', "print the maximum display width")
	getopt.FlagLong(&intervalMs, "interval", 'i',
		fmt.Sprintf("set update interval in ms (default %d ms)", DEFAULT_INTERVAL))
	getopt.FlagLong(&config.Help, "help", 'h', "display this help and exit")
	getopt.FlagLong(&config.Version, "version", 'V', "output version information and exit")
	getopt.Parse()
	config.Interval = time.Duration(intervalMs * 1e6)
	config.Files = getopt.Args()
	if !(config.Lines || config.Words || config.Chars || config.Bytes) {
		config.Lines = true
		config.Words = true
		config.Bytes = true
	}
	return config
}

func PrintUsage() {
	getopt.PrintUsage(os.Stdout)
}
