package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
)

var _ body.Body = DeleteAnswer{}

type DeleteAnswer struct {
	simpleOK
}

func (p DeleteAnswer) Command() command.Code {
	return command.Answer
}
