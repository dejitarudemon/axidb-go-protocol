package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Answer = DeleteAnswer{}

type DeleteAnswer struct {
	simpleOK
}

func (d DeleteAnswer) Command() fields.Command {
	return fields.Answer
}

func (d DeleteAnswer) IsResponseTo() fields.Command {
	return fields.Delete
}

func (d DeleteAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	d.IsResponseTo().Encode(buf)
}
