package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = UntypedArray{}

// UntypedArrayLenFieldSize is the encoded size in bytes of the element count prefix for [UntypedArray].
const UntypedArrayLenFieldSize = 4

// UntypedArray is a heterogeneous sequence of [value.V] elements with type code [fields.UntypedArray].
type UntypedArray []value.V

// realLen returns the number of non-nil elements.
func (ua UntypedArray) realLen() int {
	realLen := 0
	for _, elem := range ua {
		if elem != nil {
			realLen += 1
		}
	}

	return realLen
}

// Encode writes the wire encoding of the value into buf.
// Nil elements are omitted from the encoded sequence.
func (ua UntypedArray) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(ua.realLen()))

	for _, elem := range ua {
		if elem != nil {
			elem.Type().Encode(buf)
			elem.Encode(buf)
		}
	}
}

// Size returns the encoded value size in bytes.
func (ua UntypedArray) Size() int {
	size := UntypedArrayLenFieldSize

	for _, elem := range ua {
		if elem != nil {
			size += elem.Size() + elem.Type().Size()
		}
	}

	return size
}

// Type returns [fields.UntypedArray].
func (ua UntypedArray) Type() fields.Type {
	return fields.UntypedArray
}

// IsValid reports whether the value satisfies type-specific rules.
// Elements must be non-nil and each element must pass its own IsValid check.
func (ua UntypedArray) IsValid() error {
	for i, elem := range ua {
		if elem == nil {
			return err.NewValidationError(
				"nil elem in UntypedArray",
				"index", i,
			)
		}

		if err := elem.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
