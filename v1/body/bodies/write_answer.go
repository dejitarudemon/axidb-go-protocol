package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Answer = WriteAnswer{}

// WriteAnswer is a successful reply to [fields.Write].
type WriteAnswer struct {
	simpleOK
}

// Command returns [fields.Answer].
func (w WriteAnswer) Command() fields.Command {
	return fields.Answer
}

// IsResponseTo returns [fields.Write].
func (w WriteAnswer) IsResponseTo() fields.Command {
	return fields.Write
}

// Encode writes the success flag and the replied-to command into buf.
func (w WriteAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	w.IsResponseTo().Encode(buf)
}
