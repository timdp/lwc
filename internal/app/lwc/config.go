package lwc

import (
	"bufio"
	"fmt"
	"time"

	getopt "github.com/pborman/getopt/v2"
	"github.com/timdp/lwc/internal/pkg/lwcutil"
)

const defaultInterval int = 100

// Config represents a configuration
type Config struct {
	Lines         bool
	Words         bool
	Chars         bool
	Bytes         bool
	MaxLineLength bool
	Files0From    string
	Interval      time.Duration
	Help          bool
	Version       bool
	Files         []string
	g             *getopt.Set
}

// PrintUsage prints usage
func (c *Config) PrintUsage() {
	writer := lwcutil.GetStdout()
	c.g.PrintUsage(writer)
	fmt.Fprintln(writer, "\nFull documentation at: <https://github.com/timdp/lwc>")
}

// NewConfig creates a new Config
func NewConfig(args []string) *Config {
	intervalMs := defaultInterval
	var c Config
	c.g = getopt.New()
	c.g.SetParameters("[file ...]")
	c.g.FlagLong(&c.Lines, "lines", 'l', "print the newline counts")
	c.g.FlagLong(&c.Words, "words", 'w', "print the word counts")
	c.g.FlagLong(&c.Chars, "chars", 'm', "print the character counts")
	c.g.FlagLong(&c.Bytes, "bytes", 'c', "print the byte counts")
	c.g.FlagLong(&c.MaxLineLength, "max-line-length", 'L', "print the maximum display width")
	c.g.FlagLong(&c.Files0From, "files0-from", 0,
		"read input from the files specified by NUL-terminated names in file F",
		"F")
	c.g.FlagLong(&intervalMs, "interval", 'i',
		fmt.Sprintf("set the update interval to T ms (default %d ms)", defaultInterval),
		"T")
	c.g.FlagLong(&c.Help, "help", 0, "display this help and exit")
	c.g.FlagLong(&c.Version, "version", 0, "output version information and exit")
	c.g.Parse(args)
	if intervalMs < 0 {
		lwcutil.Fatal("Update interval cannot be negative")
	}
	c.Interval = time.Duration(intervalMs) * time.Millisecond
	c.Files = c.g.Args()
	if !(c.Lines || c.Words || c.Chars || c.Bytes || c.MaxLineLength) {
		c.Lines = true
		c.Words = true
		c.Bytes = true
	}
	return &c
}

// Processors creates the processors enabled by the configuration
func (c *Config) Processors() []Processor {
	var temp [5]Processor
	i := 0
	if c.Lines {
		temp[i] = Processor{lwcutil.ScanBytes(LineFeed, true), ScanCount}
		i++
	}
	if c.Words {
		temp[i] = Processor{bufio.ScanWords, ScanCount}
		i++
	}
	if c.Chars {
		temp[i] = Processor{bufio.ScanRunes, ScanCount}
		i++
	}
	if c.Bytes {
		temp[i] = Processor{bufio.ScanBytes, ScanCount}
		i++
	}
	if c.MaxLineLength {
		temp[i] = Processor{lwcutil.ScanLines, ScanMaxLength}
		i++
	}
	return temp[0:i]
}

// FilesChan creates the channel sending the filenames in the configuration
func (c *Config) FilesChan() *chan string {
	if c.Files0From != "" {
		reader := lwcutil.OpenFile(c.Files0From)
		return lwcutil.NewFilesChanFromReader(reader, byte(0))
	}
	return lwcutil.NewFilesChanFromSlice(c.Files)
}
