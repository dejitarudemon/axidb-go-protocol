package bodies

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compression"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
)

var _ body.Body = HandshakeAnswer{}

type HandshakeAnswer struct {
	Compressions []compression.Code
}

func NewHandshakeAnswer(compressions []compression.Code) HandshakeAnswer {
	return HandshakeAnswer{
		Compressions: filter(compressions),
	}
}

func (h HandshakeAnswer) Size() int {
	size := min(len(h.Compressions), MaxCompressionsPerOneHandshake)

	return ResultFieldSize + compression.FieldSize + size*compression.FieldSize
}

func (h HandshakeAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	buf.AppendUint8(uint8(min(len(h.Compressions), MaxCompressionsPerOneHandshake)))

	for i, compression := range h.Compressions {
		if i > MaxCompressionsPerOneHandshake {
			break
		}

		compression.Encode(buf)
	}
}

func (h HandshakeAnswer) Command() command.Code {
	return command.Answer
}

// Не проверяем Compression на валидность, т.к. по спеке могут быть кастомные алгоритмы.
func (h HandshakeAnswer) IsValid() error {
	if len(h.Compressions) > MaxCompressionsPerOneHandshake {
		return errs.NewErrorMalformedValue(fmt.Sprintf("got %v compressions, but max compression per a handshake is %v", len(h.Compressions), MaxCompressionsPerOneHandshake))
	}

	return nil
}
