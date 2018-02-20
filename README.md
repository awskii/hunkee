# hunkee
Convenient way to parse logs

Currently unstable. No autotests, no benchmarks. But example works fine (with debug output for now).

## Usage
`go get github.com/awskii/hunkee`

```go
package main

import (
	"fmt"
	"time"

	"github.com/awskii/hunkee"
)

type Beach struct {
	ID   uint16    `hunk:"id"`
	Name string    `hunk:"name"`
	LoAc uint8     `hunk:"lo_ac"`
	Temp float32   `hunk:"temp"`
	Time time.Time `hunk:"time"`
	next int
}

var (
	format = ":time :id :name :lo_ac :temp"
	entry  = "02 Jan 06 15:04 MST 17522 Brighton 20 25.6"
)

func main() {
	var err error
	bch := new(Beach)

	workersAmount := 10
	parser, err := hunkee.NewParser(format, bch, workersAmount)
	if err != nil {
		panic(err)
	}

	iters := 100000000
	for i := 0; i < iters; i++ {
		if i%10000 == 0 {
			fmt.Printf("\r Progress: %d/%d", i, iters)
		}
		if err := parser.ParseLine(entry, bch); err != nil {
			fmt.Println(err)
		}
	}
}

```

## Benchmarks
```
BenchmarkParse-4     	 2000000	       953 ns/op	      32 B/op	       1 allocs/op
BenchmarkParseRE-4   	  500000	      2482 ns/op	     448 B/op	       6 allocs/op
```
