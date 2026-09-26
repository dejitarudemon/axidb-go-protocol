package frame_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"slices"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor/compressors"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/specs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

type failingCompressor struct{ err error }

func (f failingCompressor) Code() fields.Compression          { return fields.Zstd }
func (f failingCompressor) Compress([]byte) ([]byte, error)   { return nil, f.err }
func (f failingCompressor) Decompress([]byte) ([]byte, error) { return nil, f.err }

func encode(t *testing.T, f frame.Frame, c compressor.Compressor) []byte {
	t.Helper()

	buf := buffer.Slice{}
	buf.Preallocate(f.Size())

	if e := f.Encode(&buf, c); e != nil {
		t.Fatalf("Encode() = %v", e)
	}

	return buf.Bytes()
}

func TestFrame_Specs(t *testing.T) {
	for _, tt := range specs.Frames {
		t.Run(tt.Name, func(t *testing.T) {
			if got := encode(t, tt.Frame, nil); !bytes.Equal(got, tt.Encoded) {
				t.Errorf("Encode() = % X, want % X", got, tt.Encoded)
			}

			if got := tt.Frame.Size(); got != len(tt.Encoded) {
				t.Errorf("Size() = %v, want %v", got, len(tt.Encoded))
			}

			if got := tt.Frame.Version(); got != fields.CurrentVersion {
				t.Errorf("Version() = %v, want %v", got, fields.CurrentVersion)
			}

			if e := tt.Frame.IsValid(); e != nil {
				t.Errorf("IsValid() = %v", e)
			}
		})
	}
}

func TestFrame_EncodeWithCompression(t *testing.T) {
	zstd, _ := compressors.NewZstd(1 << 12)
	s2, _ := compressors.NewS2(1 << 12)

	tests := []struct {
		name       string
		compressor compressor.Compressor
		f          frame.Frame
	}{
		{"zstd read", zstd, frame.Frame{RequestID: 1, Body: bodies.Read("key")}},
		{"zstd read answer", zstd, frame.Frame{RequestID: 2, Body: bodies.ReadAnswer{Value: values.Bytes(make([]byte, 256))}}},
		{"s2 read", s2, frame.Frame{RequestID: 1, Body: bodies.Read("key")}},
		{"s2 read answer", s2, frame.Frame{RequestID: 2, Body: bodies.ReadAnswer{Value: values.Bytes(make([]byte, 256))}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encode(t, tt.f, tt.compressor)

			headers, compressed, checksum := got[:13], got[13:len(got)-4], got[len(got)-4:]

			if want := byte(tt.f.Body.Command()); headers[3] != want {
				t.Errorf("command byte = %v, want %v", headers[3], want)
			}

			if want := byte(tt.compressor.Code()); headers[8] != want {
				t.Errorf("compression byte = %v, want %v", headers[8], want)
			}

			if bodyLen := binary.BigEndian.Uint32(headers[9:]); int(bodyLen) != len(compressed) {
				t.Errorf("body len = %v, but body has %v bytes", bodyLen, len(compressed))
			}

			if want := fields.NewChecksum(got[:len(got)-4]); binary.BigEndian.Uint32(checksum) != uint32(want) {
				t.Errorf("checksum = % X, want %08X", checksum, uint32(want))
			}

			body, e := tt.compressor.Decompress(compressed)
			if e != nil {
				t.Fatalf("Decompress() = %v", e)
			}

			if want := testutil.Encode(t, tt.f.Body); !bytes.Equal(body, want) {
				t.Errorf("decompressed body = % X, want % X", body, want)
			}
		})
	}
}

func TestFrame_EncodeErrs(t *testing.T) {
	sentinel := errors.New("compress failed")

	tests := []struct {
		name       string
		f          frame.Frame
		compressor compressor.Compressor
		want       error
	}{
		{"nil body", frame.Frame{RequestID: 1}, nil, nil},
		{"nil body with compressor", frame.Frame{RequestID: 1}, failingCompressor{sentinel}, nil},
		{"compressor error", frame.Frame{RequestID: 1, Body: bodies.Ping{}}, failingCompressor{sentinel}, sentinel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Slice{}

			e := tt.f.Encode(&buf, tt.compressor)
			if e == nil {
				t.Fatal("Encode() = nil, want error")
			}

			if tt.want != nil && !errors.Is(e, tt.want) {
				t.Errorf("Encode() = %v, want %v", e, tt.want)
			}
		})
	}
}

func TestFrame_NilBody(t *testing.T) {
	f := frame.Frame{RequestID: 1}

	if got, want := f.Size(), 17; got != want {
		t.Errorf("Size() = %v, want %v", got, want)
	}

	if e := f.IsValid(); e == nil {
		t.Error("IsValid() = nil, want error")
	}
}

func TestFrame_IsValid(t *testing.T) {
	tooManyCompressions := slices.Repeat([]fields.Compression{fields.Zstd}, int(bodies.MaxCompressionsPerOneHandshake)+1)

	tests := []struct {
		name string
		f    frame.Frame
	}{
		{"handshake with too many compressions", frame.Frame{RequestID: 1, Body: bodies.Handshake{Compressions: tooManyCompressions}}},
		{"handshake answer with too many compressions", frame.Frame{RequestID: 1, Body: bodies.HandshakeAnswer{Compressions: tooManyCompressions}}},
		{"read empty key", frame.Frame{RequestID: 1, Body: bodies.Read("")}},
		{"read answer nil value", frame.Frame{RequestID: 1, Body: bodies.ReadAnswer{Value: nil}}},
		{"read answer invalid json", frame.Frame{RequestID: 1, Body: bodies.ReadAnswer{Value: values.JSON("n")}}},
		{"write empty key and nil value", frame.Frame{RequestID: 1, Body: bodies.Write{Key: fields.Key(""), Value: nil}}},
		{"write nil value", frame.Frame{RequestID: 1, Body: bodies.Write{Key: fields.Key("key"), Value: nil}}},
		{"write invalid json", frame.Frame{RequestID: 1, Body: bodies.Write{Key: fields.Key("key"), Value: values.JSON("n")}}},
		{"delete empty key", frame.Frame{RequestID: 1, Body: bodies.Delete("")}},
		{"batch with nil request body", frame.Frame{RequestID: 1, Body: bodies.Batch{Requests: []bodies.Request{{Number: 1, Body: nil}}}}},
		{"batch with invalid request", frame.Frame{RequestID: 1, Body: bodies.Batch{Requests: []bodies.Request{{Number: 1, Body: bodies.Delete("")}}}}},
		{"batch answer with nil result body", frame.Frame{RequestID: 1, Body: bodies.BatchAnswer{{Number: 1, Body: nil}}}},
		{"batch answer with nil value", frame.Frame{RequestID: 1, Body: bodies.BatchAnswer{{Number: 1, Body: bodies.ReadAnswer{Value: nil}}}}},
		{"batch answer with invalid json", frame.Frame{RequestID: 1, Body: bodies.BatchAnswer{{Number: 1, Body: bodies.ReadAnswer{Value: values.JSON("a")}}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if e := tt.f.IsValid(); e == nil {
				t.Error("IsValid() = nil, want error")
			}
		})
	}
}
