package hunkee

import (
	"fmt"
	"reflect"
)

// Mapper knows which field associated with tag,
// current position and token sequence
type Mapper struct {
	fields    map[string]*field
	pos       int
	tokensSeq []string
}

// field represents structure field
type field struct {
	index    []int
	typ      reflect.Type // field Go type
	name     string       // field key
	hasRaw   bool         // signals that corresponded field has raw field too
	after    int          // offset after token to the next token
	position int          // numeric position of token in format string
}

type namedEntry struct {
	name    string // entry name without ':' (== tag)
	str_pos int    // numeric position in format
	offset  int    // count of symols after entry to next entry
}

func NewMapper(format string, to interface{}) (*Mapper, error) {
	tokens, err := extractEntries(format)
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
		if i != tokens[i].str_pos {
			panic("i != tokens[i].str_pos")
		}
		fields[tokens[i].name].after = tokens[i].offset
		fields[tokens[i].name].position = tokens[i].str_pos
		fields[tokens[i].name].name = tokens[i].name
	}

	return &Mapper{
		fields: fields, tokensSeq: tokenSeq,
	}, nil
}

func (m *Mapper) Seek(i int) *field {
	if i >= len(m.tokensSeq) {
		return nil
	}

	f := m.fields[m.tokensSeq[i]]
	m.pos = i + 1
	return f
}

func (m *Mapper) Next() *field {
	return m.Seek(m.pos)
}

func (m *Mapper) First() *field {
	return m.Seek(0)
}

// Raw returns raw field of passed in arg
func (m *Mapper) Raw(normal *field) *field {
	f, ok := m.fields[normal.name+"_raw"]
	if !ok {
		return nil
	}
	return f
}
