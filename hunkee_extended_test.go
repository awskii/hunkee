package hunkee

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
)

// === 1. Edge Cases & Robustness Tests ===

// TestParseLine_DashAsNull_IntField verifies that "-" token into an int field
// produces zero value with no error, matching nginx/CDN log conventions where
// dash represents missing values.
func TestParseLine_DashAsNull_IntField(t *testing.T) {
	var s struct {
		Status int    `hunk:"status"`
		Name   string `hunk:"name"`
	}
	p, err := NewParser(":status :name", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	if err := p.ParseLine("- hello", &s); err != nil {
		t.Fatalf("unexpected parse error: %s", err)
	}
	if s.Status != 0 {
		t.Errorf("int with dash: expected 0, got %d", s.Status)
	}
	if s.Name != "hello" {
		t.Errorf("name: expected 'hello', got %q", s.Name)
	}
}

// TestParseLine_DashAsNull_MultipleTypes verifies that "-" into uint, float, bool,
// IP, and string fields produces appropriate zero/default values.
func TestParseLine_DashAsNull_MultipleTypes(t *testing.T) {
	var s struct {
		U  uint64  `hunk:"u"`
		F  float64 `hunk:"f"`
		B  bool    `hunk:"b"`
		IP net.IP  `hunk:"ip"`
		S  string  `hunk:"s"`
	}
	p, err := NewParser(":u :f :b :ip :s", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	if err := p.ParseLine("- - - - -", &s); err != nil {
		t.Fatalf("unexpected parse error: %s", err)
	}
	if s.U != 0 {
		t.Errorf("uint: expected 0 for dash, got %d", s.U)
	}
	if s.F != 0 {
		t.Errorf("float: expected 0 for dash, got %f", s.F)
	}
	if s.B != false {
		t.Errorf("bool: expected false for dash, got %t", s.B)
	}
	if s.IP != nil {
		t.Errorf("IP: expected nil for dash, got %v", s.IP)
	}
	// String fields get "-" literally, not zero value
	if s.S != "-" {
		t.Errorf("string: expected '-' for dash, got %q", s.S)
	}
}

// TestParseLine_EmptyToken_NumericTypes verifies that empty tokens for int/uint/float
// produce zero values with no error (guarded by `token != ""` check in processField).
func TestParseLine_EmptyToken_NumericTypes(t *testing.T) {
	var s struct {
		I int     `hunk:"i"`
		U uint64  `hunk:"u"`
		F float64 `hunk:"f"`
	}
	p, err := NewParser(":i :u :f", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	p.SetTokenSeparator('"')

	// All tokens are empty strings between quotes
	if err := p.ParseLine(`"" "" ""`, &s); err != nil {
		t.Fatalf("unexpected parse error for empty tokens: %s", err)
	}
	if s.I != 0 {
		t.Errorf("int: expected 0 for empty token, got %d", s.I)
	}
	if s.U != 0 {
		t.Errorf("uint: expected 0 for empty token, got %d", s.U)
	}
	if s.F != 0 {
		t.Errorf("float: expected 0 for empty token, got %f", s.F)
	}
}

// TestParseLine_EmptyToken_BoolError verifies that an empty token for a bool
// field produces an error (no `token != ""` guard exists for bool in processField).
func TestParseLine_EmptyToken_BoolError(t *testing.T) {
	var s struct {
		B bool `hunk:"b"`
	}
	p, err := NewParser(":b", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	p.SetTokenSeparator('"')

	err = p.ParseLine(`""`, &s)
	if err == nil {
		t.Fatal("expected error for empty bool token, got nil")
	}
}

// TestParseLine_Overflow_Uint8 verifies that parsing "256" into a uint8 field
// returns a range error (max uint8 is 255).
func TestParseLine_Overflow_Uint8(t *testing.T) {
	var s struct {
		V uint8 `hunk:"v"`
	}
	p, err := NewParser(":v", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	err = p.ParseLine("256", &s)
	if err == nil {
		t.Fatal("expected range error for uint8=256, got nil")
	}
	if !strings.Contains(err.Error(), "range") {
		t.Errorf("expected range error, got: %s", err)
	}
}

// TestParseLine_Overflow_Int8 verifies that parsing "128" into an int8 field
// returns a range error (max int8 is 127).
func TestParseLine_Overflow_Int8(t *testing.T) {
	var s struct {
		V int8 `hunk:"v"`
	}
	p, err := NewParser(":v", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	err = p.ParseLine("128", &s)
	if err == nil {
		t.Fatal("expected range error for int8=128, got nil")
	}
	if !strings.Contains(err.Error(), "range") {
		t.Errorf("expected range error, got: %s", err)
	}
}

// TestParseLine_VeryLongString verifies that a 10KB+ token does not cause
// panics or OOM.
func TestParseLine_VeryLongString(t *testing.T) {
	var s struct {
		Name string `hunk:"name"`
	}
	p, err := NewParser(":name", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	longVal := strings.Repeat("x", 10240) // 10KB
	if err := p.ParseLine(longVal, &s); err != nil {
		t.Fatalf("unexpected error for long string: %s", err)
	}
	if len(s.Name) != 10240 {
		t.Errorf("expected length 10240, got %d", len(s.Name))
	}
}

// TestParseLine_MultipleConsecutiveSpaces documents that consecutive spaces
// cause empty tokens for intermediate fields (the parser treats each space
// as a separator boundary).
func TestParseLine_MultipleConsecutiveSpaces(t *testing.T) {
	var s struct {
		A string `hunk:"a"`
		B string `hunk:"b"`
		C string `hunk:"c"`
	}
	p, err := NewParser(":a :b :c", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}

	// Double space between "hello" and "world" means field B gets an empty
	// token from the adjacent space characters.
	if err := p.ParseLine("hello  world", &s); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if s.A != "hello" {
		t.Errorf("field A: expected 'hello', got %q", s.A)
	}
	if s.B != "" {
		t.Errorf("field B: expected empty string (double space), got %q", s.B)
	}
	if s.C != "world" {
		t.Errorf("field C: expected 'world', got %q", s.C)
	}
}

// === 2. Token Count Mismatch Tests ===

// TestParseLine_FewerTokensThanFormat verifies that a line with fewer tokens
// than the format expects returns an error instead of panicking.
// Bug fix: prior to the offset<0 guard in extractors.go, this would panic
// with a negative slice index.
func TestParseLine_FewerTokensThanFormat(t *testing.T) {
	var s struct {
		A string `hunk:"a"`
		B string `hunk:"b"`
		C string `hunk:"c"`
	}
	p, err := NewParser(":a :b :c", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}

	err = p.ParseLine("hello world", &s)
	if err == nil {
		t.Fatal("expected error for fewer tokens than format, got nil")
	}
	if !strings.Contains(err.Error(), "less tokens") {
		t.Errorf("expected 'less tokens' error, got: %s", err)
	}
}

// TestParseLine_MoreTokensThanFormat verifies that extra tokens in the line
// are silently ignored — only the fields defined in the format are parsed.
func TestParseLine_MoreTokensThanFormat(t *testing.T) {
	var s struct {
		A string `hunk:"a"`
		B string `hunk:"b"`
	}
	p, err := NewParser(":a :b", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}

	err = p.ParseLine("hello world extra tokens here", &s)
	if err != nil {
		t.Fatalf("expected no error for extra tokens, got: %s", err)
	}
	if s.A != "hello" {
		t.Errorf("expected 'hello', got %q", s.A)
	}
	if s.B != "world" {
		t.Errorf("expected 'world', got %q", s.B)
	}
}

// TestParseLine_FewerTokensWithSeparator verifies that fewer tokens with a
// custom separator also returns an error (not a panic).
func TestParseLine_FewerTokensWithSeparator(t *testing.T) {
	var s struct {
		A string `hunk:"a"`
		B string `hunk:"b"`
		C string `hunk:"c"`
	}
	p, err := NewParser(":a :b :c", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	p.SetTokenSeparator('"')

	err = p.ParseLine(`"hello" "world"`, &s)
	if err == nil {
		t.Fatal("expected error for fewer tokens with separator, got nil")
	}
	if !strings.Contains(err.Error(), "less tokens") {
		t.Errorf("expected 'less tokens' error, got: %s", err)
	}
}

// === 3. Separator Variety Tests ===

// TestParseLine_PipeSeparator verifies that pipe '|' works as a wrapping
// separator, allowing values with spaces inside.
func TestParseLine_PipeSeparator(t *testing.T) {
	var s struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}
	p, err := NewParser(":id :name", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	p.SetTokenSeparator('|')

	if err := p.ParseLine("|42| |John Doe|", &s); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if s.ID != 42 {
		t.Errorf("expected ID 42, got %d", s.ID)
	}
	if s.Name != "John Doe" {
		t.Errorf("expected 'John Doe', got %q", s.Name)
	}
}

// TestParseLine_TabSeparator verifies that tab can be used as a wrapping
// separator. Note: hunkee treats the separator as a wrapping character
// (like quotes), not as a field delimiter (like TSV/CSV).
func TestParseLine_TabSeparator(t *testing.T) {
	var s struct {
		A string `hunk:"a"`
		B string `hunk:"b"`
	}
	p, err := NewParser(":a :b", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	p.SetTokenSeparator('\t')

	// Values wrapped in tabs, separated by space
	if err := p.ParseLine("\thello\t \tworld\t", &s); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if s.A != "hello" {
		t.Errorf("expected 'hello', got %q", s.A)
	}
	if s.B != "world" {
		t.Errorf("expected 'world', got %q", s.B)
	}
}

// TestParseLine_SeparatorInsideQuotedValues verifies that the default separator
// (space) appearing inside quote-wrapped tokens is preserved as part of the value.
func TestParseLine_SeparatorInsideQuotedValues(t *testing.T) {
	var s struct {
		Method string `hunk:"method"`
		Path   string `hunk:"path"`
		Agent  string `hunk:"agent"`
	}
	p, err := NewParser(":method :path :agent", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	p.SetTokenSeparator('"')

	if err := p.ParseLine(`"GET" "/index.html" "Mozilla/5.0 (X11; Linux)"`, &s); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if s.Method != "GET" {
		t.Errorf("expected 'GET', got %q", s.Method)
	}
	if s.Path != "/index.html" {
		t.Errorf("expected '/index.html', got %q", s.Path)
	}
	if s.Agent != "Mozilla/5.0 (X11; Linux)" {
		t.Errorf("expected full user agent string, got %q", s.Agent)
	}
}

// === 4. Concurrency Tests ===

// TestParseLine_ConcurrentAccess spawns multiple goroutines calling ParseLine
// on the same Parser instance with individual structs. Run with -race to detect
// data races in the worker pool.
func TestParseLine_ConcurrentAccess(t *testing.T) {
	type entry struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}
	p, err := NewParser(":id :name", &entry{})
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	const numGoroutines = 100
	var wg sync.WaitGroup
	errs := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var local entry
			line := fmt.Sprintf("%d name_%d", id, id)
			if err := p.ParseLine(line, &local); err != nil {
				errs <- fmt.Errorf("goroutine %d: parse error: %s", id, err)
				return
			}
			if local.ID != id {
				errs <- fmt.Errorf("goroutine %d: expected ID %d, got %d", id, id, local.ID)
			}
			expectedName := fmt.Sprintf("name_%d", id)
			if local.Name != expectedName {
				errs <- fmt.Errorf("goroutine %d: expected name %q, got %q", id, expectedName, local.Name)
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

// TestParseLine_ConcurrentMixedTypes verifies concurrent parsing with multiple
// field types doesn't corrupt data across goroutines.
func TestParseLine_ConcurrentMixedTypes(t *testing.T) {
	type record struct {
		ID     uint32  `hunk:"id"`
		Score  float64 `hunk:"score"`
		Name   string  `hunk:"name"`
		Active bool    `hunk:"active"`
	}
	p, err := NewParser(":id :score :name :active", &record{})
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}
	const N = 200
	var wg sync.WaitGroup
	errs := make(chan error, N)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var r record
			line := fmt.Sprintf("%d %d.%d item_%d true", id, id, id, id)
			if err := p.ParseLine(line, &r); err != nil {
				errs <- fmt.Errorf("goroutine %d: %s", id, err)
				return
			}
			if r.ID != uint32(id) {
				errs <- fmt.Errorf("goroutine %d: ID mismatch: expected %d, got %d", id, id, r.ID)
			}
			if !r.Active {
				errs <- fmt.Errorf("goroutine %d: Active should be true", id)
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

// === 5. Parser Reuse Tests ===

// TestParseLine_ParserReuse_NoStateLeak parses 1500 different lines with one
// Parser instance and verifies no state leaks between calls.
func TestParseLine_ParserReuse_NoStateLeak(t *testing.T) {
	type entry struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}
	p, err := NewParser(":id :name", &entry{})
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}

	for i := 0; i < 1500; i++ {
		var s entry
		expectedName := fmt.Sprintf("item_%d", i)
		line := fmt.Sprintf("%d %s", i, expectedName)
		if err := p.ParseLine(line, &s); err != nil {
			t.Fatalf("iteration %d: parse error: %s", i, err)
		}
		if s.ID != i {
			t.Fatalf("iteration %d: expected ID %d, got %d", i, i, s.ID)
		}
		if s.Name != expectedName {
			t.Fatalf("iteration %d: expected name %q, got %q", i, expectedName, s.Name)
		}
	}
}

// TestParseLine_StructReuse_CleanOverwrite parses different lines into the
// same struct and verifies previous values are cleanly overwritten each time.
func TestParseLine_StructReuse_CleanOverwrite(t *testing.T) {
	var s struct {
		ID   int    `hunk:"id"`
		Name string `hunk:"name"`
	}
	p, err := NewParser(":id :name", &s)
	if err != nil {
		t.Fatalf("unexpected init error: %s", err)
	}

	// First parse
	if err := p.ParseLine("100 alpha", &s); err != nil {
		t.Fatalf("first parse error: %s", err)
	}
	if s.ID != 100 || s.Name != "alpha" {
		t.Fatalf("first parse: expected {100, alpha}, got {%d, %s}", s.ID, s.Name)
	}

	// Second parse into same struct — should overwrite
	if err := p.ParseLine("200 beta", &s); err != nil {
		t.Fatalf("second parse error: %s", err)
	}
	if s.ID != 200 {
		t.Errorf("struct reuse: ID not overwritten, expected 200, got %d", s.ID)
	}
	if s.Name != "beta" {
		t.Errorf("struct reuse: Name not overwritten, expected 'beta', got %q", s.Name)
	}

	// Third parse with zero value for ID — verify it doesn't retain old value
	if err := p.ParseLine("0 gamma", &s); err != nil {
		t.Fatalf("third parse error: %s", err)
	}
	if s.ID != 0 {
		t.Errorf("struct reuse: ID not overwritten to 0, got %d", s.ID)
	}
	if s.Name != "gamma" {
		t.Errorf("struct reuse: Name not overwritten, expected 'gamma', got %q", s.Name)
	}
}
