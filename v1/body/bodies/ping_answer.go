package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
)

var _ body.Body = PingAnswer{}

type PingAnswer struct {
	simpleOK
}

func (p PingAnswer) Command() command.Code {
	return command.Answer
}
