package compressors

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestS2_Code(t *testing.T) {
	tests := []struct {
		z    S2
		want fields.Compression
	}{
		{S2{limit: 1000}, fields.S2},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestS2_Code %v", i),
			func(t *testing.T) {
				if got := tt.z.Code(); got != tt.want {
					t.Errorf("got %v want %v", got, tt.want)
				}
			},
		)
	}

}

func TestS2_EncodeDecode(t *testing.T) {
	tests := []struct {
		z    S2
		data []byte
	}{
		{S2{limit: 1}, []byte{}},
		{S2{limit: 1<<20 + 1}, randomBytes(1 << 20)},
		{S2{limit: 1 << 12}, randomBytes(1 << 11)},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestS2_EncodeDecode %v", i),
			func(t *testing.T) {
				compressed, err := tt.z.Compress(tt.data)

				if err != nil {
					t.Fatalf("encode: got err: %v", err)
					return
				}

				decompress, err := tt.z.Decompress(compressed)

				if err != nil {
					t.Fatalf("decode: got err: %v", err)
					return
				}

				if !bytes.Equal(tt.data, decompress) {
					t.Errorf("got %v want %v", decompress, tt.data)
				}
			},
		)
	}

}
