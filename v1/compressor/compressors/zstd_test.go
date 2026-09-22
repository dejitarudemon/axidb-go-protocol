package compressors

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var (
	zstdZero, _   = NewZstd(0)
	zstdNormal, _ = NewZstd(1 << 20)
)

func randomBytes(size int) []byte {
	r := rand.New(rand.NewSource(42))
	buf := make([]byte, 0, size)

	for range size {
		buf = append(buf, byte(r.Int63()))
	}

	return buf
}

func TestZstd_Code(t *testing.T) {
	tests := []struct {
		z    Zstd
		want fields.Compression
	}{
		{zstdZero, fields.Zstd},
		{zstdNormal, fields.Zstd},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestZstd_Code %v", i),
			func(t *testing.T) {
				if got := tt.z.Code(); got != tt.want {
					t.Errorf("got %v want %v", got, tt.want)
				}
			},
		)
	}

}

func TestZstd_EncodeDecode(t *testing.T) {
	tests := []struct {
		z    Zstd
		data []byte
	}{
		{zstdNormal, []byte{}},
		{zstdNormal, randomBytes(1 << 10)},
		{zstdNormal, randomBytes(1 << 11)},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestZstd_Code %v", i),
			func(t *testing.T) {
				compressed, err := tt.z.Compress(tt.data)

				if err != nil {
					t.Fatalf("got err: %v", err)
					return
				}

				decompress, err := tt.z.Decompress(compressed)

				if err != nil {
					t.Fatalf("got err: %v", err)
					return
				}

				if !bytes.Equal(tt.data, decompress) {
					t.Errorf("got %v want %v", decompress, tt.data)
				}
			},
		)
	}

}
