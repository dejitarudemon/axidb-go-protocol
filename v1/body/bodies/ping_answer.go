package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Body = PingAnswer{}

type PingAnswer struct {
	simpleOK
}

func (p PingAnswer) Command() fields.Command {
	return fields.Answer
}
