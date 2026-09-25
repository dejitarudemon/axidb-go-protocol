package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestWriteAnswer_Size(t *testing.T) {
	tests := []struct {
		w    WriteAnswer
		want int
	}{
		{WriteAnswer{}, 2},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestWriteAnswer_Size %v", tt.w),
			func(t *testing.T) {
				if got := tt.w.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestWriteAnswer_Encode(t *testing.T) {
	tests := []struct {
		w    WriteAnswer
		want []byte
	}{
		{WriteAnswer{}, []byte{0x01, 0x03}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestWriteAnswer_Encode %v", tt.w),
			func(t *testing.T) {
				buf := buffer.Slice{}
				buf.Preallocate(tt.w.Size())

				tt.w.Encode(&buf)

				got := buf.Bytes()

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestWriteAnswer_Command(t *testing.T) {
	tests := []struct {
		w    WriteAnswer
		want fields.Command
	}{
		{WriteAnswer{}, fields.Answer},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestWriteAnswer_Command %v", tt.w),
			func(t *testing.T) {
				if got := tt.w.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestWriteAnswer_IsValid(t *testing.T) {
	tests := []struct {
		w    WriteAnswer
		want bool
	}{
		{WriteAnswer{}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestWriteAnswer_IsValid %v", tt.w),
			func(t *testing.T) {
				if got := tt.w.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
