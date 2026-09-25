package fields

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestTracebackID_Size(t *testing.T) {
	tests := []struct {
		t    TracebackID
		want int
	}{
		{TracebackID([16]byte{}), 16},
		{TracebackID([16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}), 16},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.t),
			func(t *testing.T) {
				if got := tt.t.Size(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestTracebackID_Encode(t *testing.T) {
	tests := []struct {
		t    TracebackID
		want []byte
	}{
		{TracebackID([16]byte{}), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{TracebackID([16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.t),
			func(t *testing.T) {
				buf := buffer.Slice{}
				buf.Preallocate(tt.t.Size())

				tt.t.Encode(&buf)
				got := buf.Bytes()

				if len(got) != tt.t.Size() {
					t.Fatalf("expected %v bytes, got %v bytes", tt.t.Size(), len(got))
				}

				if !bytes.Equal(got, tt.want) {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			},
		)
	}
}
