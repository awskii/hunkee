package hunkee

import (
	"fmt"
	"testing"
	"time"
)

func TestNewParser(t *testing.T) {
	var s struct {
		Int int `hunk:"int"`
	}

	_, err := NewParser(":ab: :c", &s)
	if err == nil {
		t.Error("expected init error, got nil")
	}
}

func TestSetTimeLayout(t *testing.T) {
	var s struct {
		T  time.Time `hunk:"t"`
		Tr string    `hunk:"t_raw"`
	}

	p, err := NewParser(":t", &s)
	if err != nil {
		t.Error("init error " + err.Error())
	}

	p.SetTimeLayout("t", time.Kitchen)
	str := "5:43PM"
	if err := p.ParseLine(str, &s); err != nil {
		t.Error(err)
	}

	if s.T.Minute() != 43 {
		t.Errorf("unexpected value after parsing:\nhave: %d\nwant %d", s.T.Minute(), 43)
	}

	if s.Tr != "5:43PM" {
		t.Errorf("expected another raw value:\nhave: %q\nwant: %q", s.Tr, str)
	}
}

func TestSetMultiplyTimeLayouts(t *testing.T) {
	var s struct {
		A  time.Time `hunk:"a"`
		B  time.Time `hunk:"b"`
		C  time.Time `hunk:"c"`
		D  time.Time `hunk:"d"`
		E  time.Time `hunk:"e"`
		Cr string    `hunk:"c_raw"`
	}

	p, err := NewParser(":a :b :d :e :c", &s)
	if err != nil {
		t.Error("unexpected init error " + err.Error())
	}
	if p == nil {
		t.Fatal("returned nil parser")
	}

	layouts := map[string]string{
		"a": time.RFC3339,
		"b": time.RFC1123,
		"c": time.RFC1123,
		"d": time.Kitchen,
		"e": "2006-01-02 15:04:05",
	}

	v := map[string]string{
		"a": "2018-07-28T21:10:45+10:00",
		"b": "Tue, 10 Apr 2018 19:17:21 UTC",
		"c": "Tue, 10 Apr 2018 19:17:33 UTC",
		"d": "5:43PM",
		"e": "2006-01-02 03:04:05",
	}

	p.SetMultiplyTimeLayout(layouts)
	p.SetTokenSeparator('"')
	str := fmt.Sprintf("\"%s\" \"%s\" \"%s\" \"%s\" \"%s\"", v["a"], v["b"], v["d"], v["e"], v["c"])
	if err := p.ParseLine(str, &s); err != nil {
		t.Error(err)
	}

	if _, o := s.A.Zone(); o != 36000 || s.A.Year() != 2018 || s.A.Hour() != 21 {
		t.Errorf("wrong parsed time with options:\nhave: %q\nwant: %q", s.A.String(), v["a"])
	}

	if _, o := s.C.Zone(); o != 0 || s.C.Month() != 4 || s.C.Second() != 33 {
		t.Errorf("wrong parsed time with options:\nhave: %q\nwant: %q", s.C.String(), v["c"])
	}

	if s.D.Hour() != 17 || s.D.Minute() != 43 {
		t.Errorf("wrong parsed time with options:\nhave: %q\nwant: %q", s.D.String(), v["d"])
	}
}

func TestParseLine(t *testing.T) {
	var s struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}

	p, err := NewParser(":id :name", &s)
	if err != nil {
		t.Error("unexpected error: " + err.Error())
	}

	if err := p.parseLine("998 Gordon", &s); err != nil {
		t.Error(err)
	}

	if s.ID != 998 {
		t.Errorf("unexpected result of parsing commented string:\nhave: %d\nwant: %d", s.ID, 998)
	}
	if s.Name != "Gordon" {
		t.Errorf("unexpected result of parsing commented string:\nhave: %s\nwant: %s", s.Name, "Gordon")
	}
}

func TestParseCommentedLine(t *testing.T) {
	var s struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}

	p, err := NewParser(":id :name", &s)
	if err != nil {
		t.Error("unexpected error: " + err.Error())
	}

	if err := p.parseLine("#17 your_name_here\n", &s); err != nil {
		t.Error(err)
	}
	if s.ID != 0 {
		t.Errorf("unexpected result of parsing commented string:\nhave: %d\nwant: %d", s.ID, 0)
	}
}

func TestParseLineWithEscape(t *testing.T) {
	var s struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}

	p, err := NewParser(":id :name", &s)
	if err != nil {
		t.Error("unexpected error: " + err.Error())
	}

	p.SetTokenSeparator('"')
	if err := p.parseLine(`"998" "Gordon Freeman"`, &s); err != nil {
		t.Error(err)
	}

	if s.ID != 998 {
		t.Errorf("unexpected result of parsing commented string:\nhave: %d\nwant: %d", s.ID, 998)
	}
	if s.Name != "Gordon Freeman" {
		t.Errorf("unexpected result of parsing commented string:\nhave: %s\nwant: %s", s.Name, "Gordon Freeman")
	}
}
