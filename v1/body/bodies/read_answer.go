package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ body.Body = ReadAnswer{}

type ReadAnswer struct {
	Value value.V
}

func (r ReadAnswer) Size() int {
	if r.Value == nil {
		return ResultFieldSize
	}

	return ResultFieldSize + r.Value.Type().Size() + r.Value.Size()
}

func (r ReadAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)

	if r.Value != nil {
		r.Value.Type().Encode(buf)
		r.Value.Encode(buf)
	}

}

func (r ReadAnswer) Command() fields.Command {
	return fields.Answer
}

func (r ReadAnswer) IsValid() error {
	if r.Value == nil {
		return err.NewValidationError(
			"nil value",
			"target", r.Command(),
		)
	}

	return nil
}
