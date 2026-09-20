package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
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

func (r Read) IsValid() error {
	if len(r.Key) == 0 {
		return errs.NewErrorMalformedValue("expected key, but got nothing")
	}
	return nil
}
