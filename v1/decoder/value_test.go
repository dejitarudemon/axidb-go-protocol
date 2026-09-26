package decoder

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func encodeTypedValue(t testing.TB, v value.V) []byte {
	t.Helper()

	buf := buffer.Slice{}
	buf.Preallocate(fields.TypeFieldSize + v.Size())
	v.Type().Encode(&buf)
	v.Encode(&buf)
	return buf.Bytes()
}

func TestDecoder_decodeTypedValue(t *testing.T) {
	tests := []struct {
		name string
		want value.V
	}{
		{"empty bytes", values.Bytes{}},
		{"bytes", values.Bytes{0x00, 0xFF}},
		{"empty string", values.String("")},
		{"string", values.String("hello")},
		{"json", values.JSON(`{"a":1}`)},
		{"min int", values.Int(math.MinInt64)},
		{"negative int", values.Int(-1)},
		{"max uint", values.Uint(math.MaxUint64)},
		{"float", values.Float(-2.5)},
		{"empty typed array", values.TypedArray{ElemType: fields.Uint}},
		{"typed array", values.TypedArray{ElemType: fields.JSON, Elems: []value.V{values.JSON(`[]`), values.JSON(`{}`)}}},
		{"empty untyped array", values.UntypedArray{}},
		{"nested arrays", values.UntypedArray{
			values.Uint(7),
			values.JSON(`null`),
			values.Bytes{0x01},
			values.TypedArray{ElemType: fields.Float, Elems: []value.V{values.Float(1.5)}},
			values.UntypedArray{values.String("nested")},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCursor(encodeTypedValue(t, tt.want))
			got, e := NewDecoder(1024, nil).decodeTypedValue(c)
			assertDecoded(t, got, e, c, tt.want)
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
			_, e := NewDecoder(1024, nil).decodeTypedValue(newCursor(tt.data))
			assertMalformed(t, e)
		})
	}
}

func TestDecoder_decodeValue_UnknownType(t *testing.T) {
	_, e := NewDecoder(1024, nil).decodeValue(newCursor(u64(0)), fields.Type(0xFF))
	assertMalformed(t, e)
}

func TestDecoder_decodeScalarValues(t *testing.T) {
	d := NewDecoder(1024, nil)

	tests := []struct {
		name string
		data []byte
		read func(*cursor) (value.V, error)
		want value.V
	}{
		{"int", u64(math.MaxUint64), func(c *cursor) (value.V, error) { return d.decodeIntValue(c) }, values.Int(-1)},
		{"float", u64(math.Float64bits(3.25)), func(c *cursor) (value.V, error) { return d.decodeFloatValue(c) }, values.Float(3.25)},
		{"json", cat(u32(2), []byte("{}")), func(c *cursor) (value.V, error) { return d.decodeJSONValue(c) }, values.JSON("{}")},
		{"uint", u64(1), func(c *cursor) (value.V, error) { return d.decodeUintValue(c) }, values.Uint(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCursor(tt.data)
			got, e := tt.read(c)
			assertDecoded(t, got, e, c, tt.want)
		})
	}
}
