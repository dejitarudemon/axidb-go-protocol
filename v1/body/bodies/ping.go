package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Body = Ping{}

type Ping struct{}

func (p Ping) Size() int {
	return 0
}

func (p Ping) Encode(buf buffer.Appender) {}

func (p Ping) Command() fields.Command {
	return fields.Ping
}

func (p Ping) IsValid() error {
	return nil
}
