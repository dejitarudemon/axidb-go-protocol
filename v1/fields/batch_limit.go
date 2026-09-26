package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// BatchLimit is a maximum allowed number of requests in a batch.
type BatchLimit uint32

// BatchLimitFieldSize is the encoded size in bytes of a [BatchLimit].
const BatchLimitFieldSize = 4

// Encode writes the wire encoding of the batch limit into buf.
func (b BatchLimit) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(b))
}

// Size returns the encoded size in bytes.
func (b BatchLimit) Size() int {
	return BatchLimitFieldSize
}
