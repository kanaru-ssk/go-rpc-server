package query

import (
	"net/url"
	"testing"
	"time"
)

func TestDecode_SetsSupportedTypes(t *testing.T) {
	values := url.Values{
		"string":         {"abc"},
		"strings":        {"a", "b"},
		"int":            {"10"},
		"ints":           {"0", "1", "10"},
		"int8":           {"-8"},
		"int16":          {"-16"},
		"int32":          {"-32"},
		"int64":          {"-64"},
		"uint":           {"12"},
		"uints":          {"5", "6"},
		"uint8":          {"8"},
		"uint16":         {"16"},
		"uint32":         {"32"},
		"uint64":         {"64"},
		"bool":           {"true"},
		"bools":          {"true", "false", "true"},
		"float32":        {"3.14"},
		"float64":        {"2.718"},
		"floats":         {"1.23", "4.56"},
		"complex64":      {"(1+2i)"},
		"complex128":     {"(3+4i)"},
		"complexes":      {"(1+0i)", "(0+1i)"},
		"datetime":       {"2025-01-02T15:04:05Z"},
		"datetimeOffset": {"2025-01-02T15:04:05+09:00"},
		"date":           {"2025-01-02"},
		"time":           {"15:04:05"},
		"times":          {"2025-01-02T15:04:05Z", "2025-01-03"},
	}

	var dst struct {
		String         string       `query:"string"`
		Strings        []string     `query:"strings"`
		Int            int          `query:"int"`
		Ints           []int        `query:"ints"`
		Int8           int8         `query:"int8"`
		Int16          int16        `query:"int16"`
		Int32          int32        `query:"int32"`
		Int64          int64        `query:"int64"`
		Uint           uint         `query:"uint"`
		UintSlice      []uint       `query:"uints"`
		Uint8          uint8        `query:"uint8"`
		Uint16         uint16       `query:"uint16"`
		Uint32         uint32       `query:"uint32"`
		Uint64         uint64       `query:"uint64"`
		Bool           bool         `query:"bool"`
		BoolSlice      []bool       `query:"bools"`
		Float32        float32      `query:"float32"`
		Float64        float64      `query:"float64"`
		FloatSlice     []float64    `query:"floats"`
		Complex64      complex64    `query:"complex64"`
		Complex128     complex128   `query:"complex128"`
		ComplexSlice   []complex128 `query:"complexes"`
		Datetime       time.Time    `query:"datetime"`
		DatetimeOffset time.Time    `query:"datetimeOffset"`
		Date           time.Time    `query:"date"`
		Time           time.Time    `query:"time"`
		TimeSlice      []time.Time  `query:"times"`
	}

	if err := Decode(values, &dst); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if dst.String != "abc" {
		t.Fatalf("string mismatch: %+v", dst)
	}
	if len(dst.Strings) != 2 || dst.Strings[0] != "a" || dst.Strings[1] != "b" {
		t.Fatalf("strings mismatch: %+v", dst)
	}
	if dst.Int != 10 || dst.Int8 != -8 || dst.Int16 != -16 || dst.Int32 != -32 || dst.Int64 != -64 {
		t.Fatalf("int mismatch: %+v", dst)
	}

	if dst.Uint != 12 || dst.Uint8 != 8 || dst.Uint16 != 16 || dst.Uint32 != 32 || dst.Uint64 != 64 {
		t.Fatalf("uint mismatch: %+v", dst)
	}
	if !dst.Bool {
		t.Fatalf("bool mismatch: %+v", dst)
	}
	if dst.Float32 != 3.14 || dst.Float64 != 2.718 {
		t.Fatalf("float mismatch: %+v", dst)
	}
	if dst.Complex64 != complex64(1+2i) || dst.Complex128 != complex128(3+4i) {
		t.Fatalf("complex mismatch: %+v", dst)
	}
	if dst.Datetime.Format(time.RFC3339) != "2025-01-02T15:04:05Z" {
		t.Fatalf("datetime mismatch: %+v", dst)
	}
	if dst.DatetimeOffset.UTC().Format(time.RFC3339) != "2025-01-02T06:04:05Z" {
		t.Fatalf("datetime mismatch: %+v", dst)
	}
	if dst.Date.Format(time.DateOnly) != "2025-01-02" {
		t.Fatalf("date mismatch: %+v", dst)
	}
	if dst.Time.Format(time.TimeOnly) != "15:04:05" {
		t.Fatalf("time mismatch: %+v", dst)
	}
	if len(dst.Ints) != 3 || dst.Ints[0] != 0 || dst.Ints[1] != 1 || dst.Ints[2] != 10 {
		t.Fatalf("ints slice mismatch: %+v", dst.Ints)
	}
	if len(dst.UintSlice) != 2 || dst.UintSlice[0] != 5 || dst.UintSlice[1] != 6 {
		t.Fatalf("uints slice mismatch: %+v", dst.UintSlice)
	}
	if len(dst.BoolSlice) != 3 || !dst.BoolSlice[0] || dst.BoolSlice[1] || !dst.BoolSlice[2] {
		t.Fatalf("bools slice mismatch: %+v", dst.BoolSlice)
	}
	if len(dst.FloatSlice) != 2 || dst.FloatSlice[0] != 1.23 || dst.FloatSlice[1] != 4.56 {
		t.Fatalf("float slice mismatch: %+v", dst.FloatSlice)
	}
	if len(dst.ComplexSlice) != 2 || dst.ComplexSlice[0] != complex128(1+0i) || dst.ComplexSlice[1] != complex128(0+1i) {
		t.Fatalf("complex slice mismatch: %+v", dst.ComplexSlice)
	}
	if len(dst.TimeSlice) != 2 || dst.TimeSlice[0].Format(time.RFC3339) != "2025-01-02T15:04:05Z" || dst.TimeSlice[1].Format(time.DateOnly) != "2025-01-03" {
		t.Fatalf("time slice mismatch: %+v", dst.TimeSlice)
	}
}

func TestDecode_InvalidDestination(t *testing.T) {
	if err := Decode(url.Values{}, nil); err == nil {
		t.Fatalf("expected error when dest is nil")
	}

	var notStruct int
	if err := Decode(url.Values{}, &notStruct); err == nil {
		t.Fatalf("expected error when dest is not struct pointer")
	}
}
