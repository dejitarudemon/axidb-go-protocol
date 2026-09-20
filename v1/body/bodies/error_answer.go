package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
)

var ResultNotOK = []byte{0x00}

var _ body.Body = ErrorAnswer{}

type ErrorAnswer struct {
	Err err.Error
}

func (e ErrorAnswer) Size() int {
	return ResultFieldSize + e.Size()
}

func (e ErrorAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultNotOK)
	e.Err.Encode(buf)
}

func (e ErrorAnswer) Command() command.Code {
	return command.Answer
}

func (e ErrorAnswer) IsValid() error {
	if !e.Err.IsValid() {
		return errs.NewErrorMalformedValue("invalid error")
	}

	return nil
}
