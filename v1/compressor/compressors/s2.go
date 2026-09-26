package compressors

import (
	"bytes"
	"io"

	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/klauspost/compress/s2"
)

var _ compressor.Compressor = S2{}

// S2 compresses frame bodies with the Snappy-framed S2 algorithm ([fields.S2]).
type S2 struct {
	limit int64
}

// NewS2 returns an S2 compressor that rejects decompressed output larger than limit bytes.
func NewS2(limit int64) (S2, error) {
	return S2{limit: limit}, nil
}

// Code returns [fields.S2].
func (s S2) Code() fields.Compression {
	return fields.S2
}

// Compress returns the S2-compressed form of data.
func (s S2) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := s2.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decompress returns the original bytes from an S2 payload.
// Output longer than the configured limit is rejected.
func (s S2) Decompress(data []byte) ([]byte, error) {
	reader := s2.NewReader(bytes.NewReader(data))
	limited := io.LimitReader(reader, s.limit+1)

	out, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(out)) > s.limit {
		return nil, errs.NewErrorBodyLimitIsExceeded(uint32(len(out)), fields.BodyLimit(s.limit))
	}
	return out, nil
}
