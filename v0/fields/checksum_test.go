package fields

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v0/internal/testutil"
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
		{"check value", []byte("123456789"), 0xBD0BE338},
		{"hello example", []byte{0x0A, 0xDB, 0x00, 0x03, 0x01, 0x02, 0x03}, 0x6E389900},
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
		{"base only", []byte("123456789"), nil, 0xBD0BE338},
		{"empty base", nil, [][]byte{[]byte("123456789")}, 0xBD0BE338},
		{"split check value", []byte("123"), [][]byte{[]byte("456"), {}, []byte("789")}, 0xBD0BE338},
		{"hello parts", []byte{0x0A, 0xDB, 0x00}, [][]byte{{0x03, 0x01, 0x02, 0x03}}, 0x6E389900},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChecksumWithParts(tt.base, tt.parts...); got != tt.want {
				t.Errorf("NewChecksumWithParts() = %#08X, want %#08X", got, tt.want)
			}
		})
	}
}

func crc32xfer(data []byte) uint32 {
	var crc uint32
	for _, b := range data {
		crc ^= uint32(b) << 24
		for range 8 {
			if crc&0x80000000 != 0 {
				crc = crc<<1 ^ 0x000000AF
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

func TestNewChecksum_MatchesBitAlgorithm(t *testing.T) {
	tests := [][]byte{
		nil,
		{},
		[]byte("123456789"),
		{0x0A, 0xDB, 0x00, 0x03, 0x01, 0x02, 0x03},
		{0x0A, 0xDB, 0x00, 0x04, 0x01, 0x04, 0x07, 0x0B},
	}

	for _, data := range tests {
		if got, want := NewChecksum(data), Checksum(crc32xfer(data)); got != want {
			t.Errorf("NewChecksum(% X) = %#08X, want %#08X", data, got, want)
		}
	}
}

func TestChecksum_Equal(t *testing.T) {
	tests := []struct {
		name string
		a, b Checksum
		want bool
	}{
		{"zero", 0, 0, true},
		{"same", 0xBD0BE338, 0xBD0BE338, true},
		{"different", 0xBD0BE338, 0x6E389900, false},
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
