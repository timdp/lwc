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

type Opts struct {
	countLines bool
	countWords bool
	countChars bool
	countBytes bool
}

func getOpts() Opts {
	var opts Opts
	for _, arg := range os.Args {
		switch {
		case arg == "-l" || arg == "--lines":
			opts.countLines = true
		case arg == "-w" || arg == "--words":
			opts.countWords = true
		case arg == "-m" || arg == "--chars":
			opts.countChars = true
		case arg == "-c" || arg == "--bytes":
			opts.countBytes = true
		}
	}
	if !(opts.countLines || opts.countWords || opts.countChars || opts.countBytes) {
		opts.countLines = true
		opts.countWords = true
		opts.countBytes = true
	}
	return opts
}

func buildSplits(opts Opts) []bufio.SplitFunc {
	var splits []bufio.SplitFunc
	if opts.countLines {
		splits = append(splits, bufio.ScanLines)
	}
	if opts.countWords {
		splits = append(splits, bufio.ScanWords)
	}
	if opts.countChars {
		splits = append(splits, bufio.ScanRunes)
	}
	if opts.countBytes {
		splits = append(splits, bufio.ScanBytes)
	}
	return splits
}

func consumeReader(reader *io.PipeReader, split bufio.SplitFunc, update chan int, i int, wg sync.WaitGroup) {
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

func printCounts(numCounts int, update chan int, wg sync.WaitGroup) {
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

func pipeStdin(pws []*io.PipeWriter, wg sync.WaitGroup) {
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
	opts := getOpts()

	// Determine which counters to use
	splits := buildSplits(opts)
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
	wg.Add(numCounts + 2)

	// Start listening for updates to counters
	go printCounts(numCounts, update, wg)

	// Start reading from pipes
	for i, split := range splits {
		go consumeReader(prs[i], split, update, i, wg)
	}

	// Start writing to pipes
	go pipeStdin(pws, wg)

	// Wait for goroutines to complete
	wg.Wait()

	// Write final newline
	fmt.Println()
}
