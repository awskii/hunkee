package hunkee

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"
)

type Beach struct {
	ID   uint16    `hunk:"id"`
	Name string    `hunk:"name"`
	LoAc uint8     `hunk:"lo_ac"`
	Temp float32   `hunk:"temp"`
	Time time.Time `hunk:"time"`
}

var (
	timeLayout = time.RFC822
	format     = ":time :id :name :lo_ac :temp"
	entry      = "02 Jan 06 15:04 MST 17522 Brighton 20 25.6"
)

func BenchmarkParse(b *testing.B) {
	SetTimeLayout(time.RFC822)
	SetDebug(true)
	bch := new(Beach)
	mp, _ := NewMapper(format, bch)
	for i := 0; i < b.N; i++ {
		if err := parse(mp, timeLayout, entry, bch); err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkParseRE(b *testing.B) {
	var LogRecordRegex = regexp.MustCompile(`^(\d{2}\s[A-Z][a-z]{2}\s[0-9]{2}\s[0-9]{2}:[0-9]{2})\s([A-Z]{3})\s+(\d)+\s+(\w)+\s+(\d{2})\s+([0-9]+\.\d+)$`)

	bch := new(Beach)
	for i := 0; i < b.N; i++ {
		tokens := LogRecordRegex.FindStringSubmatch(entry)
		u16, _ := strconv.ParseUint(tokens[1], 10, 16)
		bch.ID = uint16(u16)
		bch.Name = tokens[2]

		u8, _ := strconv.ParseUint(tokens[3], 10, 8)
		bch.LoAc = uint8(u8)
		f32, _ := strconv.ParseFloat(tokens[4], 32)
		bch.Temp = float32(f32)
		bch.Time, _ = time.Parse(timeLayout, tokens[0])
	}
}
