package hunkee

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// processField gets token and parse it into corresponded type and puts into 'final' value
func (m *mapper) processField(field *field, final reflect.Value, token string) error {
	v := final.FieldByIndex(field.index)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		// Allocate memory
		v.Set(reflect.New(deref(v.Type())))
	}

	// set raw value
	if field.hasRaw {
		raw := m.raw(field)
		if raw == nil {
			panic(fmt.Sprintf("%s field should have raw field, but it's not provided", field.name))
		}
		final.Field(raw.index[0]).Set(reflect.ValueOf(token))
	}

	// nothing to process, but if it's string, we should set token to the field
	if token == "-" {
		if field.reflectKind == reflect.String {
			v.SetString(token)
		}
		return nil
	}

	switch field.reflectKind {
	case reflect.Bool:
		b, err := strconv.ParseBool(token)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64, err := parseInt(field.reflectKind, token)
		if err != nil {
			return err
		}
		v.SetInt(i64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ui64, err := parseUint(field.reflectKind, token)
		if err != nil {
			return err
		}
		v.SetUint(ui64)
	case reflect.String:
		v.SetString(token)
	case reflect.Float32, reflect.Float64:
		fl64, err := parseFloat(field.reflectKind, token)
		if err != nil {
			return err
		}
		v.SetFloat(fl64)
	case reflect.Struct:
		if err := parseStringToStruct(v, token, field); err != nil {
			return err
		}
	default:
		if field.ftype == typeIP {
			ip := net.ParseIP(token)
			v.Set(reflect.ValueOf(ip))
		} else {
			return fmt.Errorf("type %+v is not supported", field)
		}
	}

	return nil
}

func parseUint(kind reflect.Kind, token string) (uint64, error) {
	var size int
	switch kind {
	case reflect.Uint8:
		size = 8
	case reflect.Uint16:
		size = 16
	case reflect.Uint32:
		size = 32
	case reflect.Uint64, reflect.Uint:
		size = 64
	default:
		return 0, ErrNotUint
	}

	return strconv.ParseUint(token, 10, size)
}

func parseInt(kind reflect.Kind, token string) (int64, error) {
	var size int
	switch kind {
	case reflect.Int8:
		size = 8
	case reflect.Int16:
		size = 16
	case reflect.Int32:
		size = 32
	case reflect.Int, reflect.Int64:
		size = 64
	default:
		return 0, ErrNotInt
	}

	return strconv.ParseInt(token, 10, size)
}

func parseFloat(kind reflect.Kind, token string) (float64, error) {
	var size int
	switch kind {
	case reflect.Float32:
		size = 32
	case reflect.Float64:
		size = 64
	default:
		return 0, ErrNotFloat
	}

	return strconv.ParseFloat(token, size)
}

// parseStringToStruct gets token and parses it into
// net.Addr, time.Time, time.Duration, url.URL
func parseStringToStruct(v reflect.Value, token string, field *field) (err error) {
	switch field.ftype {
	case typeTime:
		var t time.Time
		if field.timeOptions == nil {
			return ErrNilTimeOptions
		}

		if field.timeOptions.Location == nil {
			t, err = time.Parse(field.timeOptions.Layout, token)
		} else {
			t, err = time.ParseInLocation(field.timeOptions.Layout, token, field.timeOptions.Location)
		}
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(t))
	case typeURL:
		u, err := url.Parse(token)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(u))
	case typeDuration:
		d, err := time.ParseDuration(token)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(d))
	default:
		return fmt.Errorf("type %s is not supported: ftype: %d name: %s", field.reflectType, field.ftype, field.name)
	}
	return nil
}

// processTag returns full tag, normalName aka not raw name and error, if exists
func processTag(tagLine reflect.StructTag) (tag, normalName string, err error) {
	var ok bool
	tag, ok = tagLine.Lookup(libtag)
	if !ok || tag == "" || tag == unexportedTag {
		tag = unexportedTag
		return
	}

	if strings.ContainsAny(tag, ".,") {
		err = ErrComaNotSupported
		return
	}

	if strings.HasSuffix(tag, "_raw") {
		normalName = strings.TrimSuffix(tag, "_raw")
	}

	return
}
