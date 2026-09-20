package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
)

func TestDelete_Size(t *testing.T) {
	tests := []struct {
		d    Delete
		want int
	}{
		{Delete{}, 0},
		{Delete{[]byte{0x00}}, 1},
		{Delete{[]byte("ab")}, 2},
		{Delete{[]byte("фи")}, 4},
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
		{Delete{[]byte{0x00}}, []byte{0x00}},
		{Delete{[]byte("ab")}, []byte{0x61, 0x62}},
		{Delete{[]byte("фи")}, []byte{0xD1, 0x84, 0xD0, 0xB8}},
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
		want command.Code
	}{
		{Delete{}, command.Delete},
		{Delete{[]byte{0x00}}, command.Delete},
		{Delete{[]byte("ab")}, command.Delete},
		{Delete{[]byte("фи")}, command.Delete},
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
		{Delete{}, false},
		{Delete{[]byte{0x00}}, true},
		{Delete{[]byte("ab")}, true},
		{Delete{[]byte("фи")}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDelete_IsValid %v", tt.d),
			func(t *testing.T) {
				if got := tt.d.IsValid(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
