package compressor

import "github.com/dejitarudemon/axidb-go-protocol/v1/fields"

// Compressor compresses and decompresses frame bodies for a protocol compression code.
type Compressor interface {
	// Code returns the protocol compression algorithm identifier.
	Code() fields.Compression

	// Compress returns the compressed form of data.
	Compress(buf []byte) ([]byte, error)

	// Decompress returns the original bytes from a compressed payload.
	Decompress(buf []byte) ([]byte, error)
}
