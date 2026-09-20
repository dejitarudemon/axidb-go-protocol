package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestSimpleOK_Size(t *testing.T) {
	tests := []struct {
		s    simpleOK
		want int
	}{
		{simpleOK{}, 1},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSimpleOK_Size %v", tt.s),
			func(t *testing.T) {
				if got := tt.s.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSimpleOK_Encode(t *testing.T) {
	tests := []struct {
		s    simpleOK
		want []byte
	}{
		{simpleOK{}, []byte{0x01}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSimpleOK_Encode %v", tt.s),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.s.Size())

				tt.s.Encode(&buf)

				got := buf.Bytes()

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestSimpleOK_IsValid(t *testing.T) {
	tests := []struct {
		s    simpleOK
		want bool
	}{
		{simpleOK{}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSimpleOK_IsValid %v", tt.s),
			func(t *testing.T) {
				if got := tt.s.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
