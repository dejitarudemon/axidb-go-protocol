package headers

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

const (
	RequestIDFieldSize = 4
	BodyLenFieldSize   = 4
)

type Headers struct {
	Command     fields.Command
	RequestID   uint32
	Compression fields.Compression
	BodyLen     uint32
}

func (h Headers) Size() int {
	return fields.CommandFieldSize + RequestIDFieldSize + fields.CompressionFieldSize + BodyLenFieldSize
}

func (h Headers) IsValid() error {
	return nil
}

func (h Headers) Encode(buf buffer.Appender) {
	h.Command.Encode(buf)
	buf.AppendUint32(uint32(h.RequestID))
	h.Compression.Encode(buf)
	buf.AppendUint32(uint32(h.BodyLen))
}
