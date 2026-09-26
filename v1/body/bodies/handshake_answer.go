package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Answer = HandshakeAnswer{}

// HandshakeAnswer is a successful reply to [fields.Handshake] listing accepted compressions.
type HandshakeAnswer struct {
	// Compressions lists compression algorithms the server accepts.
	Compressions []fields.Compression
}

// NewHandshakeAnswer returns a HandshakeAnswer after dropping [fields.None] and duplicate compression codes.
func NewHandshakeAnswer(compressions []fields.Compression) HandshakeAnswer {
	return HandshakeAnswer{
		Compressions: filter(compressions),
	}
}

// Size returns the encoded body size in bytes.
func (h HandshakeAnswer) Size() int {
	size := min(len(h.Compressions), MaxCompressionsPerOneHandshake)

	return ResultFieldSize + fields.CompressionFieldSize + size*fields.CompressionFieldSize + h.IsResponseTo().Size()
}

// Encode writes the wire encoding of the body into buf.
func (h HandshakeAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	h.IsResponseTo().Encode(buf)
	buf.AppendUint8(uint8(min(len(h.Compressions), MaxCompressionsPerOneHandshake)))

	for i, compression := range h.Compressions {
		if i >= MaxCompressionsPerOneHandshake {
			break
		}

		compression.Encode(buf)
	}
}

// Command returns [fields.Answer].
func (h HandshakeAnswer) Command() fields.Command {
	return fields.Answer
}

// IsResponseTo returns [fields.Handshake].
func (h HandshakeAnswer) IsResponseTo() fields.Command {
	return fields.Handshake
}

// IsValid reports whether the body satisfies protocol rules.
// Compression codes are not checked because custom algorithms are allowed.
// The compression list must not exceed [MaxCompressionsPerOneHandshake].
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
