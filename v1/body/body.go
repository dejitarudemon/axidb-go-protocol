package body

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

// BodyLenFieldSize is the encoded size in bytes of the body length prefix in a frame header.
const BodyLenFieldSize = 4

// Body is a protocol v1 message body that can be validated and encoded.
type Body interface {
	// Size returns the encoded body size in bytes.
	Size() int

	// Encode writes the wire encoding of the body into buf.
	Encode(buf buffer.Appender)

	// Command returns the protocol command code for the body.
	Command() fields.Command

	// IsValid reports whether the body satisfies protocol rules.
	IsValid() error
}

// Answer is a [Body] that replies to a specific command.
type Answer interface {
	Body

	// IsResponseTo returns the command this answer replies to.
	IsResponseTo() fields.Command
}
