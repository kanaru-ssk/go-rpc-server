package query

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

// Decode は url.Values を構造体へマッピングする。
// 構造体のフィールドに `query:"name"` タグを指定して利用する。
// プリミティブ型とtime.Time型(RFC3339, DateOnly, TimeOnly)と、それらの配列をサポート
func Decode(values url.Values, dst any) error {
	if values == nil {
		values = url.Values{}
	}
	if dst == nil {
		return fmt.Errorf("query.Decode: dst must not be nil")
	}
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("query.Decode: dst must be non-nil pointer to struct")
	}
	structVal := v.Elem()
	if structVal.Kind() != reflect.Struct {
		return fmt.Errorf("query.Decode: dst must point to struct")
	}
	structType := structVal.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("query")
		if tag == "" || !field.IsExported() {
			continue
		}
		vals, ok := values[tag]
		if !ok || len(vals) == 0 {
			continue
		}
		if err := assignValue(structVal.Field(i), vals); err != nil {
			return fmt.Errorf("query.Decode: field %s: %w", field.Name, err)
		}
	}
	return nil
}

func assignValue(field reflect.Value, values []string) error {
	if !field.CanSet() {
		return fmt.Errorf("query.assignValue: cannot set field")
	}
	if field.Kind() == reflect.Slice {
		return assignSlice(field, values)
	}
	parsed, err := parseScalar(field.Type(), values[len(values)-1])
	if err != nil {
		return fmt.Errorf("query.assignValue: %w", err)
	}
	field.Set(parsed)
	return nil
}

func parseTime(value string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.DateOnly, value); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.TimeOnly, value); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("query.parseTime: : invalid time value %q", value)
}

func assignSlice(field reflect.Value, values []string) error {
	elemType := field.Type().Elem()
	slice := reflect.MakeSlice(field.Type(), 0, len(values))
	for _, v := range values {
		parsed, err := parseScalar(elemType, v)
		if err != nil {
			return fmt.Errorf("query.assignSlice: %w", err)
		}
		slice = reflect.Append(slice, parsed)
	}
	field.Set(slice)
	return nil
}

func parseScalar(t reflect.Type, value string) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(value).Convert(t), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize := t.Bits()
		n, err := strconv.ParseInt(value, 10, bitSize)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("query.parseScalar: %w", err)
		}
		val := reflect.New(t).Elem()
		val.SetInt(n)
		return val, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		bitSize := t.Bits()
		n, err := strconv.ParseUint(value, 10, bitSize)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("query.parseScalar: %w", err)
		}
		val := reflect.New(t).Elem()
		val.SetUint(n)
		return val, nil
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("query.parseScalar: %w", err)
		}
		val := reflect.New(t).Elem()
		val.SetBool(b)
		return val, nil
	case reflect.Float32, reflect.Float64:
		bitSize := t.Bits()
		f, err := strconv.ParseFloat(value, bitSize)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("query.parseScalar: %w", err)
		}
		val := reflect.New(t).Elem()
		val.SetFloat(f)
		return val, nil
	case reflect.Complex64, reflect.Complex128:
		bitSize := t.Bits()
		c, err := strconv.ParseComplex(value, bitSize)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("query.parseScalar: %w", err)
		}
		val := reflect.New(t).Elem()
		val.SetComplex(c)
		return val, nil
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			tm, err := parseTime(value)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("query.parseScalar: %w", err)
			}
			return reflect.ValueOf(tm), nil
		}
	}
	return reflect.Value{}, fmt.Errorf("query.parseScalar: unsupported kind %s", t.Kind())
}
