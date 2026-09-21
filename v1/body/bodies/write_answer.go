package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Body = WriteAnswer{}

type WriteAnswer struct {
	simpleOK
}

func (p WriteAnswer) Command() fields.Command {
	return fields.Answer
}
