package compressors

import (
	"bytes"
	"io"

	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/klauspost/compress/zstd"
)

var _ compressor.Compressor = Zstd{}

// Zstd compresses frame bodies with Zstandard ([fields.Zstd]).
type Zstd struct {
	limit uint32
}

// NewZstd returns a Zstd compressor that rejects decompressed output larger than limit bytes.
func NewZstd(limit uint32) (Zstd, error) {
	return Zstd{
		limit: limit,
	}, nil
}

// Code returns [fields.Zstd].
func (z Zstd) Code() fields.Compression {
	return fields.Zstd
}

// Compress returns the Zstandard-compressed form of data.
func (z Zstd) Compress(data []byte) ([]byte, error) {
	writer, err := zstd.NewWriter(
		nil,
		zstd.WithEncoderLevel(zstd.SpeedDefault),
	)

	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return writer.EncodeAll(data, nil), nil
}

// Decompress returns the original bytes from a Zstandard payload.
// Output longer than the configured limit is rejected.
func (z Zstd) Decompress(data []byte) ([]byte, error) {
	reader, err := zstd.NewReader(
		bytes.NewReader(data),
		zstd.WithDecoderMaxMemory(uint64(z.limit)),
		zstd.WithDecoderMaxWindow(uint64(z.limit)),
		zstd.WithDecoderConcurrency(1),
	)

	if err != nil {
		return nil, err
	}
	defer reader.Close()

	limited := io.LimitReader(reader, int64(z.limit)+1)
	result, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}

	if len(result) > int(z.limit) {
		return nil, zstd.ErrCompressedSizeTooBig
	}

	return result, nil
}
