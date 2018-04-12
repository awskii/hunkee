package main

import (
	"sync"
	"time"

	"bufio"
	"log"
	"os"

	"github.com/awskii/hunkee"
)

// Fill tags on your structure
type Beach struct {
	ID   uint16    `hunk:"id"`
	Name string    `hunk:"name"`
	LoAc uint8     `hunk:"lo_ac"`
	Temp float32   `hunk:"temp"`
	Time time.Time `hunk:"time"`
}

var (
	// names should match tags in struct
	format = ":time :id :name :lo_ac :temp "

	filePath = "./my.log"
)

func main() {
	bch := new(Beach)

	// Initialize parser with format string and structure
	parser, err := hunkee.NewParser(format, bch)
	if err != nil {
		panic(err)
	}

	// let hunkee to know that we await time from log at our format
	// Tag 'time' tells that field with such tag will be proceeded with
	// provided time format
	parser.SetTimeLayout("time", time.RFC822)

	f, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}
	r := bufio.NewReader(f)

	data := make([]*Beach, 0, 100)
	result := make(chan *Beach)
	wg := new(sync.WaitGroup)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			// io.EOF expected
			break
		}

		wg.Add(1)
		go func(line string, result chan *Beach) {
			// and parse each line into structure
			var b Beach
			if err := parser.ParseLine(line, &b); err != nil {
				log.Println(err)
			}
			result <- &b
		}(line, result)
	}

	go func() {
		for b := range result {
			data = append(data, b)
			wg.Done()
		}
	}()
	wg.Wait()
	log.Println("Done")

	f.Close()
}
