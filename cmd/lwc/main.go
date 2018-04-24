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
)

type Config struct {
	countLines bool
	countWords bool
	countChars bool
	countBytes bool
	version    bool
	files      []string
}

type Update struct {
	channel int
	count   uint64
	done    bool
}

const COUNT_FORMAT string = "%8d"
const UPDATE_INTERVAL time.Duration = 50 * time.Millisecond
const CARRIAGE_RETURN byte = 13
const SPACE byte = 32

var version = "master"

func buildConfig(config *Config) {
	for _, arg := range os.Args[1:] {
		switch {
		case arg == "-l" || arg == "--lines":
			config.countLines = true
		case arg == "-w" || arg == "--words":
			config.countWords = true
		case arg == "-m" || arg == "--chars":
			config.countChars = true
		case arg == "-c" || arg == "--bytes":
			config.countBytes = true
		case arg == "--version":
			config.version = true
		case arg != "-" && arg[0] == '-':
			log.Fatalf("Invalid option: %s", arg)
		default:
			config.files = append(config.files, arg)
		}
	}
	if !(config.countLines || config.countWords || config.countChars || config.countBytes) {
		config.countLines = true
		config.countWords = true
		config.countBytes = true
	}
}

func buildSplits(config *Config, splits *[]bufio.SplitFunc) {
	if config.countLines {
		*splits = append(*splits, bufio.ScanLines)
	}
	if config.countWords {
		*splits = append(*splits, bufio.ScanWords)
	}
	if config.countChars {
		*splits = append(*splits, bufio.ScanRunes)
	}
	if config.countBytes {
		*splits = append(*splits, bufio.ScanBytes)
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

func printCounts(counts []uint64, label string, cr bool) {
	var sb strings.Builder
	if cr {
		sb.WriteByte(CARRIAGE_RETURN)
	}
	sb.WriteString(fmt.Sprintf(COUNT_FORMAT, counts[0]))
	for i := 1; i < len(counts); i++ {
		sb.WriteByte(SPACE)
		sb.WriteString(fmt.Sprintf(COUNT_FORMAT, counts[i]))
	}
	if label != "" {
		sb.WriteByte(SPACE)
		sb.WriteString(label)
	}
	os.Stdout.WriteString(sb.String())
}

func consumeReader(reader *io.PipeReader, split bufio.SplitFunc, in chan bool, out chan Update, index int, wg *sync.WaitGroup) {
	defer wg.Done()
	var count uint64
	go func() {
		for <-in {
			out <- Update{index, atomic.LoadUint64(&count), false}
		}
	}()
	scanner := bufio.NewScanner(reader)
	scanner.Split(split)
	for scanner.Scan() {
		// fmt.Printf("%v: %v\n", i, scanner.Text())
		atomic.AddUint64(&count, 1)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	} else {
		out <- Update{index, atomic.LoadUint64(&count), true}
	}
}

func collectCounts(name string, numCounts int, totals *[]uint64, updates chan Update, wg *sync.WaitGroup) {
	defer wg.Done()

	counts := make([]uint64, numCounts)
	// Print zeroes straightaway in case file is empty
	printCounts(counts, name, false)

	completed := 0
	for completed < numCounts {
		update := <-updates
		counts[update.channel] = update.count
		if totals != nil {
			(*totals)[update.channel] = update.count
		}
		printCounts(counts, name, true)
		if update.done {
			completed++
		}
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

func processFile(file *os.File, name string, splits []bufio.SplitFunc, totals *[]uint64) {
	numCounts := len(splits)

	// For each counter, set up a pipe
	prs := make([]*io.PipeReader, numCounts)
	pws := make([]*io.PipeWriter, numCounts)
	for i := 0; i < numCounts; i++ {
		prs[i], pws[i] = io.Pipe()
	}

	// Set up channels
	ins := make([]chan bool, numCounts)
	for i := range ins {
		ins[i] = make(chan bool)
	}
	out := make(chan Update)

	// Set up WaitGroup for our goroutines
	var wg sync.WaitGroup
	wg.Add(numCounts + 1)

	// Request an update at fixed intervals
	tick := time.NewTicker(UPDATE_INTERVAL)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-tick.C:
				for _, in := range ins {
					in <- true
				}
			case <-done:
				for _, in := range ins {
					close(in)
				}
				break
			}
		}
	}()

	// Start listening for updates to counters
	go collectCounts(name, numCounts, totals, out, &wg)

	// Start reading from pipes
	for index, split := range splits {
		go consumeReader(prs[index], split, ins[index], out, index, &wg)
	}

	// Write to pipes
	pipeSource(file, pws)

	// Wait for goroutines to complete
	wg.Wait()

	done <- true
	tick.Stop()

	// Write final newline
	fmt.Println()
}

func processFiles(files []string, splits []bufio.SplitFunc) {
	// If no files given, process stdin
	if len(files) == 0 {
		processFile(os.Stdin, "", splits, nil)
		return
	}

	numCounts := len(splits)

	// If more than one file given, also calculate totals
	var totals []uint64
	if len(files) > 1 {
		totals = make([]uint64, numCounts)
	} else {
		totals = nil
	}

	// Process files sequentially
	for _, name := range files {
		file := openFile(name)
		var totalsPtr *[]uint64
		if totals != nil {
			totalsPtr = &totals
		}
		processFile(file, name, splits, totalsPtr)
	}

	// If we were keeping totals, print them now
	if totals != nil {
		printCounts(totals, "total\n", false)
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

	// Determine which counters to use
	var splits []bufio.SplitFunc
	buildSplits(&config, &splits)

	// All set, let's go
	processFiles(config.files, splits)
}
