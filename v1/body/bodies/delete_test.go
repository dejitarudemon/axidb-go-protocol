package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestDelete_Size(t *testing.T) {
	tests := []struct {
		d    Delete
		want int
	}{
		{Delete{}, 0},
		{Delete{0x00}, 1},
		{Delete("ab"), 2},
		{Delete("фи"), 4},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDelete_Size %v", tt.d),
			func(t *testing.T) {
				if got := tt.d.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestDelete_Encode(t *testing.T) {
	tests := []struct {
		d    Delete
		want []byte
	}{
		{Delete{}, []byte{}},
		{Delete{0x00}, []byte{0x00}},
		{Delete("ab"), []byte{0x61, 0x62}},
		{Delete("фи"), []byte{0xD1, 0x84, 0xD0, 0xB8}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDelete_Encode %v", tt.d),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.d.Size())

				tt.d.Encode(&buf)

				got := buf.Bytes()

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestDelete_Command(t *testing.T) {
	tests := []struct {
		d    Delete
		want fields.Command
	}{
		{Delete{}, fields.Delete},
		{Delete{0x00}, fields.Delete},
		{Delete("ab"), fields.Delete},
		{Delete("фи"), fields.Delete},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDelete_Command %v", tt.d),
			func(t *testing.T) {
				if got := tt.d.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestDelete_IsValid(t *testing.T) {
	tests := []struct {
		d    Delete
		want bool
	}{
		{Delete{}, true},
		{Delete{0x00}, false},
		{Delete("ab"), false},
		{Delete("фи"), false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDelete_IsValid %v", tt.d),
			func(t *testing.T) {
				if got := tt.d.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
