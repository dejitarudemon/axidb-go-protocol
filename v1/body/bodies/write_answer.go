package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Answer = WriteAnswer{}

type WriteAnswer struct {
	simpleOK
}

func (w WriteAnswer) Command() fields.Command {
	return fields.Answer
}

func (w WriteAnswer) IsResponseTo() fields.Command {
	return fields.Write
}

func (w WriteAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	w.IsResponseTo().Encode(buf)
}
