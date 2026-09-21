package headers

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compression"
)

const (
	RequestIDFieldSize = 4
	BodyLenFieldSize   = 4
)

type Headers struct {
	Command     command.Code
	RequestID   uint32
	Compression compression.Code
	BodyLen     uint32
}

func (h Headers) Size() int {
	return command.FieldSize + RequestIDFieldSize + compression.FieldSize + BodyLenFieldSize
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
