package compressor

import "github.com/dejitarudemon/axidb-go-protocol/v1/fields"

type Compressor interface {
	Code() fields.Compression
	Compress(buf []byte) ([]byte, error)
	Decompress(buf []byte) ([]byte, error)
}
