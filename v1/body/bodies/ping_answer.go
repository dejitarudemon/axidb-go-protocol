package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Answer = PingAnswer{}

type PingAnswer struct {
	simpleOK
}

func (p PingAnswer) Command() fields.Command {
	return fields.Answer
}
func (p PingAnswer) IsResponseTo() fields.Command {
	return fields.Ping
}

func (p PingAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	p.IsResponseTo().Encode(buf)
}
