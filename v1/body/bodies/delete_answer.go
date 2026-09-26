package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Answer = DeleteAnswer{}

// DeleteAnswer is a successful reply to [fields.Delete].
type DeleteAnswer struct {
	simpleOK
}

// Command returns [fields.Answer].
func (d DeleteAnswer) Command() fields.Command {
	return fields.Answer
}

// IsResponseTo returns [fields.Delete].
func (d DeleteAnswer) IsResponseTo() fields.Command {
	return fields.Delete
}

// Encode writes the success flag and the replied-to command into buf.
func (d DeleteAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	d.IsResponseTo().Encode(buf)
}
