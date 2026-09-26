package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Answer = PingAnswer{}

// PingAnswer is a successful reply to [fields.Ping].
type PingAnswer struct {
	simpleOK
}

// Command returns [fields.Answer].
func (p PingAnswer) Command() fields.Command {
	return fields.Answer
}

// IsResponseTo returns [fields.Ping].
func (p PingAnswer) IsResponseTo() fields.Command {
	return fields.Ping
}

// Encode writes the success flag and the replied-to command into buf.
func (p PingAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	p.IsResponseTo().Encode(buf)
}
