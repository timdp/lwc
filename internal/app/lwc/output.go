package lwc

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const COUNT_FORMAT string = "%8d"
const CARRIAGE_RETURN byte = 13
const LINE_FEED byte = 10
const SPACE byte = 32

func FormatCounts(counts *[]uint64, label string, cr bool, lf bool) string {
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
	if lf {
		sb.WriteByte(LINE_FEED)
	}
	return sb.String()
}

func PrintCounts(counts *[]uint64, label string, cr bool, lf bool) {
	os.Stdout.WriteString(FormatCounts(counts, label, cr, lf))
}

func PollCounts(name string, counts *[]uint64, interval time.Duration, done chan bool) {
	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			PrintCounts(counts, name, true, false)
		case <-done:
			break
		}
	}
}
