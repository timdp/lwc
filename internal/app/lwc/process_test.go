package lwc

import (
	"testing"
	"time"
)

func TestBuildProcessors(t *testing.T) {
	config := Config{
		true, true, true, true, true,
		time.Millisecond,
		false, false,
		[]string{},
		nil,
	}
	actualProcs := BuildProcessors(&config)
	actualCount := len(actualProcs)
	expectedCount := 5
	if expectedCount != actualCount {
		t.Fatalf("Expecting %d processors, got %d", expectedCount, actualCount)
	}
}
