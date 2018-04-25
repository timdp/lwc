package lwc

import (
	"fmt"
	"os"
	"time"

	getopt "github.com/pborman/getopt/v2"
)

const DEFAULT_INTERVAL int = 100

type Config struct {
	CountLines    bool
	CountWords    bool
	CountChars    bool
	CountBytes    bool
	MaxLineLength bool
	Interval      time.Duration
	Help          bool
	Version       bool
	Files         []string
}

func BuildConfig() Config {
	var config Config
	intervalMs := DEFAULT_INTERVAL
	getopt.FlagLong(&config.CountLines, "lines", 'l', "print the newline counts")
	getopt.FlagLong(&config.CountWords, "words", 'w', "print the word counts")
	getopt.FlagLong(&config.CountChars, "chars", 'm', "print the character counts")
	getopt.FlagLong(&config.CountBytes, "bytes", 'c', "print the byte counts")
	getopt.FlagLong(&config.MaxLineLength, "max-line-length", 'L', "print the maximum display width")
	getopt.FlagLong(&intervalMs, "interval", 'i',
		fmt.Sprintf("set update interval in ms (default %d ms)", DEFAULT_INTERVAL))
	getopt.FlagLong(&config.Help, "help", 'h', "display this help and exit")
	getopt.FlagLong(&config.Version, "version", 'V', "output version information and exit")
	getopt.Parse()
	config.Interval = time.Duration(intervalMs * 1e6)
	config.Files = getopt.Args()
	if !(config.CountLines || config.CountWords || config.CountChars || config.CountBytes) {
		config.CountLines = true
		config.CountWords = true
		config.CountBytes = true
	}
	return config
}

func PrintUsage() {
	getopt.PrintUsage(os.Stdout)
}
