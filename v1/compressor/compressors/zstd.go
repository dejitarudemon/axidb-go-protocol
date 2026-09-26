package compressors

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/klauspost/compress/zstd"
)

var _ compressor.Compressor = Zstd{}

// Zstd compresses frame bodies with Zstandard ([fields.Zstd]).
type Zstd struct {
	limit uint32
	enc   *zstd.Encoder
	dec   *zstd.Decoder
}

// NewZstd returns a Zstd compressor that rejects decompressed output larger than limit bytes.
func NewZstd(limit uint32) (Zstd, error) {
	enc, err := zstd.NewWriter(nil,
		zstd.WithEncoderLevel(zstd.SpeedDefault),
		zstd.WithEncoderConcurrency(1),
	)
	if err != nil {
		return Zstd{}, err
	}

	window := uint64(limit)
	if window < uint64(zstd.MinWindowSize) {
		window = uint64(zstd.MinWindowSize)
	}

	dec, err := zstd.NewReader(nil,
		zstd.WithDecoderMaxMemory(window),
		zstd.WithDecoderMaxWindow(window),
		zstd.WithDecoderConcurrency(1),
	)
	if err != nil {
		_ = enc.Close()
		return Zstd{}, err
	}

	return Zstd{
		limit: limit,
		enc:   enc,
		dec:   dec,
	}, nil
}

// Code returns [fields.Zstd].
func (z Zstd) Code() fields.Compression {
	return fields.Zstd
}

// Compress returns the Zstandard-compressed form of data.
func (z Zstd) Compress(data []byte) ([]byte, error) {
	return z.enc.EncodeAll(data, nil), nil
}

// Decompress returns the original bytes from a Zstandard payload.
// Output longer than the configured limit is rejected.
func (z Zstd) Decompress(data []byte) ([]byte, error) {
	if z.limit < uint32(zstd.MinWindowSize) {
		return nil, zstd.ErrWindowSizeTooSmall
	}

	out, err := z.dec.DecodeAll(data, nil)
	if err != nil {
		return nil, err
	}
	if uint32(len(out)) > z.limit {
		return nil, zstd.ErrCompressedSizeTooBig
	}
	return out, nil
}
