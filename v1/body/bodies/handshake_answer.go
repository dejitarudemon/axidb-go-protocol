package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Answer = HandshakeAnswer{}

type HandshakeAnswer struct {
	Compressions []fields.Compression
}

func NewHandshakeAnswer(compressions []fields.Compression) HandshakeAnswer {
	return HandshakeAnswer{
		Compressions: filter(compressions),
	}
}

func (h HandshakeAnswer) Size() int {
	size := min(len(h.Compressions), MaxCompressionsPerOneHandshake)

	return ResultFieldSize + fields.CompressionFieldSize + size*fields.CompressionFieldSize + h.IsResponseTo().Size()
}

func (h HandshakeAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	h.IsResponseTo().Encode(buf)
	buf.AppendUint8(uint8(min(len(h.Compressions), MaxCompressionsPerOneHandshake)))

	for i, compression := range h.Compressions {
		if i > MaxCompressionsPerOneHandshake {
			break
		}

		compression.Encode(buf)
	}
}

func (h HandshakeAnswer) Command() fields.Command {
	return fields.Answer
}

func (h HandshakeAnswer) IsResponseTo() fields.Command {
	return fields.Handshake
}

// Не проверяем Compression на валидность, т.к. по спеке могут быть кастомные алгоритмы.
func (h HandshakeAnswer) IsValid() error {
	if len(h.Compressions) > MaxCompressionsPerOneHandshake {
		return err.NewValidationError(
			"too many compressions",
			"compressions", len(h.Compressions),
			"target", h.Command(),
		)
	}

	return nil
}
