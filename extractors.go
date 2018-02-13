package hunkee

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"reflect"
	"time"
	"unicode"
)

const (
	unexportedTag = "-"
	libtag        = "hunk"
)

var (
	// Pre-defined supported structure types
	typeTime     = reflect.TypeOf(time.Time{})
	typeDuration = reflect.TypeOf(time.Second)
	typeURL      = reflect.TypeOf(url.URL{})
	typeIP       = reflect.TypeOf(net.IP{})
)

// splitTokens returns raw values of corresponding line
func splitTokens(mp *Mapper, line string, timeFormat string) []string {
	var (
		end    int
		offset int

		tokens = make([]string, len(mp.tokensSeq))
	)

	if debug {
		log.Printf("Entry: %q Len: %d\n", line, len(line))
	}

	for field, i := mp.First(), 0; field != nil; field = mp.Next() {
		// mapper guarantee that all names has fields
		if field.typ == typeTime {
			// if it's time make offset from offset to the end of value
			end = offset + len(_timeLayout)
		} else {
			end = findNextSpace(line, offset)
		}

		// findNextSpace returns -1 if no other space found
		// so if no space found - read line from current position
		// to the end of line, else read all between offset and end
		if end < offset {
			tokens[i] = line[offset:]
		} else {
			tokens[i] = line[offset:end]
		}

		if debug {
			log.Printf("Token: %q [%d:%d] After: %d\n", line[offset:end], offset, end, field.after)
		}

		// update current offset
		offset = end + field.after
		i++
	}

	return tokens
}

func findNextSpace(line string, start int) int {
	if start >= len(line) {
		return -1
	}

	for i := start + 1; i < len(line); i++ {
		if line[i] == ' ' {
			return i
		}
	}
	return -1
}

// TODO how to work with slices?
func extractFieldsOnTags(arg interface{}) (map[string]*field, error) {
	v := reflect.ValueOf(arg)

	// Maybe it's a pointer to struct - get it's value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// If v is still not a struct - done with error
	if v.Kind() != reflect.Struct {
		return nil, ErrOnlyStructs
	}

	index := make(map[string]*field)

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		// Ignore anonymous and unexported fields
		if f.Anonymous || !v.CanSet() {
			continue
		}

		tag, normalizedTag, err := procTag(f.Tag)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Check if field already indexed
		if _, ok := index[tag]; ok {
			index[tag].index = f.Index
			index[tag].typ = f.Type
		} else {
			index[tag] = &field{
				index: f.Index,
				typ:   f.Type,
			}
		}

		// Set .hasRaw flag to normal (non-raw) tag
		if normalizedTag != "" {
			if _, ok := index[normalizedTag]; ok {
				index[normalizedTag].hasRaw = true
			} else {
				index[normalizedTag] = &field{hasRaw: true}
			}
		}
	}
	return index, nil
}

func extractEntries(format string) ([]*namedEntry, error) {
	var (
		valid = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		names = make([]*namedEntry, 0)

		s      = []byte(format)
		pos    int
		offt   int
		inName bool
		name   string
	)

	for i := 0; i < len(s); i++ {
		if !inName {
			if s[i] == ':' {
				inName = true
				if len(names) != 0 {
					names[len(names)-1].offset = offt
					offt = 0
				}
				continue
			}
			// just another not interesting symbol
			offt++
			continue
		}

		if inName {
			if s[i] == ':' {
				// another one ':' and currently in word - error
				return nil, ErrUnexpectedColon
			}

			if unicode.IsSpace(rune(s[i])) {
				names = append(names, &namedEntry{
					name: name, str_pos: pos,
				})
				if debug {
					log.Printf("Field %q: %+v\n", name, names[len(names)-1])
				}

				inName = false
				name = ""
				pos++
				offt++
				continue
			}

			if !bytes.ContainsAny(s[i:i+1], valid) {
				return nil, fmt.Errorf("unsupported symol %q in format string at %d", s[i], i)
			}
			name += string(s[i])
		}
	}

	if debug {
		log.Printf("format string has been succesfully parsed")
	}
	return names, nil
}

// deref is Indirect for reflect.Types
func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
