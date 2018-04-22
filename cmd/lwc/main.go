package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type Config struct {
	countLines bool
	countWords bool
	countChars bool
	countBytes bool
	version    bool
	files      []string
}

const COUNT_FORMAT string = "%8d"
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

func printCounts(counts []int, label string, cr bool) {
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

func consumeReader(reader *io.PipeReader, split bufio.SplitFunc, update chan int, i int, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(reader)
	scanner.Split(split)
	for scanner.Scan() {
		// fmt.Printf("%v: %v\n", i, scanner.Text())
		update <- i
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	} else {
		update <- -1
	}
}

func collectCounts(name string, numCounts int, totals *[]int, update chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	counts := make([]int, numCounts)
	// Print zeroes straightaway in case file is empty
	printCounts(counts, name, false)

	completed := 0
	for completed < numCounts {
		index := <-update
		if index < 0 {
			// Received -1, so done processing
			completed++
		} else {
			// Received index, so increment corresponding counter(s)
			counts[index]++
			if totals != nil {
				(*totals)[index]++
			}
			// Go to start of line and print updated counts
			printCounts(counts, name, true)
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

func processFile(file *os.File, name string, splits []bufio.SplitFunc, totals *[]int) {
	numCounts := len(splits)

	// For each counter, set up a pipe for stdin
	prs := make([]*io.PipeReader, numCounts)
	pws := make([]*io.PipeWriter, numCounts)
	for i := 0; i < numCounts; i++ {
		prs[i], pws[i] = io.Pipe()
	}

	// Set up channel where counters will send updates
	update := make(chan int)

	// Set up WaitGroup for our goroutines
	var wg sync.WaitGroup
	wg.Add(numCounts + 1)

	// Start listening for updates to counters
	go collectCounts(name, numCounts, totals, update, &wg)

	// Start reading from pipes
	for i, split := range splits {
		go consumeReader(prs[i], split, update, i, &wg)
	}

	// Write to pipes
	pipeSource(file, pws)

	// Wait for goroutines to complete
	wg.Wait()

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
	var totals []int
	if len(files) > 1 {
		totals = make([]int, numCounts)
	} else {
		totals = nil
	}

	// Process files sequentially
	for _, name := range files {
		file := openFile(name)
		var totalsPtr *[]int
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
