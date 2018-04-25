package lwc

import (
	"bufio"
	"strings"
	"testing"
)

func testScanner(t *testing.T, scan ScanFunc, input string, expectedCount uint64, withTotal bool, initialTotal uint64, expectedTotal uint64) {
	var actualCount uint64
	var actualTotalPtr *uint64
	if withTotal {
		actualTotalPtr = &initialTotal
	} else {
		actualTotalPtr = nil
	}
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	scan(scanner, &actualCount, actualTotalPtr)
	if expectedCount != actualCount {
		t.Errorf("Expecting count %d, got %d", expectedCount, actualCount)
	}
	if withTotal && expectedTotal != *actualTotalPtr {
		t.Errorf("Expecting total %d, got %d", expectedTotal, *actualTotalPtr)
	}
}

func TestScanCountWithoutTotal(t *testing.T) {
	testScanner(t,
		ScanCount,
		"one two three four five six",
		6,
		false,
		0,
		0)
}

func TestScanCountWithTotal(t *testing.T) {
	testScanner(t,
		ScanCount,
		"one two three four five six",
		6,
		true,
		0,
		6)
}

func TestScanMaxLengthWithoutTotal(t *testing.T) {
	testScanner(t,
		ScanMaxLength,
		"one two three four five six",
		5,
		false,
		0,
		0)
}

func TestScanMaxLengthWithLowerTotal(t *testing.T) {
	testScanner(t,
		ScanMaxLength,
		"one two three four five six",
		5,
		true,
		0,
		5)
}

func TestScanMaxLengthWithHigherTotal(t *testing.T) {
	testScanner(t,
		ScanMaxLength,
		"one two three four five six",
		5,
		true,
		6,
		6)
}
