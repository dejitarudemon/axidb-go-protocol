package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
)

func TestPing_Size(t *testing.T) {
	tests := []struct {
		p    Ping
		want int
	}{
		{Ping{}, 0},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestPing_Size %v", tt.p),
			func(t *testing.T) {
				if got := tt.p.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestPing_Encodt(t *testing.T) {
	tests := []struct {
		p    Ping
		want []byte
	}{
		{Ping{}, []byte{}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestPing_Encode %v", tt.p),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.p.Size())

				tt.p.Encode(&buf)

				got := buf.Bytes()

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestPing_Command(t *testing.T) {
	tests := []struct {
		p    Ping
		want command.Code
	}{
		{Ping{}, command.Ping},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestPing_Command %v", tt.p),
			func(t *testing.T) {
				if got := tt.p.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestPing_IsValid(t *testing.T) {
	tests := []struct {
		p    Ping
		want bool
	}{
		{Ping{}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestPing_IsValid %v", tt.p),
			func(t *testing.T) {
				if got := tt.p.IsValid(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
