package err

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

// TracebackIDFieldSize is the number of bytes reserved for a traceback ID in an error message.
const TracebackIDFieldSize = 16

// ProtocolError is a protocol v1 error that can be encoded in a message.
type ProtocolError interface {
	error

	// Size returns the encoded error message size in bytes.
	Size() int

	// TracebackID returns the error traceback ID.
	TracebackID() fields.TracebackID

	// Encode writes the wire encoding of the error into buf.
	Encode(buf buffer.Appender)

	// Code returns the protocol error code.
	Code() fields.Error
}
