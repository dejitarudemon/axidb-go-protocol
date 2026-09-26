package bodies

import (
	"slices"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestHandshakeAnswer(t *testing.T) {
	tests := []struct {
		name    string
		b       HandshakeAnswer
		want    []byte
		wantErr bool
	}{
		{"empty", HandshakeAnswer{}, []byte{0x01, 0x00, 0x00}, false},
		{"empty compressions", HandshakeAnswer{[]fields.Compression{}}, []byte{0x01, 0x00, 0x00}, false},
		{"none", HandshakeAnswer{[]fields.Compression{fields.None}}, []byte{0x01, 0x00, 0x01, 0x00}, false},
		{"two compressions", HandshakeAnswer{[]fields.Compression{fields.None, fields.Zstd}}, []byte{0x01, 0x00, 0x02, 0x00, 0x01}, false},
		{"repeated compressions", HandshakeAnswer{[]fields.Compression{fields.None, fields.Zstd, fields.None}}, []byte{0x01, 0x00, 0x03, 0x00, 0x01, 0x00}, false},
		{"custom compression", HandshakeAnswer{[]fields.Compression{fields.None, fields.Zstd, 0xFF}}, []byte{0x01, 0x00, 0x03, 0x00, 0x01, 0xFF}, false},
		{
			"max compressions",
			HandshakeAnswer{generateManyCompressions(MaxCompressionsPerOneHandshake)},
			append([]byte{0x01, 0x00}, encodeCompressions(generateManyCompressions(MaxCompressionsPerOneHandshake))...),
			false,
		},
		{
			"too many compressions",
			HandshakeAnswer{generateManyCompressions(1000)},
			append([]byte{0x01, 0x00}, encodeCompressions(generateManyCompressions(1000))...),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertAnswer(t, tt.b, fields.Handshake, tt.want, tt.wantErr)
		})
	}
}

func TestNewHandshakeAnswer(t *testing.T) {
	tests := []struct {
		name string
		got  HandshakeAnswer
		want HandshakeAnswer
	}{
		{"nil compressions", NewHandshakeAnswer(nil), HandshakeAnswer{}},
		{
			"drops none",
			NewHandshakeAnswer([]fields.Compression{0x01, 0x00, 0x02}),
			HandshakeAnswer{[]fields.Compression{0x01, 0x02}},
		},
		{
			"drops none and duplicates",
			NewHandshakeAnswer([]fields.Compression{0x01, 0x00, 0x02, 0x01, 0x03}),
			HandshakeAnswer{[]fields.Compression{0x01, 0x02, 0x03}},
		},
		{
			"drops duplicates",
			NewHandshakeAnswer([]fields.Compression{0x01, 0x01}),
			HandshakeAnswer{[]fields.Compression{0x01}},
		},
		{
			"drops duplicates beyond limit",
			NewHandshakeAnswer(generateManyCompressions(1000)),
			HandshakeAnswer{generateManyCompressions(254)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !slices.Equal(tt.got.Compressions, tt.want.Compressions) {
				t.Errorf("Compressions = %v, want %v", tt.got.Compressions, tt.want.Compressions)
			}

			testutil.AssertSameEncoding(t, tt.got, tt.want)
			testutil.AssertErr(t, tt.got.IsValid(), false)
		})
	}
}
