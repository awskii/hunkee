package hunkee

import (
	"strings"
	"testing"
	"time"
)

// === Benchmark types ===

type benchStr1 struct {
	Name string `hunk:"name"`
}

type benchStr5 struct {
	A string `hunk:"a"`
	B string `hunk:"b"`
	C string `hunk:"c"`
	D string `hunk:"d"`
	E string `hunk:"e"`
}

type bench15 struct {
	F1  string  `hunk:"f1"`
	F2  int     `hunk:"f2"`
	F3  uint64  `hunk:"f3"`
	F4  string  `hunk:"f4"`
	F5  float64 `hunk:"f5"`
	F6  string  `hunk:"f6"`
	F7  int     `hunk:"f7"`
	F8  string  `hunk:"f8"`
	F9  uint32  `hunk:"f9"`
	F10 string  `hunk:"f10"`
	F11 int64   `hunk:"f11"`
	F12 string  `hunk:"f12"`
	F13 uint16  `hunk:"f13"`
	F14 float32 `hunk:"f14"`
	F15 string  `hunk:"f15"`
}

type bench25 struct {
	F1  string `hunk:"f1"`
	F2  string `hunk:"f2"`
	F3  string `hunk:"f3"`
	F4  string `hunk:"f4"`
	F5  string `hunk:"f5"`
	F6  string `hunk:"f6"`
	F7  string `hunk:"f7"`
	F8  string `hunk:"f8"`
	F9  string `hunk:"f9"`
	F10 string `hunk:"f10"`
	F11 string `hunk:"f11"`
	F12 string `hunk:"f12"`
	F13 string `hunk:"f13"`
	F14 string `hunk:"f14"`
	F15 string `hunk:"f15"`
	F16 string `hunk:"f16"`
	F17 string `hunk:"f17"`
	F18 string `hunk:"f18"`
	F19 string `hunk:"f19"`
	F20 string `hunk:"f20"`
	F21 string `hunk:"f21"`
	F22 string `hunk:"f22"`
	F23 string `hunk:"f23"`
	F24 string `hunk:"f24"`
	F25 string `hunk:"f25"`
}

type benchAllInts struct {
	A int   `hunk:"a"`
	B int64 `hunk:"b"`
	C int32 `hunk:"c"`
	D int16 `hunk:"d"`
	E int   `hunk:"e"`
}

type benchMixed struct {
	I int     `hunk:"i"`
	U uint64  `hunk:"u"`
	F float64 `hunk:"f"`
	S string  `hunk:"s"`
	B bool    `hunk:"b"`
}

type benchTimeHeavy struct {
	T1 time.Time `hunk:"t1"`
	T2 time.Time `hunk:"t2"`
	T3 time.Time `hunk:"t3"`
	T4 time.Time `hunk:"t4"`
	T5 time.Time `hunk:"t5"`
}

// === Benchmark data ===

var (
	fmt1     = ":name"
	entry1   = "Brighton"
	fmt5     = ":a :b :c :d :e"
	entry5   = "alpha beta gamma delta epsilon"
	fmt15    = ":f1 :f2 :f3 :f4 :f5 :f6 :f7 :f8 :f9 :f10 :f11 :f12 :f13 :f14 :f15"
	entry15  = "alpha 42 99999 beta 3.14 gamma 7 delta 12345 epsilon 9876543210 zeta 1024 2.718 omega"
	fmt25    = ":f1 :f2 :f3 :f4 :f5 :f6 :f7 :f8 :f9 :f10 :f11 :f12 :f13 :f14 :f15 :f16 :f17 :f18 :f19 :f20 :f21 :f22 :f23 :f24 :f25"
	entry25  = "alpha beta gamma delta epsilon zeta eta theta iota kappa lambda mu nu xi omicron pi rho sigma tau upsilon phi chi psi omega terminus"
	fmtInts  = ":a :b :c :d :e"
	entInts  = "42 1234567890 9999 100 77"
	fmtMixed = ":i :u :f :s :b"
	entMixed = "42 99999 3.14159 hello true"
	fmtTime  = ":t1 :t2 :t3 :t4 :t5"
	entTime  = "2024-01-15T10:30:00Z 2024-02-20T14:45:30Z 2024-03-25T08:15:00Z 2024-04-01T23:59:59Z 2024-05-10T12:00:00Z"
)

// === Field count benchmarks ===

func BenchmarkParseByFieldCount(b *testing.B) {
	b.Run("1", func(b *testing.B) {
		b.ReportAllocs()
		v := new(benchStr1)
		p, err := NewParser(fmt1, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(entry1, v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("5", func(b *testing.B) {
		b.ReportAllocs()
		v := new(benchStr5)
		p, err := NewParser(fmt5, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(entry5, v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("15", func(b *testing.B) {
		b.ReportAllocs()
		v := new(bench15)
		p, err := NewParser(fmt15, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(entry15, v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("25", func(b *testing.B) {
		b.ReportAllocs()
		v := new(bench25)
		p, err := NewParser(fmt25, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(entry25, v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// === Type complexity benchmarks ===

func BenchmarkParseByTypeComplexity(b *testing.B) {
	b.Run("AllStrings", func(b *testing.B) {
		b.ReportAllocs()
		v := new(benchStr5)
		p, err := NewParser(fmt5, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(entry5, v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("AllInts", func(b *testing.B) {
		b.ReportAllocs()
		v := new(benchAllInts)
		p, err := NewParser(fmtInts, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(entInts, v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Mixed", func(b *testing.B) {
		b.ReportAllocs()
		v := new(benchMixed)
		p, err := NewParser(fmtMixed, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(entMixed, v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("TimeHeavy", func(b *testing.B) {
		b.ReportAllocs()
		v := new(benchTimeHeavy)
		p, err := NewParser(fmtTime, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(entTime, v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// === Line length benchmarks ===

func BenchmarkParseByLineLength(b *testing.B) {
	shortEntry := "hello world foo bar baz"
	mediumEntry := strings.Repeat("a", 40) + " " + strings.Repeat("b", 40) + " " +
		strings.Repeat("c", 40) + " " + strings.Repeat("d", 40) + " " + strings.Repeat("e", 40)
	longEntry := strings.Repeat("a", 200) + " " + strings.Repeat("b", 200) + " " +
		strings.Repeat("c", 200) + " " + strings.Repeat("d", 200) + " " + strings.Repeat("e", 200)

	b.Run("Short_23B", func(b *testing.B) {
		b.ReportAllocs()
		v := new(benchStr5)
		p, err := NewParser(fmt5, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(shortEntry, v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Medium_204B", func(b *testing.B) {
		b.ReportAllocs()
		v := new(benchStr5)
		p, err := NewParser(fmt5, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(mediumEntry, v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Long_1004B", func(b *testing.B) {
		b.ReportAllocs()
		v := new(benchStr5)
		p, err := NewParser(fmt5, v)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := p.ParseLine(longEntry, v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// === Concurrent benchmark ===

func BenchmarkParseConcurrent(b *testing.B) {
	b.ReportAllocs()
	v := new(benchMixed)
	p, err := NewParser(fmtMixed, v)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		local := new(benchMixed)
		for pb.Next() {
			if err := p.ParseLine(entMixed, local); err != nil {
				b.Error(err)
			}
		}
	})
}
