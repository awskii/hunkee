package hunkee

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

const (
	unexportedTag = "-"
	libtag        = "hunk"
)

var debug bool

// parseLine processing one log line into structure
func (p *Parser) parseLine(line string, dest interface{}) (err error) {
	if line == "" || line == "\n" {
		return ErrEmptyLine
	}

	var (
		offset  int
		lineLen = len(line)
		w       = p.mapper.aquireWorker()
	)

	defer w.free()

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

	destination := reflect.Indirect(reflect.ValueOf(dest))
	for field := w.first(); field != nil; field = w.next() {
		var token string

		start := offset
		// if token separator is 0, no need to search first occurrence
		if p.mapper.tokenSep != 0 {
			start = findNextSep(line, offset, p.mapper.tokenSep)
			if start < 0 {
				return errors.New("provided line has less tokens than expected")
			}
		}
		end := findNextSep(line, start, p.mapper.tokenSep)

		// findNextSpace returns -1 if no other space found
		// so if no space found - read line from current position
		// to the end of line, else read all between offset and end
		if end < offset || end >= lineLen-1 {
			token = line[offset:]
		} else {
			token = line[offset:end]
		}
		if field.ftype == typeIgnored {
			offset = end
			continue
		}

		token = strings.Trim(strings.TrimSpace(token), string(p.mapper.tokenSep))

		if debug {
			log.Printf("Token: %q [%d:%d] TimeOption: %#+v\n", token, offset, end, field.timeOptions)
		}

		if err = p.mapper.processField(field, destination, token); err != nil {
			return err
		}

		offset = end
	}
	return
}

// if provided sep is empty, space lookup will be used instead
func findNextSep(line string, start int, sep byte) int {
	if start >= len(line) {
		return -1
	}
	if sep == 0 {
		sep = ' '
	}

	for i := start + 1; i < len(line); i++ {
		if line[i] == sep || line[i] == '\n' {
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
