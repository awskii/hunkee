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
	return &Parser{
		mapper: mapper,
	}, nil
}

type TimeOption struct {
	Layout   string
	Location *time.Location
}

// SetDebug makes hunkee more verbose
func (p *Parser) SetDebug(val bool) {
	p.debug = val
}

// SetTimeLayout setups provided time layout for time.Time
// fields in log entry. By default it's corresponded to
// RFC822 (02 Jan 06 15:04 MST)
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

// ParseLine gets line of input and structure to parse in
func (p *Parser) ParseLine(line string, to interface{}) error {
	return p.parseLine(line, to)
}

func DefaultTimeOptions() *TimeOption {
	return &TimeOption{
		Layout: time.RFC822, // default time layout 02 Jan 06 15:04 MST
	}
}
