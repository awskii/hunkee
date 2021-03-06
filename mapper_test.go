package hunkee

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestInitMapper(t *testing.T) {
	t.Parallel()

	type (
		tooEasy struct {
			ID    int       `hunk:"id"`
			Name  string    `hunk:"name"`
			Added time.Time `hunk:"added"`
		}

		easy struct {
			Id    uint64   `hunk:"id"`
			Token string   `hunk:"token"`
			Temp  float64  `hunk:"temp"`
			Nice  bool     `hunk:"nice"`
			IP    net.Addr `hunk:"ip"`
		}

		notSoEasy struct {
			Id          uint64        `hunk:"id"`
			IdRaw       string        `hunk:"id_raw"`
			Token       string        `hunk:"token"`
			TokenRaw    string        `hunk:"token_raw"`
			Temp        float64       `hunk:"temp"`
			Nice        bool          `hunk:"nice"`
			Ch          rune          `hunk:"ch"`
			Date        time.Time     `hunk:"date"`
			Dur         time.Duration `hunk:"dur"`
			ExplicitURL url.URL       `hunk:"explicit_url"`
			IP          net.Addr      `hunk:"ip"`
			ignoreIt    bool          `hunk:"ignore_it"`
			failWithIt  []byte        `hunk:"fail_with_it"`
		}

		embed struct {
			// tooEasy `hunk:"too_easy"`
			SoWhat bool `hunk:"so_what"`
			ST     struct {
				InTime   bool      `hunk:"in_time"`
				Arrival  time.Time `hunk:"arrival"`
				Token    string    `hunk:"token"`
				TicketID uint64    `hunk:"ticket_id"`
			} `hunk:"st"`
		}
		raw struct {
			Size    uint64 `hunk:"size"`
			SizeRaw string `hunk:"size_raw"`
		}
	)

	var (
		te  tooEasy
		e   easy
		nse notSoEasy
		em  embed
		r   raw

		tef  = ":id :name :added"
		ef   = ":id :temp :token :ip :nice"
		nsef = ":id :temp :token :ip :nice :ch :date :dur :explicit_url :fail_with_it"
		emf  = ":so_what :in_time :arrival :token :ticket_id"
		rawf = ":size"

		badWithPoint = ":id :name.name :added"
	)

	_, err := initMapper(tef, &te)
	if err != nil {
		t.Fatalf("Mapper initialization over %q should be finished without error, but have: %s", tef, err)
	}
	_, err = initMapper(badWithPoint, &te)
	if err == nil {
		t.Fatalf("Mapper initialization over %q should be finished with error of unexpected symbol, but no error occured", badWithPoint)
	}
	_, err = initMapper(ef, &e)
	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}
	_, err = initMapper(nsef, &nse)
	if err == nil {
		t.Fatal("expected error about absence of filed 'fail_with_it', got nil")
	}
	_, err = initMapper(emf, &em)
	if err == nil {
		t.Fatalf("Unexpected successfull finish of maper initialization")
	}
	_, err = initMapper(rawf, &r)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
}

func TestExtractNames(t *testing.T) {
	t.Parallel()

	// valid formats
	a := ":id :temp :token :ip :nice :ch :date :dur :explicit_url :ignore_it :fail_with_it"
	b := ":so_what :in_time :starrival :token :ticket_id"
	// c := `":id" ":temp" ":token" ":ip" ":nice" ":ch" ":date" ":dur" ":explicit_url" ":ignore_it" ":fail_with_it"`

	// invalid formats
	ia := ":id :temp :token :ip :nice :en:e"
	ib := ":id :temp.far :token :ip :nice :en:e"
	ic := ":so-what :ticket-id"

	// case A
	p, err := extractNames(a)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	la := 11
	if len(p) != la {
		t.Fatalf("%q - wrong length or extracted names: %d elements instead of %d", a, len(p), la)
	}

	// case B
	p, err = extractNames(b)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	lb := 5
	if len(p) != lb {
		t.Fatalf("wrong length or extracted names: %d elements instead of %d", len(p), lb)
	}

	// case IA (invalid A)
	expErr := "unexpected"
	p, err = extractNames(ia)
	for i := 0; i < len(p); i++ {
		fmt.Println(p[i].name)
	}
	if err == nil || !strings.Contains(err.Error(), expErr) {
		t.Fatalf("%q - expected error %q, have %q", ia, expErr+"..", err)
	}

	// case IB
	expErr = "unsupported symbol"
	_, err = extractNames(ib)
	if err == nil || !strings.Contains(err.Error(), expErr) {
		t.Fatalf("%q - expected error %q, have %q", ib, expErr+"..", err)
	}

	// case IC
	expErr = "unsupported symbol"
	p, err = extractNames(ic)
	if err == nil || !strings.Contains(err.Error(), expErr) {
		t.Fatalf("%q - expected error %q, have %q", ic, expErr+"..", err)
	}
}

func TestExtractFieldsOnTags(t *testing.T) {
	type (
		notSoEasy struct {
			Id            uint64        `hunk:"id"`
			IdRaw         string        `hunk:"id_raw"`
			Token         string        `hunk:"token"`
			TokenRaw      string        `hunk:"token_raw"`
			Temp          float64       `hunk:"temp"`
			TempRaw       string        `hunk:"temp_raw"`
			Nice          bool          `hunk:"nice"`
			Ch            rune          `hunk:"ch"`
			Date          time.Time     `hunk:"date"`
			Dur           time.Duration `hunk:"dur"`
			ExplicitURL   url.URL       `hunk:"explicit_url"`
			IP            net.Addr      `hunk:"ip"`
			ignoreIt      bool          `hunk:"ignore_it"`
			failWithIt    []byte        `hunk:"fail_with_it"`
			FailWithItToo []byte        `hunk:"Fail_with_it_too"`
		}

		withReader struct {
			Name string    `hunk:"name"`
			R    io.Reader `hunk:"r"`
		}
	)

	var (
		nse notSoEasy
		wr  withReader
	)

	f, err := extractFieldsOnTags(nse)
	if err != nil {
		t.Fatal(err.Error())
	}

	in := func(k string, valid []string) bool {
		for i := 0; i < len(valid); i++ {
			if valid[i] == k {
				return true
			}
		}
		return false
	}

	hasRaw := []string{"id", "token", "temp"}

	for k, v := range f {
		if in(k, hasRaw) && !v.hasRaw {
			t.Fatalf("field %q has raw in structure but that field was not detected", k)
		}
	}

	_, err = extractFieldsOnTags(wr)
	if err != nil {
		t.Fatal(err.Error())
	}

	var abc interface{}
	_, err = extractFieldsOnTags(abc)
	if err == nil {
		t.Fatalf("expected %s have nil error", ErrOnlyStructs)
	}
}

func TestMapperRaw(t *testing.T) {
	t.Parallel()

	var te struct {
		IDr   string    `hunk:"id_raw"`
		ID    int       `hunk:"id"`
		Name  string    `hunk:"name"`
		Added time.Time `hunk:"added"`
	}

	tef := ":id :name :added"
	m, err := initMapper(tef, &te)
	if err != nil {
		t.Fatalf("Mapper initialization over %q should be finished without error, but have: %s", tef, err)
	}

	f := m.raw(m.getField("id"))
	if f == nil {
		t.Error("expected non-nil _raw value, have nil")
	}

	f = m.raw(m.getField("name"))
	if f != nil {
		t.Error("unexpected non-nil _raw value, want nil")
	}
}
