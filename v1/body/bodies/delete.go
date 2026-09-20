package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
)

var _ body.Body = Delete{}

type Delete struct {
	Key []byte
}

func (d Delete) Size() int {
	return len(d.Key)
}

func (d Delete) Encode(buf buffer.Appender) {
	buf.Append(d.Key)
}

func (d Delete) Command() command.Code {
	return command.Delete
}

func (d Delete) IsValid() bool {
	return len(d.Key) > 0
}
