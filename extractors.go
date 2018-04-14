package hunkee

import (
	"log"
	"reflect"
	"strings"
)

const (
	unexportedTag = "-"
	libtag        = "hunk"
)

var (
	debug bool
)

// parseLine processing one log line into structure
func (p *Parser) parseLine(line string, dest interface{}) (err error) {
	if line == "" || line == "\n" {
		return ErrEmptyLine
	}

	var (
		end     int
		offset  int
		lineLen = len(line)
		w       = p.mapper.aquireWorker()
	)

	if debug {
		log.Printf("Entry: %q Len: %d\n", line, len(line))
	}

	// Check if line has commentary prefix. If so, skip
	if p.mapper.prefixActive && strings.HasPrefix(line, p.mapper.comPrefix) {
		if debug {
			log.Printf("Entry: %q skipped due to matched prefix %q", line, p.mapper.comPrefix)
		}
		return
	}

	for field := w.first(); field != nil; field = w.next() {
		var token string
		// mapper guarantee that all names has fields
		if field.ftype == typeTime {
			// if it's time make offset from current offset to the end of value
			to := p.TimeOption(field.name)
			if to == nil {
				panic("not initialized TimeOptions for field " + field.name)
			}
			// TODO find other way to distinct end of time value. Issued when using time formats with
			// non-constant value length like RFC 1123
			end = offset + len(to.Layout)
		} else {
			end = findNextSpace(line, offset)
		}

		// findNextSpace returns -1 if no other space found
		// so if no space found - read line from current position
		// to the end of line, else read all between offset and end
		if end < offset || end >= lineLen-1 {
			token = line[offset:]
		} else {
			token = line[offset:end]
		}

		if debug {
			log.Printf("Token: %q [%d:%d] After: %d TimeOption: %#+v\n", token, offset, end, field.after, field.timeOptions)
		}

		destination := reflect.Indirect(reflect.ValueOf(dest))
		if err = p.mapper.processField(field, destination, token); err != nil {
			return err
		}

		// update current offset
		offset = end + field.after
	}
	return
}

func findNextSpace(line string, start int) int {
	if start >= len(line) {
		return -1
	}

	for i := start + 1; i < len(line); i++ {
		if line[i] == ' ' || line[i] == '\n' {
			return i
		}
	}
	return -1
}

// deref is Indirect for reflect.Types
func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
