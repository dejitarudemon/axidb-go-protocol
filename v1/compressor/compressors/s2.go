package compressors

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/klauspost/compress/s2"
)

var _ compressor.Compressor = S2{}

type S2 struct{}

func (s S2) Code() fields.Compression {
	return fields.S2
}

func (s S2) Compress(data []byte) ([]byte, error) {
	compressed := make([]byte, 0, len(data))
	return s2.Encode(compressed, data), nil
}

func (s S2) Decompress(data []byte) ([]byte, error) {
	decompressed := make([]byte, 0, len(data))
	return s2.Decode(data, decompressed)
}
