package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
)

var _ body.Body = Read{}

type Read struct {
	Key []byte
}

func (r Read) Size() int {
	return len(r.Key)
}

func (r Read) Encode(buf buffer.Appender) {
	buf.Append(r.Key)
}

func (r Read) Command() command.Code {
	return command.Read
}

func (r Read) IsValid() bool {
	return len(r.Key) > 0
}
