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
	parser     *Parser
	bch        = new(Beach)
	timeLayout = time.RFC822
	format     = ":time :id :name :lo_ac :temp"
	entry      = "02 Jan 06 15:04 MST 17522 Brighton 20 25.6"

	formatWithoutTime = ":id :name :lo_ac :temp"
	entryWithoutTime  = "17522 Brighton 20 25.6"

	LogRecordRegex   = regexp.MustCompile(`^(\d{2}\s[A-Z][a-z]{2}\s[0-9]{2}\s[0-9]{2}:[0-9]{2})\s([A-Z]{3})\s+(\d)+\s+(\w)+\s+(\d{2})\s+([0-9]+\.\d+)$`)
	LogRecordWTRegex = regexp.MustCompile(`^(\d{5})+\s+(\w)+\s+(\d{2})\s+([0-9]+\.\d+)$`)
)

func init() {
}

func BenchmarkParse(b *testing.B) {
	parser, err := NewParser(format, bch)
	if err != nil {
		panic(err)
	}
	parser.SetTimeLayout("time", timeLayout)

	for i := 0; i < b.N; i++ {
		if err := parser.ParseLine(entry, bch); err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkParseWithoutTime(b *testing.B) {
	parser, err := NewParser(formatWithoutTime, bch)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		if err := parser.ParseLine(entryWithoutTime, bch); err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkParseRE(b *testing.B) {
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

func BenchmarkParseREWithoutTime(b *testing.B) {
	bch := new(Beach)
	for i := 0; i < b.N; i++ {
		tokens := LogRecordWTRegex.FindStringSubmatch(entryWithoutTime)
		u16, _ := strconv.ParseUint(tokens[0], 10, 16)
		bch.ID = uint16(u16)
		bch.Name = tokens[1]
		u8, _ := strconv.ParseUint(tokens[2], 10, 8)
		bch.LoAc = uint8(u8)
		f32, _ := strconv.ParseFloat(tokens[3], 32)
		bch.Temp = float32(f32)
	}
}
