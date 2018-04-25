package lwc

import (
	"fmt"
	"log"
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
	g             *getopt.Set
}

func (c *Config) PrintUsage() {
	c.g.PrintUsage(os.Stdout)
}

func BuildConfig(args []string) Config {
	intervalMs := DEFAULT_INTERVAL
	g := getopt.New()
	var config Config
	config.g = g
	g.FlagLong(&config.Lines, "lines", 'l', "print the newline counts")
	g.FlagLong(&config.Words, "words", 'w', "print the word counts")
	g.FlagLong(&config.Chars, "chars", 'm', "print the character counts")
	g.FlagLong(&config.Bytes, "bytes", 'c', "print the byte counts")
	g.FlagLong(&config.MaxLineLength, "max-line-length", 'L', "print the maximum display width")
	g.FlagLong(&intervalMs, "interval", 'i',
		fmt.Sprintf("set update interval in ms (default %d ms)", DEFAULT_INTERVAL))
	g.FlagLong(&config.Help, "help", 'h', "display this help and exit")
	g.FlagLong(&config.Version, "version", 'V', "output version information and exit")
	g.Parse(args)
	if intervalMs < 0 {
		log.Fatal("Update interval cannot be negative")
	}
	config.Interval = time.Duration(intervalMs) * time.Millisecond
	config.Files = g.Args()
	if !(config.Lines || config.Words || config.Chars || config.Bytes) {
		config.Lines = true
		config.Words = true
		config.Bytes = true
	}
	return config
}
