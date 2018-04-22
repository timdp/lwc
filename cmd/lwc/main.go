package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var countLines bool
var countWords bool
var countChars bool
var countBytes bool

var wg sync.WaitGroup

func getOpts() {
	for _, arg := range os.Args {
		switch {
		case arg == "-l" || arg == "--lines":
			countLines = true
		case arg == "-w" || arg == "--words":
			countWords = true
		case arg == "-m" || arg == "--chars":
			countChars = true
		case arg == "-c" || arg == "--bytes":
			countBytes = true
		}
	}
	if !(countLines || countWords || countChars || countBytes) {
		countLines = true
		countWords = true
		countBytes = true
	}
}

func buildSplits() []bufio.SplitFunc {
	var splits []bufio.SplitFunc
	if countLines {
		splits = append(splits, bufio.ScanLines)
	}
	if countWords {
		splits = append(splits, bufio.ScanWords)
	}
	if countChars {
		splits = append(splits, bufio.ScanRunes)
	}
	if countBytes {
		splits = append(splits, bufio.ScanBytes)
	}
	return splits
}

func consumeReader(reader *io.PipeReader, split bufio.SplitFunc, update chan int, i int) {
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

func printCounts(numCounts int, update chan int) {
	defer wg.Done()
	counts := make([]int, numCounts)
	out := make([]string, numCounts)
	for i := 0; i < numCounts; i++ {
		out[i] = "0"
	}
	done := 0
	for done < numCounts {
		i := <-update
		if i < 0 {
			done++
		} else {
			counts[i]++
			out[i] = strconv.Itoa(counts[i])
			fmt.Printf("\r%s", strings.Join(out, " "))
		}
	}
}

func pipeStdin(pws []*io.PipeWriter) {
	defer wg.Done()
	numCounts := len(pws)
	writers := make([]io.Writer, numCounts)
	for i := 0; i < numCounts; i++ {
		defer pws[i].Close()
		writers[i] = io.Writer(pws[i])
	}
	writer := io.MultiWriter(writers...)
	if _, err := io.Copy(writer, os.Stdin); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Read command-line args
	getOpts()

	// Determine which counters to use
	splits := buildSplits()
	numCounts := len(splits)

	// For each counter, set up a pipe for stdin
	prs := make([]*io.PipeReader, numCounts)
	pws := make([]*io.PipeWriter, numCounts)
	for i := 0; i < numCounts; i++ {
		prs[i], pws[i] = io.Pipe()
	}

	// Set up channel where counters will send updates
	update := make(chan int)

	// We have a goroutine for each counter, plus one for piping input, plus one
	// for producing output
	wg.Add(numCounts + 2)

	// Listen for updates to counters
	go printCounts(numCounts, update)

	// Start reading from pipes
	for i, split := range splits {
		go consumeReader(prs[i], split, update, i)
	}

	// Write to pipes
	go pipeStdin(pws)

	// Wait for goroutines to complete
	wg.Wait()

	// Write final newline
	fmt.Println()
}
