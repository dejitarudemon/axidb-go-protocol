package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = TypedArray{}

// TypedArrayLenFieldSize is the encoded size in bytes of the element count prefix for [TypedArray].
const TypedArrayLenFieldSize = 4

// TypedArray is a homogeneous sequence of [value.V] elements with type code [fields.TypedArray].
type TypedArray struct {
	// ElemType is the wire type shared by all non-nil elements.
	ElemType fields.Type
	// Elems holds the array elements in order.
	Elems []value.V
}

// realLen returns the number of non-nil elements.
func (ta TypedArray) realLen() int {
	realLen := 0
	for _, elem := range ta.Elems {
		if elem != nil {
			realLen += 1
		}
	}

	return realLen
}

// Encode writes the wire encoding of the value into buf.
// Nil elements are omitted from the encoded sequence.
func (ta TypedArray) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(ta.realLen()))
	ta.ElemType.Encode(buf)

	for _, elem := range ta.Elems {
		if elem != nil {
			elem.Encode(buf)
		}
	}
}

// Size returns the encoded value size in bytes.
func (ta TypedArray) Size() int {
	size := 0

	for _, elem := range ta.Elems {
		if elem != nil {
			size += elem.Size()
		}
	}

	return TypedArrayLenFieldSize + ta.ElemType.Size() + size
}

// Type returns [fields.TypedArray].
func (ta TypedArray) Type() fields.Type {
	return fields.TypedArray
}

// IsValid reports whether the value satisfies type-specific rules.
// Elements must be non-nil, match ElemType, and pass their own IsValid check.
func (ta TypedArray) IsValid() error {
	for i, elem := range ta.Elems {
		if elem == nil {
			return err.NewValidationError(
				"nil elem in TypedArray",
				"index", i,
			)
		}
		if elem.Type() != ta.ElemType {
			return err.NewValidationError(
				"wrong elem's type in TypedArray",
				"expected", ta.ElemType,
				"got", elem.Type(),
				"index", i,
			)
		}
		if err := elem.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
