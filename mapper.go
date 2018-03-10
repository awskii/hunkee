package hunkee

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"sync"
	"unicode"
)

// mapper knows which field associated with tag,
// current position and token sequence
type mapper struct {
	// no mutexes because we write to fields and tokenSeq
	// only once when building up structure
	mu           sync.RWMutex
	fields       map[string]*field
	tokensSeq    []string
	comPrefix    string // skip line if line has such prefix
	prefixActive bool   // if false, prefix check will be disabled
}

// field represents structure field
type field struct {
	index       []int
	typ         reflect.Type // field Go type
	name        string       // field key
	hasRaw      bool         // signals that corresponded field has raw field too
	after       int          // offset after token to the next token
	position    int          // numeric position of token in format string
	timeOptions *TimeOption
}

type namedParameter struct {
	name   string // entry name without ':' (== tag)
	strPos int    // numeric position in format
	offset int    // count of symols after entry to next entry
}

func initMapper(format string, to interface{}) (*mapper, error) {
	// get info about entry
	tokens, err := extractNames(format)
	if err != nil {
		return nil, err
	}

	fields, err := extractFieldsOnTags(to)
	if err != nil {
		return nil, err
	}

	tokenSeq := make([]string, len(tokens))
	for i := 0; i < len(tokens); i++ {
		tokenSeq[i] = tokens[i].name

		if _, ok := fields[tokens[i].name]; !ok {
			return nil, fmt.Errorf("passed struct has no field with tag %q", tokens[i].name)
		}
		if i != tokens[i].strPos {
			panic("i != tokens[i].strPos")
		}
		fields[tokens[i].name].after = tokens[i].offset
		fields[tokens[i].name].position = tokens[i].strPos
		fields[tokens[i].name].name = tokens[i].name
	}

	return &mapper{
		fields: fields, tokensSeq: tokenSeq,
	}, nil
}

// raw returns raw field of passed in arg
func (m *mapper) raw(normal *field) *field {
	f, ok := m.fields[normal.name+"_raw"]
	if !ok {
		return nil
	}
	return f
}

// assume no any data writes here
func (m *mapper) getField(tag string) *field {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.fields[tag]
}

func (m *mapper) writeField(tag string, f *field) {
	m.mu.Lock()
	m.fields[tag] = f
	m.mu.Unlock()
}

func (m *mapper) aquireWorker() *worker {
	return &worker{parent: m}
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
		fmt.Println(v.Kind().String())
		return nil, ErrOnlyStructs
	}

	index := make(map[string]*field)

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		// Ignore anonymous and unexported fields
		if f.Anonymous || !v.CanSet() {
			continue
		}

		tag, normalizedTag, err := processTag(f.Tag)
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

		if f.Type == typeTime {
			index[tag].timeOptions = DefaultTimeOptions()
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

func extractNames(format string) ([]*namedParameter, error) {
	var (
		valid = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		names = make([]*namedParameter, 0)

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
				names = append(names, &namedParameter{
					name: name, strPos: pos,
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
