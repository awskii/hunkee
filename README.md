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
	timeLayout = time.RFC822
	format     = ":time :id :name :lo_ac :temp"
	entry      = "02 Jan 06 15:04 MST 17522 Brighton 20 25.6"
)

func main() {
	b := new(Beach)
	hunkee.SetTimeLayout(time.RFC822)
	if err := hunkee.Parse(format, entry, b); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", b)
}

```

## Benchmarks
```
BenchmarkParse-4       	 1000000	      1151 ns/op	     112 B/op	       2 allocs/op
BenchmarkParse-16      	 1000000	      1129 ns/op	     112 B/op	       2 allocs/op
BenchmarkParseRE-4     	 1000000	      2367 ns/op	     448 B/op	       6 allocs/op
BenchmarkParseRE-16    	 1000000	      2364 ns/op	     448 B/op	       6 allocs/op
```
