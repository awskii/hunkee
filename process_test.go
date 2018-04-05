package hunkee

import (
	"reflect"
	"testing"
	"time"
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

func TestParseInt(t *testing.T) {
	t.Parallel()

	ts := []string{"-127", "32767", "-2147483647", "-9223372036854775800"}
	te := []int64{-127, 32767, -2147483647, -9223372036854775800}

	r, err := parseInt(reflect.Int8, ts[0])
	if err != nil {
		t.Error(err)
	} else if r != te[0] {
		t.Errorf("expected %d, got %d", te[0], r)
	}

	r, err = parseInt(reflect.Int16, ts[1])
	if err != nil {
		t.Error(err)
	} else if r != te[1] {
		t.Errorf("expected %d, got %d", te[1], r)
	}

	r, err = parseInt(reflect.Int32, ts[2])
	if err != nil {
		t.Error(err)
	} else if r != te[2] {
		t.Errorf("expected %d, got %d", te[2], r)
	}

	// int should be at least 32 bits
	r, err = parseInt(reflect.Int, ts[2])
	if err != nil {
		t.Error(err)
	} else if r != te[2] {
		t.Errorf("expected %d, got %d", te[2], r)
	}

	// so it should works perfectly with 64 bit number too
	r, err = parseInt(reflect.Int, ts[3])
	if err != nil {
		t.Error(err)
	} else if r != te[3] {
		t.Errorf("expected %d, got %d", te[3], r)
	}

	r, err = parseInt(reflect.Int64, ts[3])
	if err != nil {
		t.Error(err)
	} else if r != te[3] {
		t.Errorf("expected %d, got %d", te[3], r)
	}

	// tests on bad inputs
	_, err = parseInt(reflect.Uint, ts[3])
	if err == nil {
		t.Errorf("expected %q, got nil error", ErrNotInt)
	}

	_, err = parseInt(reflect.Int, "32y1")
	if err == nil {
		t.Errorf("expected %q, got nil error", ErrNotInt)
	}
}

func TestParseFloat(t *testing.T) {
	t.Parallel()

	ts := []string{"1234212.2916", "32767.999999999999999999999991", "2147483647.0", "-9223372036854775800.1532164436244444"}
	te := []float64{1234212.2916, 32767.999999999999999999999991, 2147483647.0, -9223372036854775800.1532164436244444}

	r, err := parseFloat(reflect.Float32, ts[0])
	if err != nil {
		t.Error(err)
	} else if (r - te[0]) > 0.5 {
		t.Errorf("expected %f, got %f Maximum difference of 0.5 is exceeded", te[0], r)
	}

	r, err = parseFloat(reflect.Float32, ts[1])
	if err != nil {
		t.Error(err)
	} else if (r - te[1]) > 0.5 {
		t.Errorf("expected %f, got %f Maximum difference of 0.5 is exceeded", te[1], r)
	}

	r, err = parseFloat(reflect.Float64, ts[2])
	if err != nil {
		t.Error(err)
	} else if r != te[2] {
		t.Errorf("expected %f, got %f", te[2], r)
	}

	r, err = parseFloat(reflect.Float64, ts[3])
	if err != nil {
		t.Error(err)
	} else if r != te[3] {
		t.Errorf("expected %f, got %f", te[3], r)
	}

	// tests on bad inputs
	_, err = parseFloat(reflect.Uint, ts[3])
	if err == nil {
		t.Errorf("expected %q, got nil error", ErrNotFloat)
	}

	_, err = parseFloat(reflect.Float32, "32y1")
	if err == nil {
		t.Errorf("expected %q, got nil error", ErrNotFloat)
	}
}

func TestParseStringToStructTime(t *testing.T) {
	t.Parallel()

	var (
		f   = new(field)
		tim time.Time
	)

	if err := parseStringToStruct(reflect.ValueOf(&tim), "", f); err == nil {
		t.Errorf("expected %s, got nil error", ErrNilTimeOptions)
	}

	f.timeOptions = &TimeOption{
		Layout: time.ANSIC,
	}

	now := time.Now()
	err := parseStringToStruct(reflect.ValueOf(&tim).Elem(), now.Format(time.ANSIC), f)
	if err != nil {
		t.Error(err)
	} else if tim.Format(time.ANSIC) != now.Format(time.ANSIC) {
		t.Errorf("parsed and source time not matched: %s != %s", now, tim)
	}

	now = time.Now()
	err = parseStringToStruct(reflect.ValueOf(&tim).Elem(), "ggwp2018", f)
	if err == nil {
		t.Error("expected parsing error, got nil")
	}

	f.timeOptions.Location, _ = time.LoadLocation("Local")

	now = time.Now()
	err = parseStringToStruct(reflect.ValueOf(&tim).Elem(), now.Format(time.ANSIC), f)
	if err != nil {
		t.Error(err)
	} else if tim.Format(time.ANSIC) != now.Format(time.ANSIC) {
		t.Errorf("parsed and source time not matched: %s != %s", now, tim)
	}

	now = time.Now()
	err = parseStringToStruct(reflect.ValueOf(&tim).Elem(), "wp2018gg", f)
	if err == nil {
		t.Error(err)
	}

}
