package decoder

import (
	"fmt"
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func encodeTypedValue(v value.V) []byte {
	buf := buffer.Slice{}
	buf.Preallocate(fields.TypeFieldSize + v.Size())
	v.Type().Encode(&buf)
	v.Encode(&buf)
	return buf.Bytes()
}

func TestDecoder_decodeTypedValue_RoundTrip(t *testing.T) {
	tests := []value.V{
		values.Bytes{},
		values.Bytes{0x00, 0xFF},
		values.String(""),
		values.String("hello"),
		values.JSON(`{"a":1}`),
		values.Int(math.MinInt64),
		values.Int(-1),
		values.Uint(math.MaxUint64),
		values.Float(-2.5),
		values.TypedArray{ElemType: fields.Uint, Elems: []value.V{}},
		values.TypedArray{ElemType: fields.JSON, Elems: []value.V{values.JSON(`[]`), values.JSON(`{}`)}},
		values.UntypedArray{},
		values.UntypedArray{
			values.Uint(7),
			values.JSON(`null`),
			values.Bytes{0x01},
			values.TypedArray{ElemType: fields.Float, Elems: []value.V{values.Float(1.5)}},
			values.UntypedArray{values.String("nested")},
		},
	}

	for i, want := range tests {
		t.Run(fmt.Sprintf("%v_%v", i, want.Type()), func(t *testing.T) {
			d := NewDecoder(1024, nil)
			c := newCursor(encodeTypedValue(want))

			got, e := d.decodeTypedValue(c)
			if e != nil {
				t.Fatalf("decodeTypedValue: got err %v", e)
			}

			if e := c.expectEnd(); e != nil {
				t.Fatalf("decodeTypedValue: %v", e)
			}

			compareValues(t, got, want)
		})
	}
}

func TestDecoder_decodeTypedValue_Errs(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"empty", nil},
		{"unknown type", []byte{byte(fields.JSON) + 1}},
		{"unknown type 255", []byte{0xFF}},
		{"bytes without len", []byte{byte(fields.Bytes), 0x00}},
		{"bytes shorter than len", cat([]byte{byte(fields.Bytes)}, u32(3), []byte{0x01})},
		{"string shorter than len", cat([]byte{byte(fields.String)}, u32(1))},
		{"json shorter than len", cat([]byte{byte(fields.JSON)}, u32(1))},
		{"int short", cat([]byte{byte(fields.Int)}, u32(1))},
		{"uint short", cat([]byte{byte(fields.Uint)}, u32(1))},
		{"float short", cat([]byte{byte(fields.Float)}, u32(1))},
		{"typed array without len", []byte{byte(fields.TypedArray)}},
		{"typed array without elem type", cat([]byte{byte(fields.TypedArray)}, u32(1))},
		{"typed array unknown elem type", cat([]byte{byte(fields.TypedArray)}, u32(0), []byte{0xFF})},
		{"typed array short elem", cat([]byte{byte(fields.TypedArray)}, u32(1), []byte{byte(fields.Int)}, u32(0))},
		{"typed array huge len", cat([]byte{byte(fields.TypedArray)}, u32(math.MaxUint32), []byte{byte(fields.Int)}, u64(1))},
		{"untyped array without len", []byte{byte(fields.UntypedArray)}},
		{"untyped array missing elem type", cat([]byte{byte(fields.UntypedArray)}, u32(1))},
		{"untyped array unknown elem type", cat([]byte{byte(fields.UntypedArray)}, u32(1), []byte{0xFF})},
		{"untyped array short elem", cat([]byte{byte(fields.UntypedArray)}, u32(1), []byte{byte(fields.Float)})},
		{"untyped array huge len", cat([]byte{byte(fields.UntypedArray)}, u32(math.MaxUint32), []byte{byte(fields.Int)}, u64(1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDecoder(1024, nil)

			_, e := d.decodeTypedValue(newCursor(tt.data))
			assertMalformed(t, e)
		})
	}
}

func TestDecoder_decodeValue_UnknownType(t *testing.T) {
	d := NewDecoder(1024, nil)

	_, e := d.decodeValue(newCursor(u64(0)), fields.Type(0xFF))
	assertMalformed(t, e)
}

func TestDecoder_decodeScalarValues(t *testing.T) {
	d := NewDecoder(1024, nil)

	i, e := d.decodeIntValue(newCursor(u64(math.MaxUint64)))
	if e != nil || i != values.Int(-1) {
		t.Errorf("decodeIntValue: got %v, %v", i, e)
	}

	f, e := d.decodeFloatValue(newCursor(u64(math.Float64bits(3.25))))
	if e != nil || f != values.Float(3.25) {
		t.Errorf("decodeFloatValue: got %v, %v", f, e)
	}

	j, e := d.decodeJSONValue(newCursor(cat(u32(2), []byte("{}"))))
	if e != nil || string(j) != "{}" {
		t.Errorf("decodeJSONValue: got %s, %v", j, e)
	}
}
