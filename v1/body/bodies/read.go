package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Body = Read{}

type Read fields.Key

func (r Read) Size() int {
	return len(r)
}

func (r Read) Encode(buf buffer.Appender) {
	buf.Append(r)
}

func (r Read) Command() fields.Command {
	return fields.Read
}

func (r Read) IsValid() error {
	if len(r) == 0 {
		return err.NewValidationError(
			"empty key",
			"target", r.Command(),
		)
	}
	return nil
}
