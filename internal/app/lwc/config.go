package lwc

import (
	"bufio"
	"fmt"
	"time"

	getopt "github.com/pborman/getopt/v2"
	"github.com/timdp/lwc/internal/pkg/lwcutil"
)

const DEFAULT_INTERVAL int = 100

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

func (c *Config) PrintUsage() {
	c.g.PrintUsage(lwcutil.GetStdout())
}

func NewConfig(args []string) *Config {
	intervalMs := DEFAULT_INTERVAL
	var c Config
	c.g = getopt.New()
	c.g.FlagLong(&c.Lines, "lines", 'l', "print the newline counts")
	c.g.FlagLong(&c.Words, "words", 'w', "print the word counts")
	c.g.FlagLong(&c.Chars, "chars", 'm', "print the character counts")
	c.g.FlagLong(&c.Bytes, "bytes", 'c', "print the byte counts")
	c.g.FlagLong(&c.MaxLineLength, "max-line-length", 'L', "print the maximum display width")
	c.g.FlagLong(&c.Files0From, "files0-from", 0, "read input from the files specified by NUL-terminated names in file F")
	c.g.FlagLong(&intervalMs, "interval", 'i',
		fmt.Sprintf("set update interval in ms (default %d ms)", DEFAULT_INTERVAL))
	c.g.FlagLong(&c.Help, "help", 'h', "display this help and exit")
	c.g.FlagLong(&c.Version, "version", 'V', "output version information and exit")
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

func (config *Config) Processors() []Processor {
	var temp [5]Processor
	i := 0
	if config.Lines {
		temp[i] = Processor{lwcutil.SplitOnByte(LINE_FEED, true), ScanCount}
		i++
	}
	if config.Words {
		temp[i] = Processor{bufio.ScanWords, ScanCount}
		i++
	}
	if config.Chars {
		temp[i] = Processor{bufio.ScanRunes, ScanCount}
		i++
	}
	if config.Bytes {
		temp[i] = Processor{bufio.ScanBytes, ScanCount}
		i++
	}
	if config.MaxLineLength {
		temp[i] = Processor{bufio.ScanLines, ScanMaxLength}
		i++
	}
	return temp[0:i]
}

func (config *Config) FilesChan() *chan string {
	if config.Files0From != "" {
		reader := lwcutil.OpenFile(config.Files0From)
		return lwcutil.NewFilesChanFromReader(reader, byte(0))
	} else {
		return lwcutil.NewFilesChanFromSlice(config.Files)
	}
}
