package bodies

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compression"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
)

const (
	LoginLenFieldSize              = 4
	MaxCompressionsPerOneHandshake = 256
)

var _ body.Body = Handshake{}

type Handshake struct {
	Login        string
	Hash         [32]byte
	Compressions []compression.Code
}

func NewHandshake(login string, hash [32]byte, compressions []compression.Code) Handshake {
	return Handshake{
		Login:        login,
		Hash:         hash,
		Compressions: filter(compressions),
	}
}

func filter(compressions []compression.Code) []compression.Code {
	filtered := make([]compression.Code, 0, len(compressions))
	used := make(map[compression.Code]struct{}, len(compressions))

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

	return LoginLenFieldSize + cap(h.Hash) + len(h.Login) + compression.FieldSize + size*compression.FieldSize
}

func (h Handshake) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(h.Login)))
	buf.AppendString(h.Login)
	buf.Append(h.Hash[:])
	buf.AppendUint8(uint8(min(len(h.Compressions), MaxCompressionsPerOneHandshake)))

	for i, compression := range h.Compressions {
		if i > MaxCompressionsPerOneHandshake {
			break
		}

		compression.Encode(buf)
	}
}

func (h Handshake) Command() command.Code {
	return command.Handshake
}

// Не проверяем Compression на валидность, т.к. по спеке могут быть кастомные алгоритмы.
func (h Handshake) IsValid() error {
	if len(h.Compressions) > MaxCompressionsPerOneHandshake {
		return errs.NewErrorMalformedValue(fmt.Sprintf("got %v compressions, but max compression per a handshake is %v", len(h.Compressions), MaxCompressionsPerOneHandshake))
	}

	return nil
}
