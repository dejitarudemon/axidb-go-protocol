package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// RequestNumber is a zero-based index of a request inside a batch.
type RequestNumber uint32

// RequestNumberFieldSize is the encoded size in bytes of a [RequestNumber].
const RequestNumberFieldSize = 4

// Encode writes the wire encoding of the request number into buf.
func (r RequestNumber) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(r))
}

// Size returns the encoded size in bytes.
func (r RequestNumber) Size() int {
	return RequestNumberFieldSize
}
