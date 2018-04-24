package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pborman/getopt/v2"
)

const DEFAULT_INTERVAL int = 100
const COUNT_FORMAT string = "%8d"
const CARRIAGE_RETURN byte = 13
const SPACE byte = 32

type ScanFunc func(*bufio.Scanner, *uint64)

type Config struct {
	countLines    bool
	countWords    bool
	countChars    bool
	countBytes    bool
	maxLineLength bool
	interval      time.Duration
	help          bool
	version       bool
	files         []string
}

type Processor struct {
	split bufio.SplitFunc
	scan  ScanFunc
}

var version = "master"

func buildConfig(config *Config) {
	intervalMs := DEFAULT_INTERVAL
	getopt.FlagLong(&config.countLines, "lines", 'l', "print the newline counts")
	getopt.FlagLong(&config.countWords, "words", 'w', "print the word counts")
	getopt.FlagLong(&config.countChars, "chars", 'm', "print the character counts")
	getopt.FlagLong(&config.countBytes, "bytes", 'c', "print the byte counts")
	getopt.FlagLong(&config.maxLineLength, "max-line-length", 'L', "print the maximum display width")
	getopt.FlagLong(&intervalMs, "interval", 'i',
		fmt.Sprintf("set update interval in ms (default %d ms)", DEFAULT_INTERVAL))
	getopt.FlagLong(&config.help, "help", 'h', "display this help and exit")
	getopt.FlagLong(&config.version, "version", 'V', "output version information and exit")
	getopt.Parse()
	config.interval = time.Duration(intervalMs * 1e6)
	config.files = getopt.Args()
	if !(config.countLines || config.countWords || config.countChars || config.countBytes) {
		config.countLines = true
		config.countWords = true
		config.countBytes = true
	}
}

func scanMaxLength(scanner *bufio.Scanner, count *uint64) {
	var max uint64
	var length uint64
	for scanner.Scan() {
		length = uint64(len(scanner.Text()))
		if length > max {
			max = length
			atomic.StoreUint64(count, max)
		}
	}
}

func scanCount(scanner *bufio.Scanner, count *uint64) {
	for scanner.Scan() {
		atomic.AddUint64(count, 1)
	}
}

func buildProcessors(config *Config, processors *[]Processor) {
	if config.countLines {
		*processors = append(*processors, Processor{bufio.ScanLines, scanCount})
	}
	if config.countWords {
		*processors = append(*processors, Processor{bufio.ScanWords, scanCount})
	}
	if config.countChars {
		*processors = append(*processors, Processor{bufio.ScanRunes, scanCount})
	}
	if config.countBytes {
		*processors = append(*processors, Processor{bufio.ScanBytes, scanCount})
	}
	if config.maxLineLength {
		*processors = append(*processors, Processor{bufio.ScanLines, scanMaxLength})
	}
}

func openFile(name string) *os.File {
	if name == "-" {
		return os.Stdin
	}
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func printCounts(counts *[]uint64, label string, cr bool) {
	var sb strings.Builder
	if cr {
		sb.WriteByte(CARRIAGE_RETURN)
	}
	sb.WriteString(fmt.Sprintf(COUNT_FORMAT, (*counts)[0]))
	for i := 1; i < len(*counts); i++ {
		sb.WriteByte(SPACE)
		sb.WriteString(fmt.Sprintf(COUNT_FORMAT, (*counts)[i]))
	}
	if label != "" {
		sb.WriteByte(SPACE)
		sb.WriteString(label)
	}
	os.Stdout.WriteString(sb.String())
}

func pollCounts(name string, counts *[]uint64, interval time.Duration, done chan bool) {
	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			printCounts(counts, name, true)
		case <-done:
			break
		}
	}
}

func consumeReader(reader *io.PipeReader, processor Processor, count *uint64, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(reader)
	scanner.Split(processor.split)
	processor.scan(scanner, count)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func pipeSource(file *os.File, pws []*io.PipeWriter) {
	numCounts := len(pws)
	writers := make([]io.Writer, numCounts)
	for i := 0; i < numCounts; i++ {
		defer pws[i].Close()
		writers[i] = io.Writer(pws[i])
	}
	writer := io.MultiWriter(writers...)
	if _, err := io.Copy(writer, file); err != nil {
		log.Fatal(err)
	}
}

func processFile(file *os.File, name string, processors []Processor, totals *[]uint64, interval time.Duration) {
	numCounts := len(processors)

	// Create counters
	counts := make([]uint64, numCounts)

	// Write zeroes straightaway in case file is empty
	printCounts(&counts, name, false)

	// For each counter, set up a pipe
	prs := make([]*io.PipeReader, numCounts)
	pws := make([]*io.PipeWriter, numCounts)
	for i := 0; i < numCounts; i++ {
		prs[i], pws[i] = io.Pipe()
	}

	// Update stdout at fixed intervals
	done := make(chan bool)
	go pollCounts(name, &counts, interval, done)

	// Set up WaitGroup for our goroutines
	var wg sync.WaitGroup
	wg.Add(numCounts)

	// Start reading from pipes
	for index, processor := range processors {
		go consumeReader(prs[index], processor, &counts[index], &wg)
	}

	// Write to pipes
	pipeSource(file, pws)

	// Wait for goroutines to complete
	wg.Wait()

	// Stop polling
	done <- true

	// Write final counts
	printCounts(&counts, name, true)
	fmt.Println()
}

func processFiles(config *Config, processors []Processor) {
	// If no files given, process stdin
	if len(config.files) == 0 {
		processFile(os.Stdin, "", processors, nil, config.interval)
		return
	}

	numCounts := len(processors)

	// If more than one file given, also calculate totals
	var totals []uint64
	if len(config.files) > 1 {
		totals = make([]uint64, numCounts)
	} else {
		totals = nil
	}

	// Process files sequentially
	for _, name := range config.files {
		file := openFile(name)
		var totalsPtr *[]uint64
		if totals != nil {
			totalsPtr = &totals
		}
		processFile(file, name, processors, totalsPtr, config.interval)
	}

	// If we were keeping totals, print them now
	if totals != nil {
		printCounts(&totals, "total\n", false)
	}
}

func main() {
	// Read command-line args
	var config Config
	buildConfig(&config)

	// If --version was passed, print version and exit
	if config.version {
		fmt.Println(version)
		return
	}

	// If --help was passed, print help and exit
	if config.help {
		getopt.PrintUsage(os.Stdout)
		return
	}

	// Determine which processors to use
	var processors []Processor
	buildProcessors(&config, &processors)

	// All set, let's go
	processFiles(&config, processors)
}
