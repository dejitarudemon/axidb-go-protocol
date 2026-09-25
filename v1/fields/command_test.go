package fields

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestCommand_Size(t *testing.T) {
	tests := []struct {
		c    Command
		want int
	}{
		{Handshake, 1},
		{Answer, 1},
		{Read, 1},
		{Write, 1},
		{Delete, 1},
		{Batch, 1},
		{Ping, 1},
		{Command(7), 1},
		{Command(255), 1},
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

func TestCommand_Encode(t *testing.T) {
	tests := []struct {
		c    Command
		want []byte
	}{
		{Handshake, []byte{0x00}},
		{Answer, []byte{0x01}},
		{Read, []byte{0x02}},
		{Write, []byte{0x03}},
		{Delete, []byte{0x04}},
		{Batch, []byte{0x05}},
		{Ping, []byte{0x06}},
		{Command(7), []byte{0x07}},
		{Command(255), []byte{0xFF}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				buf := buffer.Slice{}
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

func TestCommand_String(t *testing.T) {
	tests := []struct {
		c    Command
		want string
	}{
		{Handshake, "Handshake"},
		{Answer, "Answer"},
		{Read, "Read"},
		{Write, "Write"},
		{Delete, "Delete"},
		{Batch, "Batch"},
		{Ping, "Ping"},
		{Command(7), "Unknown (7)"},
		{Command(255), "Unknown (255)"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.String(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCommand_IsValid(t *testing.T) {
	tests := []struct {
		c    Command
		want bool
	}{
		{Handshake, true},
		{Answer, true},
		{Read, true},
		{Write, true},
		{Delete, true},
		{Batch, true},
		{Ping, true},
		{Command(7), false},
		{Command(255), false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.IsValid(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
