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
