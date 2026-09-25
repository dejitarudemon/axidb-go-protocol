package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var ResultNotOK = []byte{0x00}

var _ body.Answer = ErrorAnswer{}

type ErrorAnswer struct {
	Err err.ProtocolError
}

func (e ErrorAnswer) Size() int {
	if e.Err == nil {
		return ResultFieldSize
	}
	return ResultFieldSize + e.Err.Size()
}

func (e ErrorAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultNotOK)
	if e.Err != nil {
		e.Err.Encode(buf)
	}
}

func (e ErrorAnswer) IsResponseTo() fields.Command {
	return fields.Answer
}

func (e ErrorAnswer) Command() fields.Command {
	return fields.Answer
}

func (e ErrorAnswer) IsValid() error {
	if e.Err == nil {
		return err.NewValidationError(
			"expected ProtocolError, but got nil",
			"target", e.Command(),
		)
	}
	return nil
}
