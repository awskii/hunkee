package hunkee

import (
	"errors"
	"time"
)

var (
	ErrSyntax           = errors.New("syntax error")
	ErrOnlyStructs      = errors.New("only struct types supported")
	ErrNotSpecified     = errors.New("tag not specified")
	ErrComaNotSupported = errors.New("coma-separated tag options is not supported")
	ErrUnexpectedColon  = errors.New("unexpected ':' while parsing format string")
)

type Parser struct {
	mapper *mapper
	debug  bool
}

func NewParser(format string, to interface{}) (*Parser, error) {
	mapper, err := initMapper(format, to)
	if err != nil {
		return nil, err
	}
	p := &Parser{
		mapper: mapper,
	}
	p.SetCommentPrefix("#")
	return p, nil
}

type TimeOption struct {
	Layout   string
	Location *time.Location
}

// SetDebug makes hunkee more verbose
func (p *Parser) SetDebug(val bool) {
	p.debug = val
	debug = val
}

// SetTimeLayout setups provided time layout for time.Time
// fields in log entry. By default it's corresponded to
// RFC3339 - "2006-01-02T15:04:05Z07:00"
func (p *Parser) SetTimeLayout(tag, timeLayout string) {
	p.mapper.fields[tag].timeOptions.Layout = timeLayout
}

func (p *Parser) SetMultiplyTimeLayout(tagToLayouts map[string]string) {
	for tag, layout := range tagToLayouts {
		p.SetTimeLayout(tag, layout)
	}
}

// SetTimeLocation used to parse time in provided location.
func (p *Parser) SetTimeLocation(tag string, loc *time.Location) {
	if loc == nil {
		panic("passed nil location")
	}
	p.mapper.fields[tag].timeOptions.Location = loc
}

// SetTimeOption sets provided timeOption to provided tag.
// Make sure you do it once at start, no andy dynamic behavior
func (p *Parser) SetTimeOption(tag string, to *TimeOption) {
	if to == nil {
		return
	}
	p.mapper.fields[tag].timeOptions = to
}

// TimeOption returns corresponded TimeOptions for tag
func (p *Parser) TimeOption(tag string) *TimeOption {
	if to := p.mapper.fields[tag].timeOptions; to != nil {
		return to
	}
	return nil
}

// SetCommentPrefix recieves prefix which from will be start commented lines.
// As soon parser will get string with such prefix, that line will be ignored.
// Default commentary prefix is '#'.
func (p *Parser) SetCommentPrefix(pref string) {
	p.mapper.comPrefix = pref
	p.mapper.prefixActive = true
}

// ParseLine gets line of input and structure to parse in
func (p *Parser) ParseLine(line string, to interface{}) error {
	return p.parseLine(line, to)
}

func DefaultTimeOptions() *TimeOption {
	return &TimeOption{
		Layout: time.RFC3339, // default time layout "2006-01-02T15:04:05Z07:00"
	}
}
