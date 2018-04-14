package hunkee

import (
	"fmt"
	"io"
	"net"
	"net/url"
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

func TestParseStringToStructDuration(t *testing.T) {
	t.Parallel()

	var (
		f = new(field)
		d time.Duration
	)

	if err := parseStringToStruct(reflect.ValueOf(&d), "", f); err == nil {
		t.Errorf("expected %s, got nil error", "some error")
	}

	err := parseStringToStruct(reflect.ValueOf(&d).Elem(), "17s", f)
	if err != nil {
		t.Error(err)
	} else if d.Seconds() != 17 {
		t.Errorf("parsed and source duration are not matched: %s != %d", d, 17)
	}

	err = parseStringToStruct(reflect.ValueOf(&d).Elem(), "20x", f)
	if err == nil {
		t.Error("expected parsing error, got nil")
	}
}

func TestParseStringToStructURL(t *testing.T) {
	t.Parallel()

	var (
		f = &field{
			ftype: typeURL,
		}
		u *url.URL
	)

	if err := parseStringToStruct(reflect.ValueOf(&u), "", f); err == nil {
		t.Errorf("expected %s, got nil error", "some error")
	}

	err := parseStringToStruct(reflect.ValueOf(&u).Elem(), "http://localhost:81/pattern?p1=a&p2=b&p3=c#wow", f)
	if err != nil {
		t.Error(err)
	} else if u.Hostname() != "localhost" || u.Query().Get("p2") != "b" || u.Fragment != "wow" {
		t.Errorf("parsed and source URLs are not matched: %q != %q", "http://localhost:81/pattern?p1=a&p2=b&p3=c#wow", u.String())
	}
}

func TestProcessFieldInt(t *testing.T) {
	t.Parallel()

	type st struct {
		Ui uint64 `hunk:"ui"`
		I  int64  `hunk:"i"`
		Ir string `hunk:"i_raw"`
		B  bool   `hunk:"b"`
		S  string `hunk:"s"`
	}
	s := new(st)

	m, err := initMapper(":ui :i :b :s", s)
	if err != nil {
		t.Error(err)
	}

	var (
		i  int64  = 924
		is string = fmt.Sprintf("%d", i)
	)

	err = m.processField(m.getField("i"), reflect.Indirect(reflect.ValueOf(s)), is)
	if err != nil {
		t.Error(err)
	}
	if s.I != i {
		t.Errorf("int was parsed wrong, expect %d, got %d", i, s.I)
	}
	if s.Ir != is {
		t.Errorf("bad value present in _raw field, expect %s, got %s", is, s.Ir)
	}

	err = m.processField(m.getField("i"), reflect.Indirect(reflect.ValueOf(s)), "924.1")
	if err == nil {
		t.Error("expected parsing error while parse bad int64 value, got nil")
	}
}

func TestProcessFieldUint(t *testing.T) {
	t.Parallel()

	type st struct {
		Ui uint64 `hunk:"ui"`
		I  int64  `hunk:"i"`
		Ir string `hunk:"i_raw"`
		B  bool   `hunk:"b"`
		S  string `hunk:"s"`
	}
	s := new(st)

	m, err := initMapper(":ui :i :b :s", s)
	if err != nil {
		t.Error(err)
	}

	var ui uint64 = 924
	err = m.processField(m.getField("ui"), reflect.Indirect(reflect.ValueOf(s)), fmt.Sprintf("%d", ui))
	if err != nil {
		t.Error(err)
	}
	if s.Ui != ui {
		t.Errorf("uint was parsed wrong, expect %d, got %d", ui, s.I)
	}

	err = m.processField(m.getField("ui"), reflect.Indirect(reflect.ValueOf(s)), "-1")
	if err == nil {
		t.Error("expected parsing error while parse bad uint64 value, got nil")
	}
}

func TestProcessFieldBool(t *testing.T) {
	t.Parallel()

	type st struct {
		B  bool   `hunk:"b"`
		Br string `hunk:"b_raw"`
	}
	s := new(st)

	m, err := initMapper(":b ", s)
	if err != nil {
		t.Error(err)
	}

	err = m.processField(m.getField("b"), reflect.Indirect(reflect.ValueOf(s)), "true")
	if err != nil {
		t.Error(err)
	}
	if !s.B {
		t.Errorf("bool was parsed wrong, expect %t, got %t", true, s.B)
	}
	if s.Br != "true" {
		t.Errorf("bad raw value for bool: %s, want %s", s.Br, "true")
	}

	err = m.processField(m.getField("b"), reflect.Indirect(reflect.ValueOf(s)), "ftrue")
	if err == nil {
		t.Error("expected parsing error while parse bad boolean value, got nil")
	}
}

func TestProcessFieldFloat(t *testing.T) {
	t.Parallel()

	type st struct {
		F  float32 `hunk:"f"`
		Fr string  `hunk:"f_raw"`
	}
	s := new(st)

	m, err := initMapper(":f", s)
	if err != nil {
		t.Error(err)
	}

	var f float32 = 32.1684
	err = m.processField(m.getField("f"), reflect.Indirect(reflect.ValueOf(s)), fmt.Sprintf("%f", f))
	if err != nil {
		t.Error(err)
	}
	if s.F != f {
		t.Errorf("bool was parsed wrong, expect %f, got %f", 32.1684, s.F)
	}

	err = m.processField(m.getField("f"), reflect.Indirect(reflect.ValueOf(s)), "32.123.1")
	if err == nil {
		t.Error("expected parsing error while parse bad float32 value, got nil")
	}
}

func TestProcessFieldString(t *testing.T) {
	t.Parallel()

	type st struct {
		S string `hunk:"s"`
	}
	s := new(st)

	m, err := initMapper(":s", s)
	if err != nil {
		t.Error(err)
	}

	msg := "some wow message here"
	err = m.processField(m.getField("s"), reflect.Indirect(reflect.ValueOf(s)), msg)
	if err != nil {
		t.Error(err)
	}
	if s.S != msg {
		t.Errorf("string was parsed wrong, expect %s, got %s", msg, s.S)
	}
}

func TestProcessFieldNotSupported(t *testing.T) {
	t.Parallel()

	type st struct {
		W io.Writer `hunk:"s"`
	}
	s := new(st)

	m, err := initMapper(":s", s)
	if err != nil {
		t.Error(err)
	}

	w := "172.17.254.1"
	err = m.processField(m.getField("s"), reflect.Indirect(reflect.ValueOf(s)), w)
	if err == nil {
		t.Error("expected error of not supported, got nil")
	}
}

func TestProcessFieldIP(t *testing.T) {
	t.Parallel()

	type st struct {
		IP     net.IP `hunk:"ip"`
		IPrawr string `hunk:"ip_raw"`
	}
	s := new(st)

	m, err := initMapper(":ip", s)
	if err != nil {
		t.Error(err)
	}

	ip := "241.74.91.45"
	err = m.processField(m.getField("ip"), reflect.Indirect(reflect.ValueOf(s)), ip)
	if err != nil {
		t.Error(err)
	}
	if !s.IP.Equal(net.IPv4(241, 74, 91, 45)) {
		t.Errorf("net.IP was parsed wrong, expect %q, got %q", ip, s.IP)
	}
	if s.IPrawr != ip {
		t.Errorf("net.IP raw was not stored, expect %q, got %q", ip, s.IPrawr)
	}
}

func TestProcessTag(t *testing.T) {
	tag := reflect.StructTag("")
	_, _, err := processTag(tag)
	if err != nil {
		t.Error(err)
	}

	tag = reflect.StructTag(`hunk:""`)
	_, _, err = processTag(tag)
	if err != nil {
		t.Error(err)
	}

	tag = reflect.StructTag(`hunk:"alia,s"`)
	_, _, err = processTag(tag)
	if err != ErrComaNotSupported {
		t.Errorf("expected %s error, got %v", ErrComaNotSupported, err)
	}
}
