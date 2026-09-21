package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Body = Delete{}

type Delete fields.Key

func (d Delete) Size() int {
	return len(d)
}

func (d Delete) Encode(buf buffer.Appender) {
	buf.Append(d)
}

func (d Delete) Command() fields.Command {
	return fields.Delete
}

func (d Delete) IsValid() error {
	if len(d) == 0 {
		return err.NewValidationError(
			"empty key",
			"target", d.Command(),
		)
	}

	return nil
}
