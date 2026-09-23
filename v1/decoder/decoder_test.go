package decoder

import (
	"bufio"
	"bytes"
	"fmt"
	"slices"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func TestDecoder_BySpecs(t *testing.T) {
	tests := []struct {
		encoded []byte
		decoded frame.Frame
	}{
		{
			encoded: []byte{
				0x0A, 0xDB, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x2B, 0x00, 0x00, 0x00,
				0x04, 0x75, 0x73, 0x65, 0x72, 0x01, 0x02, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x01, 0x02,
				0x7A, 0x39, 0x15, 0xBF},
			decoded: frame.Frame{
				RequestID: fields.RequestID(0),
				Body: bodies.Handshake{
					Login:        "user",
					Hash:         [32]byte{0x01, 0x02},
					Compressions: []fields.Compression{fields.Zstd, fields.S2},
				},
			},
		},
		{
			encoded: []byte{
				0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x04, 0x01, 0x00, 0x01,
				0x01, 0x1B, 0x08, 0x33, 0x19},
			decoded: frame.Frame{
				RequestID: fields.RequestID(0),
				Body: bodies.HandshakeAnswer{
					Compressions: []fields.Compression{fields.Zstd},
				},
			},
		},
		{
			encoded: []byte{
				0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x02,
				0x6B, 0xAE, 0x36, 0x80, 0x76, 0x95, 0x47, 0x42,
				0x92, 0xAF, 0x0A, 0x6F, 0xEC, 0x33, 0x98, 0x25,
				0x00, 0xF0, 0x8E, 0x28, 0x51},
			decoded: frame.Frame{
				RequestID: fields.RequestID(0),
				Body: bodies.ErrorAnswer{
					Err: errs.NewErrorUnexpectedCommandWithTracebackID(
						fields.Batch,
						fields.Handshake,
						fields.TracebackID{0x6B, 0xAE, 0x36, 0x80, 0x76, 0x95, 0x47, 0x42, 0x92, 0xAF, 0xA, 0x6F, 0xEC, 0x33, 0x98, 0x25},
					),
				},
			},
		},
		{
			encoded: []byte{
				0x0A, 0xDB, 0x01, 0x02, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x08, 0x73, 0x6F, 0x6D,
				0x65, 0x2D, 0x6B, 0x65, 0x79, 0xE1, 0x5A, 0x14,
				0x00},
			decoded: frame.Frame{
				RequestID: fields.RequestID(1),
				Body:      bodies.Read("some-key"),
			},
		},
		{
			encoded: []byte{
				0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x10, 0x01, 0x02, 0x06,
				0x00, 0x00, 0x00, 0x09, 0x73, 0x6F, 0x6D, 0x65,
				0x2D, 0x64, 0x61, 0x74, 0x61, 0x5E, 0x79, 0xB8,
				0x83},
			decoded: frame.Frame{
				RequestID: fields.RequestID(1),
				Body:      bodies.ReadAnswer{Value: values.String("some-data")},
			},
		},
		{
			encoded: []byte{
				0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x10, 0x01, 0x02, 0x06,
				0x00, 0x00, 0x00, 0x09, 0x73, 0x6F, 0x6D, 0x65,
				0x2D, 0x64, 0x61, 0x74, 0x61, 0x5E, 0x79, 0xB8,
				0x83},
			decoded: frame.Frame{
				RequestID: fields.RequestID(1),
				Body:      bodies.ReadAnswer{Value: values.String("some-data")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDecoder_BySpecs_%v", i),
			func(t *testing.T) {
				decoder := NewDecoder(1024, nil)
				reader := bufio.NewReader(bytes.NewReader(tt.encoded))
				version, err := decoder.DecodePreamble(reader)

				if err != nil {
					t.Fatalf("DecodePreamble: got err %v", err)
					return
				}

				if version != fields.Version(1) {
					t.Fatalf("DecodePreamble: got unexpected version %v", version)
					return
				}

				f, err := decoder.DecodeFrame(reader)
				if err != nil {
					t.Fatalf("DecodeFrame: got err %+v", err)
					return
				}

				if f.RequestID != tt.decoded.RequestID {
					t.Fatalf("Compare Request IDs: got %v, want %v", f.RequestID, tt.decoded.RequestID)
				}

				switch want := tt.decoded.Body.(type) {
				case bodies.Handshake:
					got, ok := f.Body.(bodies.Handshake)
					if !ok {
						t.Errorf("Mismatched typed between got %v and want %v", f.Body, want)
						return
					}

					if got.Login != want.Login {
						t.Errorf("Compare Logins: got %v and want %v", got.Login, want.Login)
					}

					if !bytes.Equal(got.Hash[:], want.Hash[:]) {
						t.Errorf("Compare Hashes: got % X and want % X", got.Hash, want.Hash)
					}

					if !slices.Equal(got.Compressions, want.Compressions) {
						t.Errorf("Compare Compressions: got %v and want %v", got.Compressions, want.Compressions)
					}
				case bodies.HandshakeAnswer:
					got, ok := f.Body.(bodies.HandshakeAnswer)
					if !ok {
						t.Errorf("Mismatched typed between got %v and want %v", f.Body, want)
						return
					}

					if !slices.Equal(got.Compressions, want.Compressions) {
						t.Errorf("Compare Compressions: got %v and want %v", got.Compressions, want.Compressions)
					}
				case bodies.ErrorAnswer:
					got, ok := f.Body.(bodies.ErrorAnswer)
					if !ok {
						t.Errorf("Mismatched typed between got %v and want %v", f.Body, want)
						return
					}

					if want.Err.TracebackID() != got.Err.TracebackID() {
						t.Errorf("Compare Traceback ID: got % X and want % X", got.Err.TracebackID(), want.Err.TracebackID())
					}
				case bodies.Read:
					got, ok := f.Body.(bodies.Read)
					if !ok {
						t.Errorf("Mismatched typed between got %v and want %v", f.Body, want)
						return
					}

					if !bytes.Equal([]byte(got), []byte(want)) {
						t.Errorf("Compare Keys: got % X and want % X", got, want)
					}

				case bodies.ReadAnswer:
					got, ok := f.Body.(bodies.ReadAnswer)
					if !ok {
						t.Errorf("Mismatched typed between got %v and want %v", f.Body, want)
						return
					}

					if got.Value != want.Value {
						t.Errorf("Compare Values: got %v and want %v", got.Value, want.Value)
					}
				}
			},
		)
	}

}
