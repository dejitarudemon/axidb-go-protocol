package compressors

import (
	"bytes"
	"errors"
	"math/rand"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/klauspost/compress/zstd"
)

func randomBytes(size int) []byte {
	r := rand.New(rand.NewSource(42))
	buf := make([]byte, size)
	r.Read(buf)
	return buf
}

var compressorsUnderTest = []struct {
	name     string
	new      func(limit uint32) (compressor.Compressor, error)
	wantCode fields.Compression
}{
	{
		"s2",
		func(limit uint32) (compressor.Compressor, error) { return NewS2(int64(limit)) },
		fields.S2,
	},
	{
		"zstd",
		func(limit uint32) (compressor.Compressor, error) { return NewZstd(limit) },
		fields.Zstd,
	},
}

func newCompressor(t *testing.T, newFn func(uint32) (compressor.Compressor, error), limit uint32) compressor.Compressor {
	t.Helper()

	c, err := newFn(limit)
	if err != nil {
		t.Fatalf("new(%v): got err: %v", limit, err)
	}
	return c
}

func TestCompressors_Code(t *testing.T) {
	for _, cc := range compressorsUnderTest {
		t.Run(cc.name, func(t *testing.T) {
			c := newCompressor(t, cc.new, 1<<20)

			if got := c.Code(); got != cc.wantCode {
				t.Errorf("Code() = %v, want %v", got, cc.wantCode)
			}
		})
	}
}

func TestCompressors_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		limit uint32
		data  []byte
	}{
		{"nil", 1 << 10, nil},
		{"empty", 1 << 10, []byte{}},
		{"one byte", 1 << 10, []byte{0x01}},
		{"compressible", 1 << 12, bytes.Repeat([]byte("abcd"), 1<<10)},
		{"random 2 KiB", 1 << 12, randomBytes(1 << 11)},
		{"exactly at limit", 1 << 11, randomBytes(1 << 11)},
		{"random 1 MiB", 1<<20 + 1, randomBytes(1 << 20)},
	}

	for _, cc := range compressorsUnderTest {
		t.Run(cc.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					c := newCompressor(t, cc.new, tt.limit)

					compressed, err := c.Compress(tt.data)
					if err != nil {
						t.Fatalf("Compress(): got err: %v", err)
					}

					got, err := c.Decompress(compressed)
					if err != nil {
						t.Fatalf("Decompress(): got err: %v", err)
					}

					if !bytes.Equal(got, tt.data) {
						t.Errorf("Decompress() = %d bytes, want %d bytes equal to input", len(got), len(tt.data))
					}
				})
			}
		})
	}
}

func TestCompressors_ZipBomb(t *testing.T) {
	tests := []struct {
		name  string
		limit uint32
		data  []byte
	}{
		{"one byte over limit", 1 << 11, randomBytes(1<<11 + 1)},
		{"double the limit", 1 << 12, randomBytes(1 << 13)},
		{"compressible", 1 << 12, bytes.Repeat([]byte{0x00}, 1<<16)},
		{"1 MiB over limit", 1<<20 - 1, randomBytes(1 << 20)},
	}

	for _, cc := range compressorsUnderTest {
		t.Run(cc.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					large := newCompressor(t, cc.new, uint32(len(tt.data)))
					compressed, err := large.Compress(tt.data)
					if err != nil {
						t.Fatalf("Compress(): got err: %v", err)
					}

					c := newCompressor(t, cc.new, tt.limit)
					if got, err := c.Decompress(compressed); err == nil {
						t.Errorf("Decompress() = %d bytes, want err", len(got))
					}
				})
			}
		})
	}
}

func TestCompressors_DecompressGarbage(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"ascii", []byte("garbage")},
		{"random", randomBytes(1 << 10)},
		{"all ones", bytes.Repeat([]byte{0xFF}, 64)},
	}

	for _, cc := range compressorsUnderTest {
		t.Run(cc.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					c := newCompressor(t, cc.new, 1<<20)

					if got, err := c.Decompress(tt.data); err == nil {
						t.Errorf("Decompress() = %d bytes, want err", len(got))
					}
				})
			}
		})
	}
}

func TestCompressors_DecompressTruncated(t *testing.T) {
	for _, cc := range compressorsUnderTest {
		t.Run(cc.name, func(t *testing.T) {
			c := newCompressor(t, cc.new, 1<<20)

			compressed, err := c.Compress(randomBytes(1 << 12))
			if err != nil {
				t.Fatalf("Compress(): got err: %v", err)
			}

			if got, err := c.Decompress(compressed[:len(compressed)/2]); err == nil {
				t.Errorf("Decompress() = %d bytes, want err", len(got))
			}
		})
	}
}

func TestS2_DecompressLimitError(t *testing.T) {
	tests := []struct {
		name  string
		limit int64
		data  []byte
	}{
		{"zero limit", 0, []byte{0x01}},
		{"one byte over limit", 1 << 10, randomBytes(1<<10 + 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewS2(tt.limit)
			if err != nil {
				t.Fatalf("NewS2(%v): got err: %v", tt.limit, err)
			}

			compressed, err := s.Compress(tt.data)
			if err != nil {
				t.Fatalf("Compress(): got err: %v", err)
			}

			_, err = s.Decompress(compressed)
			if !errors.As(err, new(errs.ErrorBodyLimitIsExceeded)) {
				t.Errorf("Decompress(): got err %v, want %T", err, errs.ErrorBodyLimitIsExceeded{})
			}
		})
	}
}

func TestZstd_DecompressLimitError(t *testing.T) {
	frame, err := zstd.NewWriter(nil)
	if err != nil {
		t.Fatalf("zstd.NewWriter(): got err: %v", err)
	}

	oneFrame := frame.EncodeAll(randomBytes(3<<10), nil)

	tests := []struct {
		name    string
		limit   uint32
		data    []byte
		wantErr error
	}{
		{"single frame over limit", 1 << 11, oneFrame, zstd.ErrDecoderSizeExceeded},
		{"concatenated frames over limit", 1 << 12, append(bytes.Clone(oneFrame), oneFrame...), zstd.ErrCompressedSizeTooBig},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z, err := NewZstd(tt.limit)
			if err != nil {
				t.Fatalf("NewZstd(%v): got err: %v", tt.limit, err)
			}

			if _, err := z.Decompress(tt.data); !errors.Is(err, tt.wantErr) {
				t.Errorf("Decompress(): got err %v, want %v", err, tt.wantErr)
			}
		})
	}

	if err := frame.Close(); err != nil {
		t.Fatalf("zstd.Close(): got err: %v", err)
	}
}

func TestZstd_DecompressLimitBelowMinWindow(t *testing.T) {
	tests := []struct {
		name  string
		limit uint32
	}{
		{"zero", 0},
		{"one byte", 1},
		{"one below min window", zstd.MinWindowSize - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z, err := NewZstd(tt.limit)
			if err != nil {
				t.Fatalf("NewZstd(%v): got err: %v", tt.limit, err)
			}

			compressed, err := z.Compress(nil)
			if err != nil {
				t.Fatalf("Compress(): got err: %v", err)
			}

			if got, err := z.Decompress(compressed); err == nil {
				t.Errorf("Decompress() = %d bytes, want err", len(got))
			}
		})
	}
}
