package hunkee

import (
	"reflect"
	"testing"
)

func TestParseUint(t *testing.T) {
	t.Parallel()

	ts := []string{"0", "256", "65536", "4294967297"}
	te := []uint64{0, 256, 65536, 4294967297}

	r, err := parseUint(reflect.Uint8, ts[0])
	if err != nil {
		t.Error(err)
	} else if r != te[0] {
		t.Errorf("expected %d, got %d", te[0], r)
	}

	r, err = parseUint(reflect.Uint16, ts[1])
	if err != nil {
		t.Error(err)
	} else if r != te[1] {
		t.Errorf("expected %d, got %d", te[1], r)
	}

	r, err = parseUint(reflect.Uint32, ts[2])
	if err != nil {
		t.Error(err)
	} else if r != te[2] {
		t.Errorf("expected %d, got %d", te[2], r)
	}

	// uint should be at least 32 bits
	r, err = parseUint(reflect.Uint, ts[2])
	if err != nil {
		t.Error(err)
	} else if r != te[2] {
		t.Errorf("expected %d, got %d", te[2], r)
	}

	// so it should works perfectly with 64 bit number too
	r, err = parseUint(reflect.Uint, ts[3])
	if err != nil {
		t.Error(err)
	} else if r != te[3] {
		t.Errorf("expected %d, got %d", te[3], r)
	}

	r, err = parseUint(reflect.Uint64, ts[3])
	if err != nil {
		t.Error(err)
	} else if r != te[3] {
		t.Errorf("expected %d, got %d", te[3], r)
	}

	// tests on bad inputs
	_, err = parseUint(reflect.Int64, ts[3])
	if err == nil {
		t.Errorf("expected %q, got nil error", ErrNotUint)
	}

	_, err = parseUint(reflect.Uint, "y321")
	if err == nil {
		t.Errorf("expected %q, got nil error", ErrNotUint)
	}
}
