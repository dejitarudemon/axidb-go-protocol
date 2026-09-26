package bodies

import (
	"slices"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func generateManyCompressions(size uint16) []fields.Compression {
	compressions := make([]fields.Compression, 0, size)

	for i := range size {
		compressions = append(compressions, fields.Compression((i+1)%254+1))
	}

	return compressions
}

func encodeCompressions(compressions []fields.Compression) []byte {
	n := min(len(compressions), MaxCompressionsPerOneHandshake)
	encoded := []byte{byte(n)}

	for _, c := range compressions[:n] {
		encoded = append(encoded, byte(c))
	}

	return encoded
}

func hashBytes(prefix ...byte) []byte {
	var hash [HashFieldSize]byte
	copy(hash[:], prefix)
	return hash[:]
}

func TestHandshake(t *testing.T) {
	user := []byte{0x00, 0x00, 0x00, 0x04, 0x75, 0x73, 0x65, 0x72}

	tests := []struct {
		name    string
		b       Handshake
		want    []byte
		wantErr bool
	}{
		{
			"empty",
			Handshake{},
			slices.Concat([]byte{0x00, 0x00, 0x00, 0x00}, hashBytes(), []byte{0x00}),
			false,
		},
		{
			"empty compressions",
			Handshake{"", [32]byte{}, []fields.Compression{}},
			slices.Concat([]byte{0x00, 0x00, 0x00, 0x00}, hashBytes(), []byte{0x00}),
			false,
		},
		{
			"ascii login",
			Handshake{"user", [32]byte{0x01}, []fields.Compression{}},
			slices.Concat(user, hashBytes(0x01), []byte{0x00}),
			false,
		},
		{
			"utf-8 login",
			Handshake{"юзер", [32]byte{0x01, 0x02}, []fields.Compression{}},
			slices.Concat(
				[]byte{0x00, 0x00, 0x00, 0x08, 0xD1, 0x8E, 0xD0, 0xB7, 0xD0, 0xB5, 0xD1, 0x80},
				hashBytes(0x01, 0x02),
				[]byte{0x00},
			),
			false,
		},
		{
			"none compression",
			Handshake{"user", [32]byte{0x01, 0x02, 0x03}, []fields.Compression{fields.None}},
			slices.Concat(user, hashBytes(0x01, 0x02, 0x03), []byte{0x01, 0x00}),
			false,
		},
		{
			"two compressions",
			Handshake{"user", [32]byte{0x01, 0x02, 0x03, 0x04}, []fields.Compression{fields.None, fields.Zstd}},
			slices.Concat(user, hashBytes(0x01, 0x02, 0x03, 0x04), []byte{0x02, 0x00, 0x01}),
			false,
		},
		{
			"repeated compressions",
			Handshake{"user", [32]byte{0x01, 0x02, 0x03}, []fields.Compression{fields.None, fields.Zstd, fields.None}},
			slices.Concat(user, hashBytes(0x01, 0x02, 0x03), []byte{0x03, 0x00, 0x01, 0x00}),
			false,
		},
		{
			"custom compression",
			Handshake{"user", [32]byte{}, []fields.Compression{fields.None, fields.Zstd, 0xFF}},
			slices.Concat(user, hashBytes(), []byte{0x03, 0x00, 0x01, 0xFF}),
			false,
		},
		{
			"max compressions",
			Handshake{"user", [32]byte{0x01}, generateManyCompressions(MaxCompressionsPerOneHandshake)},
			slices.Concat(user, hashBytes(0x01), encodeCompressions(generateManyCompressions(MaxCompressionsPerOneHandshake))),
			false,
		},
		{
			"too many compressions",
			Handshake{"user", [32]byte{0x01}, generateManyCompressions(1000)},
			slices.Concat(user, hashBytes(0x01), encodeCompressions(generateManyCompressions(1000))),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBody(t, tt.b, fields.Handshake, tt.want, tt.wantErr)
		})
	}
}

func TestNewHandshake(t *testing.T) {
	tests := []struct {
		name string
		got  Handshake
		want Handshake
	}{
		{"nil compressions", NewHandshake("", [32]byte{}, nil), Handshake{}},
		{
			"drops none",
			NewHandshake("user", [32]byte{0x01}, []fields.Compression{0x01, 0x00, 0x02}),
			Handshake{"user", [32]byte{0x01}, []fields.Compression{0x01, 0x02}},
		},
		{
			"drops none and duplicates",
			NewHandshake("user", [32]byte{0x02}, []fields.Compression{0x01, 0x00, 0x02, 0x01, 0x03}),
			Handshake{"user", [32]byte{0x02}, []fields.Compression{0x01, 0x02, 0x03}},
		},
		{
			"drops duplicates",
			NewHandshake("user", [32]byte{0x03}, []fields.Compression{0x01, 0x01}),
			Handshake{"user", [32]byte{0x03}, []fields.Compression{0x01}},
		},
		{
			"drops repeated none",
			NewHandshake("user", [32]byte{0x01, 0x02}, []fields.Compression{0x00, 0x01, 0x00}),
			Handshake{"user", [32]byte{0x01, 0x02}, []fields.Compression{0x01}},
		},
		{
			"drops duplicates beyond limit",
			NewHandshake("user", [32]byte{}, generateManyCompressions(1000)),
			Handshake{"user", [32]byte{}, generateManyCompressions(254)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.AssertSameEncoding(t, tt.got, tt.want)
			testutil.AssertErr(t, tt.got.IsValid(), false)
		})
	}
}
