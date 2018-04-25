package lwc

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/timdp/lwc/internal/pkg/lwcio"
)

type Processor struct {
	Split bufio.SplitFunc
	Scan  ScanFunc
}

func BuildProcessors(config *Config) []Processor {
	var temp [5]Processor
	i := 0
	if config.Lines {
		temp[i] = Processor{bufio.ScanLines, ScanCount}
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

func ProcessReader(reader io.Reader, processor Processor, count *uint64, total *uint64) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(processor.Split)
	processor.Scan(scanner, count, total)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func ProcessFile(file *os.File, name string, processors []Processor, totals *[]uint64, interval time.Duration) {
	numCounts := len(processors)

	// Create counters
	counts := make([]uint64, numCounts)

	// Write zeroes straightaway in case file is empty
	PrintCounts(&counts, name, false, false)

	// For each counter, set up a pipe
	pipes := make([]lwcio.Pipe, numCounts)
	for i := 0; i < numCounts; i++ {
		pipes[i] = lwcio.NewPipe()
	}

	// Update stdout at fixed intervals
	done := make(chan bool)
	go PollCounts(name, &counts, interval, done)

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
			ProcessReader(reader, processor, count, total)
		}(pipes[index].R, processor, &counts[index], totalPtr)
	}

	// Write to pipes
	lwcio.MultiPipe(file, lwcio.GetPipeWriters(pipes))

	// Wait for goroutines to complete
	wg.Wait()

	// Stop polling
	done <- true

	// Write final counts
	PrintCounts(&counts, name, true, true)
}

func ProcessFiles(config *Config, processors []Processor) {
	// If no files given, process stdin
	if len(config.Files) == 0 {
		ProcessFile(os.Stdin, "", processors, nil, config.Interval)
		return
	}

	numCounts := len(processors)

	// If more than one file given, also calculate totals
	var totals *[]uint64
	if len(config.Files) > 1 {
		totalsRaw := make([]uint64, numCounts)
		totals = &totalsRaw
	}

	// Process files sequentially
	for _, name := range config.Files {
		file := lwcio.OpenFile(name)
		ProcessFile(file, name, processors, totals, config.Interval)
	}

	// If we were keeping totals, print them now
	if totals != nil {
		PrintCounts(totals, "total\n", false, false)
	}
}
