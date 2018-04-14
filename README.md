# hunkee
Convenient way to parse logs

All you need to parse log file - add "format stirng" and provide line to parse and structure to parse into.

You can specify raw field, which will be filled with raw value (string) of token.
Simple:
```go
type s struct {
  ID    int64  `hunk:"id"`
  IDRaw string `hunk:"id_raw"`
}
```
For that example format string might be as simple as `":id"`.

You can use raw values to parse not supported types.

Note that dots in tags are not supported. Embedded structs are not supported too.

## Supported types
* int, int8, int16, int32, int64
* uint, uint8, uint16, uint32, uint64
* bool
* string
* time.Time (with layout and timezone parsing)
* time.Duration
* net.IP
* url.URL

## Usage
Take a glance on that example (same at example/main.go):
```go
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
```

Note that all concurrency dispatch is lying on your shoulders.

## Benchmarks
```
goos: linux
goarch: amd64
pkg: github.com/awskii/hunkee
BenchmarkParse-4                	 1000000	      1061 ns/op	      32 B/op	       1 allocs/op
BenchmarkParseWithoutTime-4     	 5000000	       406 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseRE-4              	 1000000	      2192 ns/op	     448 B/op	       6 allocs/op
BenchmarkParseREWithoutTime-4   	 1000000	      1246 ns/op	     256 B/op	       4 allocs/op
```

## Don't be an enemy of yourself
If you passing an unsupported interface or structure, dont't start an issue about something goes wrong.
If you create structure with raw field of any other type than string, don't be confused.
