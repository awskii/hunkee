package hunkee

import (
	"fmt"
	"strings"
	"testing"
)

func pointToRune(line string, at int) string {
	return fmt.Sprintf("%q\n%s", line, strings.Repeat(" ", at+1)+"^")
}

func TestFindNextSpace(t *testing.T) {
	var (
		lineA = "my precious line\n"
		lineB = ""
	)

	if i := findNextSep(lineA, 0, ' '); i != 2 {
		t.Errorf("wrong index provided:\n%s", pointToRune(lineA, i))
	}

	if i := findNextSep(lineA, 2, ' '); i != 11 {
		t.Errorf("wrong index provided:\n%s", pointToRune(lineA, i))
	}

	if i := findNextSep(lineA, 11, ' '); i != 16 {
		t.Errorf("wrong index provided:\n%s", pointToRune(lineA, i))
	}

	if i := findNextSep(lineB, 0, ' '); i != -1 {
		t.Errorf("wrong index provided: %d, expected %d", i, -1)
	}
}
