package body

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// Body is a protocol v0 Hello payload that can be validated and encoded.
type Body interface {
	// Size returns the encoded body size in bytes.
	Size() int

	// Encode writes the wire encoding of the body into buf.
	Encode(buf buffer.Appender)

	// IsValid reports whether the body satisfies protocol rules.
	IsValid() error
}
