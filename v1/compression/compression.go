package compression

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

type Compression uint8

const (
	None Compression = iota
	Zstd
	Lz4
)

const CompressionFieldSize = 1

func (c Compression) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

func (c Compression) String() string {
	switch c {
	case None:
		return "None"
	case Zstd:
		return "Zstd"
	case Lz4:
		return "Lz4"
	}

	return fmt.Sprintf("Unknown (%d)", c)
}
