package fields

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestKey_Size(k *testing.T) {
	tests := []struct {
		k    Key
		want int
	}{
		{Key([]byte{}), 0},
		{Key([]byte{0xff, 0xff, 0xff, 0xff}), 4},
	}

	for _, tt := range tests {
		k.Run(
			fmt.Sprintf("%v", tt.k),
			func(k *testing.T) {
				if got := tt.k.Size(); got != tt.want {
					k.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestKey_Encode(k *testing.T) {
	tests := []struct {
		k    Key
		want []byte
	}{
		{Key([]byte{}), []byte{}},
		{Key([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	}

	for _, tt := range tests {
		k.Run(
			fmt.Sprintf("%v", tt.k),
			func(k *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.k.Size())

				tt.k.Encode(&buf)
				got := buf.Bytes()

				if len(got) != tt.k.Size() {
					k.Fatalf("expected %v bytes, got %v bytes", tt.k.Size(), len(got))
				}

				if !bytes.Equal(got, tt.want) {
					k.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			},
		)
	}
}
