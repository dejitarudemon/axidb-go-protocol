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
	return fields.Key(d).Size()
}

func (d Delete) Encode(buf buffer.Appender) {
	fields.Key(d).Encode(buf)
}

func (d Delete) Command() fields.Command {
	return fields.Delete
}

func (d Delete) IsValid() error {
	if fields.Key(d).Size() == 0 {
		return err.NewValidationError(
			"empty key",
			"target", d.Command(),
		)
	}

	return nil
}
