package lwc

import (
	"bytes"
	"fmt"
	"time"

	"github.com/timdp/lwc/internal/pkg/lwcutil"
)

const countFormat string = "%8d"

func formatCounts(counts *[]uint64, name string, cr bool, lf bool) *bytes.Buffer {
	buf := new(bytes.Buffer)
	if cr {
		buf.WriteByte(CarriageReturn)
	}
	buf.WriteString(fmt.Sprintf(countFormat, (*counts)[0]))
	for i := 1; i < len(*counts); i++ {
		buf.WriteByte(Space)
		buf.WriteString(fmt.Sprintf(countFormat, (*counts)[i]))
	}
	if name != "" {
		buf.WriteByte(Space)
		buf.WriteString(name)
	}
	if lf {
		buf.WriteByte(LineFeed)
	}
	return buf
}

// PrintCounts writes formatted counts to stdout
func PrintCounts(counts *[]uint64, name string, cr bool, lf bool) {
	lwcutil.GetStdout().Write(formatCounts(counts, name, cr, lf).Bytes())
}

// PollCounts periodically writes counts to stdout
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
