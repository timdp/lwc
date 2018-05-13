package lwc

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/timdp/lwc/internal/pkg/lwcutil"
)

// Processor represents a scanner and splitter for input
type Processor struct {
	Split bufio.SplitFunc
	Scan  ScanFunc
}

func processReader(reader io.Reader, processor Processor, count *uint64, total *uint64) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(processor.Split)
	processor.Scan(scanner, count, total)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func openFile(namePtr *string) (string, *os.File) {
	if namePtr == nil {
		return "", os.Stdin
	}
	return *namePtr, lwcutil.OpenFile(*namePtr)
}

func processFile(namePtr *string, processors []Processor, totals *[]uint64, interval time.Duration) {
	// Open input file (can be stdin)
	name, file := openFile(namePtr)

	numCounts := len(processors)

	// Create counters
	counts := make([]uint64, numCounts)

	// For each counter, set up a pipe
	pipes := make([]lwcutil.Pipe, numCounts)
	for i := 0; i < numCounts; i++ {
		pipes[i] = lwcutil.NewPipe()
	}

	// Set up WaitGroup for our goroutines
	var wg sync.WaitGroup
	wg.Add(numCounts)

	// Start reading from pipes
	for index, processor := range processors {
		var totalPtr *uint64
		if totals != nil {
			totalPtr = &(*totals)[index]
		}
		go func(reader io.Reader, processor Processor, count *uint64, total *uint64) {
			defer wg.Done()
			processReader(reader, processor, count, total)
		}(pipes[index].R, processor, &counts[index], totalPtr)
	}

	// Update stdout at fixed intervals, but only if it's a terminal
	live := lwcutil.StdoutIsTTY()
	var done chan bool
	if live {
		done = make(chan bool)
		// Write zeroes straightaway in case file is empty
		PrintCounts(&counts, name, false, false)
		// Start updating stdout
		go PollCounts(name, &counts, interval, done)
	}

	// Write to pipes
	lwcutil.MultiPipe(file, lwcutil.GetPipeWriters(pipes))

	// Wait for goroutines to complete
	wg.Wait()

	// Stop polling
	if done != nil {
		done <- true
	}

	// Write final counts
	PrintCounts(&counts, name, live, true)
}

// ProcessFiles processes files as required by the configuration
func ProcessFiles(config *Config) {
	files := config.FilesChan()
	processors := config.Processors()

	name1 := <-*files

	// If no files given, process stdin
	if name1 == "" {
		processFile(nil, processors, nil, config.Interval)
		return
	}

	numCounts := len(processors)
	var totals *[]uint64

	name2 := <-*files

	// If more than one file given, also calculate totals
	if name2 != "" {
		totalsRaw := make([]uint64, numCounts)
		totals = &totalsRaw
	}

	processFile(&name1, processors, totals, config.Interval)

	if name2 != "" {
		processFile(&name2, processors, totals, config.Interval)

		// Process files sequentially
		for name := range *files {
			if name == lwcutil.EndOfFiles {
				break
			}
			processFile(&name, processors, totals, config.Interval)
		}
	}

	// If we were keeping totals, print them now
	if totals != nil {
		PrintCounts(totals, "total\n", false, false)
	}
}
