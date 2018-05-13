package lwc

import (
	"bufio"
	"sync/atomic"
)

// ScanFunc represents a function that applies a scanner and updates counts
type ScanFunc func(*bufio.Scanner, *uint64, *uint64)

// ScanCount applies a scanner and adds one for every occurrence
func ScanCount(scanner *bufio.Scanner, count *uint64, total *uint64) {
	for scanner.Scan() {
		atomic.AddUint64(count, 1)
		if total != nil {
			atomic.AddUint64(total, 1)
		}
	}
}

// ScanMaxLength applies a scanner and remembers max string length
func ScanMaxLength(scanner *bufio.Scanner, count *uint64, total *uint64) {
	var localMax uint64
	var globalMax uint64
	if total != nil {
		globalMax = atomic.LoadUint64(total)
	}
	var length uint64
	for scanner.Scan() {
		length = uint64(len(scanner.Text()))
		if length > localMax {
			localMax = length
			atomic.StoreUint64(count, localMax)
		}
		if total != nil && length > globalMax {
			globalMax = length
			atomic.StoreUint64(total, globalMax)
		}
	}
}
