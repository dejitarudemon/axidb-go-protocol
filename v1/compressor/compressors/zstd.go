package compressors

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/klauspost/compress/zstd"
)

var _ compressor.Compressor = Zstd{}

type Zstd struct {
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func NewZstd(limit uint32) (Zstd, error) {
	encoder, err := zstd.NewWriter(nil,
		zstd.WithEncoderLevel(zstd.SpeedDefault),
	)
	if err != nil {
		return Zstd{}, err
	}

	decoder, err := zstd.NewReader(nil,
		zstd.WithDecoderMaxMemory(uint64(limit)),
		zstd.WithDecoderConcurrency(1),
	)
	if err != nil {
		return Zstd{}, err
	}

	return Zstd{
		encoder: encoder,
		decoder: decoder,
	}, nil
}

func (z Zstd) Close() {
	z.decoder.Close()
	z.encoder.Close()
}

func (z Zstd) Code() fields.Compression {
	return fields.Zstd
}

func (z Zstd) Compress(data []byte) ([]byte, error) {
	compressed := make([]byte, 0, len(data))
	return z.encoder.EncodeAll(data, compressed), nil
}

func (z Zstd) Decompress(data []byte) ([]byte, error) {
	decompressed := make([]byte, 0, len(data))
	return z.decoder.DecodeAll(data, decompressed)
}
