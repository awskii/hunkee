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
	timeLayout = time.RFC3339
	format     = ":time :id :name :lo_ac :temp"
	entry      = "2018-07-28T21:10:45Z 17522 Brighton 20 25.6"

	formatWithoutTime = ":id :name :lo_ac :temp"
	entryWithoutTime  = "17522 Brighton 20 25.6"

	LogRecordRegex   = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)\s+(\d+)\s+(\w+)\s+(\d+)\s+([0-9]+\.\d+)$`)
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
		bch.Time, _ = time.Parse(timeLayout, tokens[1])
		u16, _ := strconv.ParseUint(tokens[2], 10, 16)
		bch.ID = uint16(u16)
		bch.Name = tokens[3]
		u8, _ := strconv.ParseUint(tokens[4], 10, 8)
		bch.LoAc = uint8(u8)
		f32, _ := strconv.ParseFloat(tokens[5], 32)
		bch.Temp = float32(f32)
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
