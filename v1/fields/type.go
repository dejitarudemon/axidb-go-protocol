package fields

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// Type is a protocol v1 value type code.
type Type uint8

// TypeFieldSize is the encoded size in bytes of a [Type].
const TypeFieldSize = 1

// Value type codes defined by protocol v1.
const (
	Bytes Type = iota
	TypedArray
	UntypedArray
	Int
	Uint
	Float
	String
	JSON
)

// Encode writes the wire encoding of the type code into buf.
func (t Type) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(t))
}

// String returns a human-readable type name.
func (t Type) String() string {
	switch t {
	case Bytes:
		return "Bytes"
	case TypedArray:
		return "Typed Array"
	case UntypedArray:
		return "Untyped Array"
	case Int:
		return "Int"
	case Uint:
		return "Uint"
	case Float:
		return "Float"
	case String:
		return "String"
	case JSON:
		return "JSON"
	}

	return fmt.Sprintf("Unknown (%d)", t)
}

// IsValid reports whether t is a known value type code (0 through [JSON]).
func (t Type) IsValid() bool {
	return t <= JSON
}

// Size returns the encoded size in bytes.
func (t Type) Size() int {
	return TypeFieldSize
}
