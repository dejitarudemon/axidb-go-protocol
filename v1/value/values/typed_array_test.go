package values

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

func TestTypedArray_Size(t *testing.T) {
	tests := []struct {
		c    TypedArray
		want int
	}{
		{
			TypedArray{
				elemType: types.Int,
				elems:    []value.V{},
			},
			5,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
				},
			},
			13,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Int(2),
				},
			},
			21,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems:    nil,
			},
			5,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Float(1),
				},
			},
			21,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Bytes([]byte{0x01, 0x02, 0x03}),
				},
			},
			20,
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

func TestTypedArray_Encode(t *testing.T) {
	tests := []struct {
		c    TypedArray
		want []byte
	}{
		{
			TypedArray{
				elemType: types.Int,
				elems:    []value.V{},
			},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x03},
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
				},
			},
			[]byte{0x00, 0x00, 0x00, 0x01, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Int(2),
				},
			},
			[]byte{0x00, 0x00, 0x00, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
		},
		{
			TypedArray{
				elemType: types.Int,
				elems:    nil,
			},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x03},
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Float(1),
				},
			},
			[]byte{0x00, 0x00, 0x00, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Bytes([]byte{0x01, 0x02, 0x03}),
				},
			},
			[]byte{0x00, 0x00, 0x00, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x03},
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

func TestTypedArray_Type(t *testing.T) {
	tests := []struct {
		c    TypedArray
		want types.Code
	}{
		{
			TypedArray{
				elemType: types.Int,
				elems:    []value.V{},
			},
			types.TypedArray,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
				},
			},
			types.TypedArray,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Int(2),
				},
			},
			types.TypedArray,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems:    nil,
			},
			types.TypedArray,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Float(1),
				},
			},
			types.TypedArray,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Bytes([]byte{0x01, 0x02, 0x03}),
				},
			},
			types.TypedArray,
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

func TestTypedArray_IsValid(t *testing.T) {
	tests := []struct {
		c       TypedArray
		wantErr bool
	}{
		{
			TypedArray{
				elemType: types.Int,
				elems:    []value.V{},
			},
			false,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
				},
			},
			false,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Int(2),
				},
			},
			false,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems:    nil,
			},
			false,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Float(1),
				},
			},
			true,
		},
		{
			TypedArray{
				elemType: types.Int,
				elems: []value.V{
					Int(1),
					Bytes([]byte{0x01, 0x02, 0x03}),
				},
			},
			true,
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
