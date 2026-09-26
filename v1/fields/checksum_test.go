package fields

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestChecksum(t *testing.T) {
	tests := []struct {
		name string
		c    Checksum
		want []byte
	}{
		{"zero", Checksum(0), []byte{0x00, 0x00, 0x00, 0x00}},
		{"small", Checksum(3), []byte{0x00, 0x00, 0x00, 0x03}},
		{"max", Checksum(math.MaxUint32), []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.AssertEncoded(t, tt.c, tt.want)
		})
	}
}

func TestNewChecksum(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want Checksum
	}{
		{"nil", nil, 0},
		{"empty", []byte{}, 0},
		{"three bytes", []byte{0x00, 0x01, 0x02}, 0x92FD4BFA},
		{"check value", []byte("123456789"), 0xE3069283},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChecksum(tt.data); got != tt.want {
				t.Errorf("NewChecksum() = %#08X, want %#08X", got, tt.want)
			}
		})
	}
}

func TestNewChecksumWithParts(t *testing.T) {
	tests := []struct {
		name  string
		base  []byte
		parts [][]byte
		want  Checksum
	}{
		{"empty", nil, nil, 0},
		{"base only", []byte("123456789"), nil, 0xE3069283},
		{"empty base", nil, [][]byte{[]byte("123456789")}, 0xE3069283},
		{"one part", []byte{0x00}, [][]byte{{0x01, 0x02}}, 0x92FD4BFA},
		{"many parts", []byte("1"), [][]byte{[]byte("2345"), {}, []byte("6789")}, 0xE3069283},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChecksumWithParts(tt.base, tt.parts...); got != tt.want {
				t.Errorf("NewChecksumWithParts() = %#08X, want %#08X", got, tt.want)
			}
		})
	}
}

func TestChecksum_Equal(t *testing.T) {
	tests := []struct {
		name string
		a, b Checksum
		want bool
	}{
		{"zero", 0, 0, true},
		{"same", 0xE3069283, 0xE3069283, true},
		{"different", 0xE3069283, 0x92FD4BFA, false},
		{"zero and max", 0, math.MaxUint32, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.b); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
