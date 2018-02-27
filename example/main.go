package main

import (
	"sync"
	"time"

	"bufio"
	"log"
	"os"

	"github.com/awskii/hunkee"
)

type Beach struct {
	ID   uint16    `hunk:"id"`
	Name string    `hunk:"name"`
	LoAc uint8     `hunk:"lo_ac"`
	Temp float32   `hunk:"temp"`
	Time time.Time `hunk:"time"`
}

var format = ":time :id :name :lo_ac :temp "

func main() {
	var err error
	bch := new(Beach)

	parser, err := hunkee.NewParser(format, bch)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(os.Args[1])
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
