package values

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
)

func TestUntypedArray_Size(t *testing.T) {
	tests := []struct {
		c    UntypedArray
		want int
	}{
		{
			UntypedArray{nil}, 4,
		},

		{
			UntypedArray{Int(1)}, 13,
		},
		{
			UntypedArray{Int(1), Int(2)}, 22,
		},
		{
			UntypedArray{Int(1), Float(1)}, 22,
		},
		{
			UntypedArray{Int(1), Bytes([]byte{0x01, 0x02, 0x03})}, 21,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.Size(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestUntypedArray_Encode(t *testing.T) {
	tests := []struct {
		c    UntypedArray
		want []byte
	}{
		{
			UntypedArray{nil},
			[]byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			UntypedArray{},
			[]byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			UntypedArray{Int(1)},
			[]byte{0x00, 0x00, 0x00, 0x01, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			UntypedArray{Int(1), Int(2)},
			[]byte{0x00, 0x00, 0x00, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
		},
		{
			UntypedArray{Int(1), Float(1)},
			[]byte{0x00, 0x00, 0x00, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x05, 0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			UntypedArray{Int(1), Bytes([]byte{0x01, 0x02, 0x03})},
			[]byte{0x00, 0x00, 0x00, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x03},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.c.Size())

				tt.c.Encode(&buf)
				got := buf.Bytes()

				if len(got) != tt.c.Size() {
					t.Fatalf("expected %v bytes, got %v bytes", tt.c.Size(), len(got))
				}

				if !bytes.Equal(got, tt.want) {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestUntypedArray_Type(t *testing.T) {
	tests := []struct {
		c    UntypedArray
		want types.Code
	}{
		{
			UntypedArray{nil}, types.UntypedArray,
		},
		{
			UntypedArray{}, types.UntypedArray,
		},
		{
			UntypedArray{Int(1)}, types.UntypedArray,
		},
		{
			UntypedArray{Int(1), Int(2)}, types.UntypedArray,
		},
		{
			UntypedArray{nil}, types.UntypedArray,
		},
		{
			UntypedArray{Int(1), Float(1)}, types.UntypedArray,
		},
		{
			UntypedArray{Int(1), Bytes([]byte{0x01, 0x02, 0x03})}, types.UntypedArray,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				got := tt.c.Type()

				if got != tt.c.Type() {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestUntypedArray_IsValid(t *testing.T) {
	tests := []struct {
		c       UntypedArray
		wantErr bool
	}{
		{
			UntypedArray{nil}, false,
		},
		{
			UntypedArray{}, false,
		},
		{
			UntypedArray{Int(1)}, false,
		},
		{
			UntypedArray{Int(1), Int(2)}, false,
		},
		{
			UntypedArray{nil}, true,
		},
		{
			UntypedArray{Int(1), Float(1)}, false,
		},
		{
			UntypedArray{Int(1), Bytes([]byte{0x01, 0x02, 0x03})}, false,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.IsValid(); got == nil == tt.wantErr {
					t.Errorf("got %v, want %v", got, tt.wantErr)
				}
			},
		)
	}
}
