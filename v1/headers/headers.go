package headers

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

const (
	BodyLenFieldSize = 4
)

type Headers struct {
	Command     fields.Command
	RequestID   fields.RequestID
	Compression fields.Compression
	BodyLen     uint32
}

func (h Headers) Size() int {
	return fields.CommandFieldSize + h.RequestID.Size() + fields.CompressionFieldSize + BodyLenFieldSize
}

func (h Headers) IsValid() error {
	return nil
}

func (h Headers) Encode(buf buffer.Appender) {
	h.Command.Encode(buf)
	h.RequestID.Encode(buf)
	h.Compression.Encode(buf)
	buf.AppendUint32(uint32(h.BodyLen))
}
