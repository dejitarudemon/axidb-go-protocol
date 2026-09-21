package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
)

func TestPingAnswer_Size(t *testing.T) {
	tests := []struct {
		p    PingAnswer
		want int
	}{
		{PingAnswer{}, 1},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestPingAnswer_Size %v", tt.p),
			func(t *testing.T) {
				if got := tt.p.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestPingAnswer_Encode(t *testing.T) {
	tests := []struct {
		p    PingAnswer
		want []byte
	}{
		{PingAnswer{}, []byte{0x01}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestPingAnswer_Encode %v", tt.p),
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

func TestPingAnswer_Command(t *testing.T) {
	tests := []struct {
		p    PingAnswer
		want command.Code
	}{
		{PingAnswer{}, command.Answer},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestPingAnswer_Command %v", tt.p),
			func(t *testing.T) {
				if got := tt.p.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestPingAnswer_IsValid(t *testing.T) {
	tests := []struct {
		p    PingAnswer
		want bool
	}{
		{PingAnswer{}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestPingAnswer_IsValid %v", tt.p),
			func(t *testing.T) {
				if got := tt.p.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
