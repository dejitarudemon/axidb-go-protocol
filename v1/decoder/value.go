package decoder

import (
	"fmt"
	"math"

	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

// decodeType reads a value type code and rejects unknown codes.
func (d Decoder) decodeType(c *cursor) (fields.Type, error) {
	t, e := c.uint8()
	if e != nil {
		return 0, e
	}

	valueType := fields.Type(t)
	if !valueType.IsValid() {
		return 0, errs.NewErrorMalformedValue(fmt.Sprintf("unknown value type: %v", valueType))
	}

	return valueType, nil
}

// decodeTypedValue reads a type code followed by a value of that type.
func (d Decoder) decodeTypedValue(c *cursor) (value.V, error) {
	valueType, e := d.decodeType(c)
	if e != nil {
		return nil, e
	}

	return d.decodeValue(c, valueType)
}

// decodeValue decodes a value of valueType.
func (d Decoder) decodeValue(c *cursor, valueType fields.Type) (value.V, error) {
	switch valueType {
	case fields.Bytes:
		return d.decodeBytesValue(c)
	case fields.TypedArray:
		return d.decodeTypedArray(c)
	case fields.UntypedArray:
		return d.decodeUntypedArray(c)
	case fields.Int:
		return d.decodeIntValue(c)
	case fields.Uint:
		return d.decodeUintValue(c)
	case fields.Float:
		return d.decodeFloatValue(c)
	case fields.String:
		return d.decodeStringValue(c)
	case fields.JSON:
		return d.decodeJSONValue(c)
	}

	return nil, errs.NewErrorMalformedValue(fmt.Sprintf("unknown value type: %v", valueType))
}

// decodeLengthPrefixed reads a uint32 length followed by that many bytes.
func (d Decoder) decodeLengthPrefixed(c *cursor) ([]byte, error) {
	n, e := c.length()
	if e != nil {
		return nil, e
	}

	return c.bytes(n)
}

// decodeBytesValue decodes a bytes value.
func (d Decoder) decodeBytesValue(c *cursor) (values.Bytes, error) {
	b, e := d.decodeLengthPrefixed(c)
	if e != nil {
		return nil, e
	}

	return values.Bytes(b), nil
}

// decodeStringValue decodes a string value.
func (d Decoder) decodeStringValue(c *cursor) (values.String, error) {
	b, e := d.decodeLengthPrefixed(c)
	if e != nil {
		return "", e
	}

	return values.String(b), nil
}

// decodeJSONValue decodes a JSON value.
func (d Decoder) decodeJSONValue(c *cursor) (values.JSON, error) {
	b, e := d.decodeLengthPrefixed(c)
	if e != nil {
		return nil, e
	}

	return values.JSON(b), nil
}

// decodeIntValue decodes an int value.
func (d Decoder) decodeIntValue(c *cursor) (values.Int, error) {
	v, e := c.uint64()
	return values.Int(int64(v)), e
}

// decodeUintValue decodes a uint value.
func (d Decoder) decodeUintValue(c *cursor) (values.Uint, error) {
	v, e := c.uint64()
	return values.Uint(v), e
}

// decodeFloatValue decodes a float value.
func (d Decoder) decodeFloatValue(c *cursor) (values.Float, error) {
	v, e := c.uint64()
	return values.Float(math.Float64frombits(v)), e
}

// decodeUntypedArray decodes an untyped array where every element carries its own type.
func (d Decoder) decodeUntypedArray(c *cursor) (values.UntypedArray, error) {
	n, e := c.uint32()
	if e != nil {
		return nil, e
	}

	// Every element takes at least one byte, so remaining bounds a trustworthy capacity.
	elems := make(values.UntypedArray, 0, min(int(n), c.remaining()))

	for range n {
		elem, e := d.decodeTypedValue(c)
		if e != nil {
			return nil, e
		}

		elems = append(elems, elem)
	}

	return elems, nil
}

// decodeTypedArray decodes a typed array where all elements share one type.
func (d Decoder) decodeTypedArray(c *cursor) (values.TypedArray, error) {
	n, e := c.uint32()
	if e != nil {
		return values.TypedArray{}, e
	}

	elemType, e := d.decodeType(c)
	if e != nil {
		return values.TypedArray{}, e
	}

	elems := make([]value.V, 0, min(int(n), c.remaining()))

	for range n {
		elem, e := d.decodeValue(c, elemType)
		if e != nil {
			return values.TypedArray{}, e
		}

		elems = append(elems, elem)
	}

	return values.TypedArray{ElemType: elemType, Elems: elems}, nil
}
