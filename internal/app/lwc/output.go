package lwc

import (
	"bytes"
	"fmt"
	"time"

	"github.com/timdp/lwc/internal/pkg/lwcutil"
)

const COUNT_FORMAT string = "%8d"

func FormatCounts(counts *[]uint64, name string, cr bool, lf bool) *bytes.Buffer {
	buf := new(bytes.Buffer)
	if cr {
		buf.WriteByte(CARRIAGE_RETURN)
	}
	buf.WriteString(fmt.Sprintf(COUNT_FORMAT, (*counts)[0]))
	for i := 1; i < len(*counts); i++ {
		buf.WriteByte(SPACE)
		buf.WriteString(fmt.Sprintf(COUNT_FORMAT, (*counts)[i]))
	}
	if name != "" {
		buf.WriteByte(SPACE)
		buf.WriteString(name)
	}
	if lf {
		buf.WriteByte(LINE_FEED)
	}
	return buf
}

func PrintCounts(counts *[]uint64, name string, cr bool, lf bool) {
	lwcutil.GetStdout().Write(FormatCounts(counts, name, cr, lf).Bytes())
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
