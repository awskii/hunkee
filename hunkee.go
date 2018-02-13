package hunkee

import (
	"errors"
	"reflect"
	"time"
)

var (
	debug bool

	_location   *time.Location
	_timeLayout = time.RFC822 // default time layout 02 Jan 06 15:04 MST

	ErrSyntax           = errors.New("syntax error")
	ErrOnlyStructs      = errors.New("only struct types supported")
	ErrNotSpecified     = errors.New("tag not specified")
	ErrComaNotSupported = errors.New("coma-separated tag options is not supported")
	ErrUnexpectedColon  = errors.New("unexpected ':' while parsing format string")
)

// SetDebug makes hunkee more verbose
func SetDebug(val bool) {
	debug = val
}

// SetTimeLayout setups provided time layout for time.Time
// fields in log entry. By default it's corresponded to
// RFC822 (02 Jan 06 15:04 MST)
func SetTimeLayout(timeLayout string) {
	_timeLayout = timeLayout
}

// SetTimeLocation used to parse time in provided location.
func SetTimeLocation(loc *time.Location) {
	if loc == nil {
		panic("passed nil location")
	}
	_location = loc
}

func Parse(format, line string, to interface{}) error {
	fields, err := NewMapper(format, to)
	if err != nil {
		return err
	}
	return parse(fields, _timeLayout, line, to)
}

func parse(fields *Mapper, timeFormat string, line string, to interface{}) error {
	// Get desitnation pointer
	destination := reflect.Indirect(reflect.ValueOf(to))
	// Split line on tokens
	tokens := splitTokens(fields, line, timeFormat)

	for field, i := fields.First(), 0; field != nil; field = fields.Next() {
		if err := processField(fields, field, destination, tokens[i]); err != nil {
			return err
		}
		i++
	}
	return nil
}
