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

// mapper knows which field associated with tag,
// current position and token sequence
type mapper struct {
	// no mutexes because we write to fields and tokenSeq
	// only once when building up structure
	fields       map[string]*field
	tokensSeq    []string
	tokenSep     byte   // byte which stead before and right after each token
	comPrefix    string // skip line if line has such prefix
	prefixActive bool   // if false, prefix check will be disabled
	workerPool   *pool
}

type fieldType int

const (
	typeIgnored fieldType = 1 << iota
	typeBool
	typeInt
	typeUint
	typeFloat
	typeString
	typeIP
	typeDuration
	typeURL
	typeTime
)

// field represents structure field
type field struct {
	index        []int
	ftype        fieldType
	reflectType  reflect.Type // field Go type
	reflectKind  reflect.Kind
	reflectValue reflect.Value
	name         string // field key
	hasRaw       bool   // signals that corresponded field has raw field too
	position     int    // numeric position of token in format string
	timeOptions  *TimeOption
}

type namedParameter struct {
	name   string // entry name without ':' (tag)
	strPos int    // numeric position in format string
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
		if tokenSeq[i] == "-" {
			fields["-"] = &field{ftype: typeIgnored}
		}

		if _, ok := fields[tokens[i].name]; !ok {
			return nil, fmt.Errorf("passed struct has no field with tag %q", tokens[i].name)
		}
		if i != tokens[i].strPos {
			panic("i != tokens[i].strPos")
		}
		fields[tokens[i].name].position = tokens[i].strPos
		fields[tokens[i].name].name = tokens[i].name
	}

	return &mapper{
		fields:     fields,
		tokensSeq:  tokenSeq,
		workerPool: initPool(10),
	}, nil
}

func (m *mapper) aquireWorker() *worker {
	return m.workerPool.get(m)
}

func (m *mapper) gainWorkers(upTo int) {
	final := upTo - m.workerPool.size
	for i := 0; i < final; i++ {
		m.workerPool.free <- new(worker)
	}
	m.workerPool.size = upTo
}

// raw returns raw field of passed in arg
func (m *mapper) raw(normal *field) *field {
	f, ok := m.fields[normal.name+"_raw"]
	if !ok {
		return nil
	}
	return f
}

func (m *mapper) getField(tag string) *field {
	return m.fields[tag]
}

func (m *mapper) writeField(tag string, f *field) {
	m.fields[tag] = f
}

func determineType(v interface{}) (ftype fieldType) {
	switch v.(type) {
	case time.Time:
		ftype = typeTime
	case time.Duration:
		ftype = typeDuration
	case net.IP:
		ftype = typeIP
	case url.URL, *url.URL:
		ftype = typeURL
	default:
		ftype = -1
	}
	return
}

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

		val := v.FieldByIndex(f.Index)
		// Ignore anonymous and unexported fields
		if f.Anonymous || !v.CanSet() || !val.CanInterface() {
			continue
		}

		tag, normalizedTag, err := processTag(f.Tag)
		if err != nil {
			return nil, err
		}

		// Check if field already indexed
		if _, ok := index[tag]; ok {
			index[tag].index = f.Index
		} else {
			ftype := determineType(val.Interface())
			index[tag] = &field{
				index:        f.Index,
				ftype:        ftype,
				reflectValue: val,
				reflectType:  f.Type,
				reflectKind:  f.Type.Kind(),
			}

			if ftype == typeTime {
				index[tag].timeOptions = DefaultTimeOptions()
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

func extractNames(format string) ([]*namedParameter, error) {
	var (
		valid = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		names = make([]*namedParameter, 0)

		s      = []byte(format)
		pos    int
		inName bool
		name   string
	)

	for i := 0; i < len(s); i++ {
		if !inName {
			switch s[i] {
			case ':':
				inName = true
			case '-': // ignore field
				names = append(names, &namedParameter{
					name: "-", strPos: pos,
				})
				name = ""
				pos++
			}
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
				continue
			}

			if !bytes.ContainsAny(s[i:i+1], valid) && s[i] != '\n' {
				return nil,
					fmt.Errorf("'%s': unsupported symbol %q in format string at pos %d", s, s[i], i)
			}

			// last symbol
			if i == len(s)-1 {
				if s[i] != '\n' {
					name += string(s[i])
				}
				names = append(names, &namedParameter{
					name: name, strPos: pos,
				})
				break
			}

			name += string(s[i])
		}
	}

	if debug {
		log.Println("format string has been successfully parsed")
	}
	return names, nil
}
