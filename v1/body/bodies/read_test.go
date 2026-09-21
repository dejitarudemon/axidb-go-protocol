package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestRead_Size(t *testing.T) {
	tests := []struct {
		r    Read
		want int
	}{
		{Read{}, 0},
		{Read{[]byte{0x00}}, 1},
		{Read{[]byte("ab")}, 2},
		{Read{[]byte("фи")}, 4},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestRead_Size %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRead_Encode(t *testing.T) {
	tests := []struct {
		r    Read
		want []byte
	}{
		{Read{}, []byte{}},
		{Read{[]byte{0x00}}, []byte{0x00}},
		{Read{[]byte("ab")}, []byte{0x61, 0x62}},
		{Read{[]byte("фи")}, []byte{0xD1, 0x84, 0xD0, 0xB8}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestRead_Encode %v", tt.r),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.r.Size())

				tt.r.Encode(&buf)

				got := buf.Bytes()

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestRead_Command(t *testing.T) {
	tests := []struct {
		r    Read
		want fields.Command
	}{
		{Read{}, fields.Read},
		{Read{[]byte{0x00}}, fields.Read},
		{Read{[]byte("ab")}, fields.Read},
		{Read{[]byte("фи")}, fields.Read},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestRead_Command %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRead_IsValid(t *testing.T) {
	tests := []struct {
		r    Read
		want bool
	}{
		{Read{}, true},
		{Read{[]byte{0x00}}, false},
		{Read{[]byte("ab")}, false},
		{Read{[]byte("фи")}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestRead_IsValid %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
