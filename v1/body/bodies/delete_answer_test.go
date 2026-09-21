package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestDeleteAnswer_Size(t *testing.T) {
	tests := []struct {
		d    DeleteAnswer
		want int
	}{
		{DeleteAnswer{}, 1},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDeleteAnswer_Size %v", tt.d),
			func(t *testing.T) {
				if got := tt.d.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestDeleteAnswer_Encode(t *testing.T) {
	tests := []struct {
		d    DeleteAnswer
		want []byte
	}{
		{DeleteAnswer{}, []byte{0x01}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDeleteAnswer_Encode %v", tt.d),
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

func TestDeleteAnswer_Command(t *testing.T) {
	tests := []struct {
		d    DeleteAnswer
		want fields.Command
	}{
		{DeleteAnswer{}, fields.Answer},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDeleteAnswer_Command %v", tt.d),
			func(t *testing.T) {
				if got := tt.d.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestDeleteAnswer_IsValid(t *testing.T) {
	tests := []struct {
		d    DeleteAnswer
		want bool
	}{
		{DeleteAnswer{}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDeleteAnswer_IsValid %v", tt.d),
			func(t *testing.T) {
				if got := tt.d.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
