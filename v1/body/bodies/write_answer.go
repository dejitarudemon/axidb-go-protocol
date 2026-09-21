package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
)

var _ body.Body = WriteAnswer{}

type WriteAnswer struct {
	simpleOK
}

func (p WriteAnswer) Command() command.Code {
	return command.Answer
}
