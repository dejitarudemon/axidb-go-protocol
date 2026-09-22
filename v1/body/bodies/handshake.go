package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

const (
	LoginLenFieldSize              = 4
	MaxCompressionsPerOneHandshake = 255
)

var _ body.Body = Handshake{}

type Handshake struct {
	Login        string
	Hash         [32]byte
	Compressions []fields.Compression
}

func NewHandshake(login string, hash [32]byte, compressions []fields.Compression) Handshake {
	return Handshake{
		Login:        login,
		Hash:         hash,
		Compressions: filter(compressions),
	}
}

func filter(compressions []fields.Compression) []fields.Compression {
	filtered := make([]fields.Compression, 0, len(compressions))
	used := make(map[fields.Compression]struct{}, len(compressions))

	for _, compression := range compressions {
		if _, ok := used[compression]; !ok {
			used[compression] = struct{}{}
			filtered = append(filtered, compression)
		}
	}

	return filtered
}

func (h Handshake) Size() int {
	size := min(len(h.Compressions), MaxCompressionsPerOneHandshake)

	return LoginLenFieldSize + cap(h.Hash) + len(h.Login) + fields.CompressionFieldSize + size*fields.CompressionFieldSize
}

func (h Handshake) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(h.Login)))
	buf.AppendString(h.Login)
	buf.Append(h.Hash[:])
	buf.AppendUint8(uint8(min(len(h.Compressions), MaxCompressionsPerOneHandshake)))

	for i, fields := range h.Compressions {
		if i > MaxCompressionsPerOneHandshake {
			break
		}

		fields.Encode(buf)
	}
}

func (h Handshake) Command() fields.Command {
	return fields.Handshake
}

// Не проверяем Compression на валидность, т.к. по спеке могут быть кастомные алгоритмы.
func (h Handshake) IsValid() error {
	if len(h.Compressions) > MaxCompressionsPerOneHandshake {
		return err.NewValidationError(
			"too many fieldss",
			"compressions", len(h.Compressions),
			"target", h.Command(),
		)
	}

	return nil
}
