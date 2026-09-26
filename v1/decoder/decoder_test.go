package decoder

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"slices"
	"strings"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor/compressors"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

var zstd, _ = compressors.NewZstd(1 << 12)
var s2, _ = compressors.NewS2(1 << 12)

var testsSpecs = []struct {
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
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x0A,
			0x00, 0x00, 0x00, 0x00, 0x29, 0x01, 0x02, 0x02,
			0x00, 0x00, 0x00, 0x03, 0x03, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x01, 0x06, 0x00, 0x00,
			0x00, 0x0B, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x20,
			0x77, 0x6F, 0x72, 0x6C, 0x64, 0x05, 0x40, 0x00,
			0xCC, 0xCC, 0xCC, 0xCC, 0xCC, 0xCD, 0x34, 0xE8,
			0xE0, 0x75},
		decoded: frame.Frame{
			RequestID: fields.RequestID(10),
			Body: bodies.ReadAnswer{Value: values.UntypedArray{
				values.Int(1),
				values.String("hello world"),
				values.Float(2.1),
			}},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x0A,
			0xF1, 0x38, 0x9C, 0x1F, 0xA2, 0x2F, 0x40, 0x82,
			0xA2, 0x79, 0xD2, 0x01, 0x7C, 0x60, 0xF8, 0x20,
			0x88, 0xEE, 0xCC, 0xAA},
		decoded: frame.Frame{
			RequestID: fields.RequestID(1),
			Body: bodies.ErrorAnswer{
				Err: errs.NewErrorNotFoundWithTracebackID(
					fields.Key("123"),
					fields.TracebackID{0xF1, 0x38, 0x9C, 0x1F, 0xA2, 0x2F, 0x40, 0x82,
						0xA2, 0x79, 0xD2, 0x01, 0x7C, 0x60, 0xF8, 0x20},
				),
			},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x08,
			0xA5, 0x18, 0xB2, 0xC1, 0xCC, 0x27, 0x49, 0x27,
			0x94, 0x09, 0x6A, 0x10, 0x15, 0xDA, 0x67, 0xB1,
			0x9C, 0x90, 0xb4, 0x33},
		decoded: frame.Frame{
			RequestID: fields.RequestID(1),
			Body: bodies.ErrorAnswer{
				Err: errs.NewErrorInternalErrorWithTracebackID(
					nil,
					fields.TracebackID{0xA5, 0x18, 0xB2, 0xC1, 0xCC, 0x27, 0x49, 0x27,
						0x94, 0x09, 0x6A, 0x10, 0x15, 0xDA, 0x67, 0xB1},
				),
			},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x03, 0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00,
			0x04, 0x63, 0x6F, 0x64, 0x65, 0x06, 0x00, 0x00,
			0x00, 0x05, 0x49, 0x44, 0x44, 0x51, 0x44, 0x38,
			0x0A, 0x5A, 0x23},
		decoded: frame.Frame{
			RequestID: fields.RequestID(2),
			Body: bodies.Write{
				Key:   fields.Key("code"),
				Value: values.String("IDDQD"),
			},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x03, 0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00,
			0x06, 0x76, 0x65, 0x63, 0x74, 0x6F, 0x72, 0x01,
			0x00, 0x00, 0x00, 0x03, 0x03, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x03, 0xD7, 0x99, 0xDF,
			0x3D},
		decoded: frame.Frame{
			RequestID: 1,
			Body: bodies.Write{
				Key: fields.Key("vector"),
				Value: values.TypedArray{
					ElemType: fields.Int,
					Elems: []value.V{
						values.Int(1),
						values.Int(2),
						values.Int(3),
					},
				},
			},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x00, 0x00, 0x02, 0x01, 0x03, 0xA0,
			0x05, 0x87, 0x3E},
		decoded: frame.Frame{
			RequestID: 2,
			Body:      bodies.WriteAnswer{},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x08,
			0x77, 0x29, 0x53, 0x66, 0x99, 0x93, 0x4D, 0xC3,
			0x9B, 0x99, 0x19, 0x24, 0x11, 0xEE, 0x5B, 0x8C,
			0x9D, 0x90, 0x1B, 0x23},
		decoded: frame.Frame{
			RequestID: 2,
			Body: bodies.ErrorAnswer{
				Err: errs.NewErrorInternalErrorWithTracebackID(
					nil,
					fields.TracebackID{0x77, 0x29, 0x53, 0x66, 0x99, 0x93, 0x4D, 0xC3, 0x9B, 0x99, 0x19, 0x24, 0x11, 0xEE, 0x5B, 0x8C},
				),
			},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x04, 0x00, 0x00, 0x00, 0x03,
			0x00, 0x00, 0x00, 0x00, 0x0B, 0x61, 0x6E, 0x6F,
			0x74, 0x68, 0x65, 0x72, 0x2D, 0x6B, 0x65, 0x79,
			0xF4, 0xC7, 0xFF, 0x81},
		decoded: frame.Frame{
			RequestID: 3,
			Body:      bodies.Delete("another-key"),
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x03,
			0x00, 0x00, 0x00, 0x00, 0x02, 0x01, 0x04, 0x3D,
			0xF3, 0x9E, 0xF2},
		decoded: frame.Frame{
			RequestID: 3,
			Body:      bodies.DeleteAnswer{},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x03,
			0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x08,
			0xCA, 0xD1, 0xC1, 0x02, 0xCA, 0xC5, 0x4D, 0xDF,
			0x9B, 0x61, 0xEB, 0x00, 0x81, 0x8E, 0x25, 0xD1,
			0x4F, 0x23, 0x68, 0x80},
		decoded: frame.Frame{
			RequestID: 3,
			Body: bodies.ErrorAnswer{
				Err: errs.NewErrorInternalErrorWithTracebackID(
					nil,
					fields.TracebackID{0xCA, 0xD1, 0xC1, 0x02, 0xCA, 0xC5, 0x4D, 0xDF,
						0x9B, 0x61, 0xEB, 0x00, 0x81, 0x8E, 0x25, 0xD1},
				),
			},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x06, 0x00, 0x00, 0x00, 0x04,
			0x00, 0x00, 0x00, 0x00, 0x00, 0xA2, 0x1B, 0x87,
			0x9B},
		decoded: frame.Frame{
			RequestID: 4,
			Body:      bodies.Ping{},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x04,
			0x00, 0x00, 0x00, 0x00, 0x02, 0x01, 0x06, 0x26,
			0x91, 0xEB, 0x01},
		decoded: frame.Frame{
			RequestID: 4,
			Body:      bodies.PingAnswer{},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x04,
			0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x08,
			0x3D, 0x8E, 0xD3, 0xBF, 0xB8, 0xD1, 0x47, 0x45,
			0xA5, 0xDB, 0xF8, 0x50, 0x30, 0x90, 0x90, 0x2A,
			0xC9, 0x04, 0x2A, 0x26},
		decoded: frame.Frame{
			RequestID: 4,
			Body: bodies.ErrorAnswer{
				Err: errs.NewErrorInternalErrorWithTracebackID(
					nil,
					fields.TracebackID{0x3D, 0x8E, 0xD3, 0xBF, 0xB8, 0xD1, 0x47, 0x45, 0xA5, 0xDB, 0xF8, 0x50, 0x30, 0x90, 0x90, 0x2A},
				),
			},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x05, 0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x32, 0x05, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x02, 0x00,
			0x00, 0x00, 0x03, 0x6B, 0x65, 0x79, 0x00, 0x00,
			0x00, 0x02, 0x03, 0x00, 0x00, 0x00, 0x18, 0x00,
			0x00, 0x00, 0x0B, 0x61, 0x6E, 0x6F, 0x74, 0x68,
			0x65, 0x72, 0x2D, 0x6B, 0x65, 0x79, 0x06, 0x00,
			0x00, 0x00, 0x04, 0x64, 0x61, 0x74, 0x61, 0xDE,
			0x84, 0x44, 0x1C},
		decoded: frame.Frame{
			RequestID: 1,
			Body: bodies.Batch{
				IsSequentialExecution: true,
				InterruptAfterError:   false,
				IsOneAnswer:           true,

				Requests: []bodies.Request{
					{
						Number: 1,
						Body:   bodies.Read("key"),
					},
					{
						Number: 2,
						Body: bodies.Write{
							Key:   fields.Key("another-key"),
							Value: values.String("data"),
						},
					},
				},
			},
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x26, 0x01, 0x05, 0x00,
			0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00,
			0x00, 0x00, 0x0E, 0x01, 0x02, 0x06, 0x00, 0x00,
			0x00, 0x07, 0x6D, 0x65, 0x73, 0x73, 0x61, 0x67,
			0x65, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
			0x02, 0x01, 0x03, 0xB7, 0x2D, 0x2D, 0xF4},
		decoded: frame.Frame{
			RequestID: 1,
			Body: bodies.BatchAnswer([]bodies.Result{
				{
					Number: 1,
					Body: bodies.ReadAnswer{
						Value: values.String("message"),
					},
				},
				{
					Number: 2,
					Body:   bodies.WriteAnswer{},
				},
			},
			),
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x3C, 0x01, 0x05, 0x00,
			0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00,
			0x00, 0x00, 0x13, 0x00, 0x00, 0x08, 0x38, 0xDD,
			0x5C, 0x39, 0x24, 0xDB, 0x45, 0xA9, 0x85, 0x09,
			0xFF, 0xE6, 0xC6, 0x5C, 0x7F, 0x42, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00,
			0x0C, 0x49, 0x84, 0x98, 0xC5, 0xCF, 0x19, 0x40,
			0xC9, 0x95, 0x38, 0x15, 0x5C, 0x27, 0xA3, 0xCF,
			0xF1, 0x01, 0x50, 0x3D, 0x7B},
		decoded: frame.Frame{
			RequestID: 1,
			Body: bodies.BatchAnswer([]bodies.Result{
				{
					Number: 1,
					Body: bodies.ErrorAnswer{
						Err: errs.NewErrorInternalErrorWithTracebackID(
							nil,
							fields.TracebackID{0x38, 0xDD, 0x5C, 0x39, 0x24, 0xDB, 0x45, 0xA9, 0x85, 0x09, 0xFF, 0xE6, 0xC6, 0x5C, 0x7F, 0x42},
						),
					},
				},
				{
					Number: 2,
					Body: bodies.ErrorAnswer{
						Err: errs.NewErrorRequestInterruptedWithTracebackID(
							fields.RequestID(1),
							fields.TracebackID{0x49, 0x84, 0x98, 0xC5, 0xCF, 0x19, 0x40, 0xC9, 0x95, 0x38, 0x15, 0x5C, 0x27, 0xA3, 0xCF, 0xF1},
						),
					},
				},
			},
			),
		},
	},
	{
		encoded: []byte{
			0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x08,
			0xE5, 0x3A, 0x39, 0x93, 0x7B, 0x7A, 0x47, 0xEC,
			0xB3, 0x82, 0x7C, 0x2B, 0x05, 0xE3, 0xD6, 0xDD,
			0x53, 0x81, 0xF9, 0xD8},
		decoded: frame.Frame{
			RequestID: 1,
			Body: bodies.ErrorAnswer{
				Err: errs.NewErrorInternalErrorWithTracebackID(
					nil,
					fields.TracebackID{0xE5, 0x3A, 0x39, 0x93, 0x7B, 0x7A, 0x47, 0xEC, 0xB3, 0x82, 0x7C, 0x2B, 0x05, 0xE3, 0xD6, 0xDD},
				),
			},
		},
	},
}

var testsRoundTrip = []struct {
	f          frame.Frame
	compressor compressor.Compressor
}{
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshake("", [32]byte{}, []fields.Compression{})},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{fields.Zstd})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{fields.Zstd, fields.Zstd})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{fields.Zstd, fields.Compression(23)})},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshakeAnswer([]fields.Compression{})},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshakeAnswer([]fields.Compression{})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshakeAnswer([]fields.Compression{fields.Zstd})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshakeAnswer([]fields.Compression{fields.Zstd, fields.Zstd})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshakeAnswer([]fields.Compression{fields.Zstd, fields.Compression(23)})},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Read("")},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Read("key")},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Read("")},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Read("key")},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Bytes{}}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Bytes{0xFF}}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Int(0)}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Int(1)}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{}}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{values.Bytes{}}}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{values.Bytes{0xff}}}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{values.Bytes{0xff}, values.Int(1)}}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.TypedArray{
			Elems:    []value.V{values.Bytes{0xff}, values.Bytes{0x00}},
			ElemType: fields.Bytes,
		}}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Write{Key: fields.Key(""), Value: values.Bytes{}}},
		nil,
	},
	{
		frame.Frame{RequestID: 2, Body: bodies.Write{Key: fields.Key("k"), Value: values.Bytes{0xFF}}},
		nil,
	},
	{
		frame.Frame{RequestID: 3, Body: bodies.Write{Key: fields.Key("ke"), Value: values.Int(0)}},
		nil,
	},
	{
		frame.Frame{RequestID: 4, Body: bodies.Write{Key: fields.Key("key"), Value: values.Int(1)}},
		nil,
	},
	{
		frame.Frame{RequestID: 5, Body: bodies.Write{Key: fields.Key("key "), Value: values.UntypedArray{}}},
		nil,
	},
	{
		frame.Frame{RequestID: 6, Body: bodies.Write{Key: fields.Key("key k"), Value: values.UntypedArray{values.Bytes{}}}},
		nil,
	},
	{
		frame.Frame{RequestID: 7, Body: bodies.Write{Key: fields.Key("key ke"), Value: values.UntypedArray{values.Bytes{0xff}}}},
		nil,
	},
	{
		frame.Frame{RequestID: 8, Body: bodies.Write{Key: fields.Key("key key"), Value: values.UntypedArray{values.Bytes{0xff}, values.Int(1)}}},
		nil,
	},
	{
		frame.Frame{RequestID: 9, Body: bodies.Write{Key: fields.Key("key yek"), Value: values.TypedArray{
			Elems:    []value.V{values.Bytes{0xff}, values.Bytes{0x00}},
			ElemType: fields.Bytes,
		}}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.WriteAnswer{}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.WriteAnswer{}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Delete("")},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Delete("key")},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Delete("")},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Delete("key")},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.DeleteAnswer{}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.DeleteAnswer{}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Ping{}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Ping{}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.PingAnswer{}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.PingAnswer{}},
		nil,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ErrorAnswer{Err: errs.NewErrorInternalError(nil)}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorBatchLimitIsExceeded(1, 0)}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorBatchLimitIsExceeded(0, 1)}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorUnexpectedCommand(fields.Command(10), fields.Handshake)}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorUnexpectedCommand(fields.Handshake, fields.Command(10))}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: false, InterruptAfterError: false, Requests: nil}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: false, InterruptAfterError: false, Requests: nil}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: true, InterruptAfterError: false, Requests: nil}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: false, InterruptAfterError: true, Requests: nil}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: false, Requests: nil}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: true, InterruptAfterError: true, Requests: nil}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: false, InterruptAfterError: true, Requests: nil}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: nil}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
		}}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
			{Number: 0, Body: bodies.Read("key")},
		}}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
			{Number: 1, Body: bodies.Read("key")},
		}}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
			{Number: 1, Body: bodies.Delete("key")},
		}}},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer(nil)},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
		})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 0, Body: bodies.WriteAnswer{}},
		})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 0, Body: bodies.DeleteAnswer{}},
		})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 1, Body: bodies.DeleteAnswer{}},
		})},
		nil,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 1, Body: bodies.ReadAnswer{Value: values.Int(123)}},
		})},
		nil,
	},

	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshake("", [32]byte{}, []fields.Compression{})},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{fields.Zstd})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{fields.Zstd, fields.Zstd})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{fields.Zstd, fields.Compression(23)})},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshakeAnswer([]fields.Compression{})},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshakeAnswer([]fields.Compression{})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshakeAnswer([]fields.Compression{fields.Zstd})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshakeAnswer([]fields.Compression{fields.Zstd, fields.Zstd})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshakeAnswer([]fields.Compression{fields.Zstd, fields.Compression(23)})},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Read("")},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Read("key")},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Read("")},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Read("key")},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Bytes{}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Bytes{0xFF}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Int(0)}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Int(1)}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{values.Bytes{}}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{values.Bytes{0xff}}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{values.Bytes{0xff}, values.Int(1)}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.TypedArray{
			Elems:    []value.V{values.Bytes{0xff}, values.Bytes{0x00}},
			ElemType: fields.Bytes,
		}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Write{Key: fields.Key(""), Value: values.Bytes{}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 2, Body: bodies.Write{Key: fields.Key("k"), Value: values.Bytes{0xFF}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 3, Body: bodies.Write{Key: fields.Key("ke"), Value: values.Int(0)}},
		zstd,
	},
	{
		frame.Frame{RequestID: 4, Body: bodies.Write{Key: fields.Key("key"), Value: values.Int(1)}},
		zstd,
	},
	{
		frame.Frame{RequestID: 5, Body: bodies.Write{Key: fields.Key("key "), Value: values.UntypedArray{}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 6, Body: bodies.Write{Key: fields.Key("key k"), Value: values.UntypedArray{values.Bytes{}}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 7, Body: bodies.Write{Key: fields.Key("key ke"), Value: values.UntypedArray{values.Bytes{0xff}}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 8, Body: bodies.Write{Key: fields.Key("key key"), Value: values.UntypedArray{values.Bytes{0xff}, values.Int(1)}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 9, Body: bodies.Write{Key: fields.Key("key yek"), Value: values.TypedArray{
			Elems:    []value.V{values.Bytes{0xff}, values.Bytes{0x00}},
			ElemType: fields.Bytes,
		}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.WriteAnswer{}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.WriteAnswer{}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Delete("")},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Delete("key")},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Delete("")},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Delete("key")},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.DeleteAnswer{}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.DeleteAnswer{}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Ping{}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Ping{}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.PingAnswer{}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.PingAnswer{}},
		zstd,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ErrorAnswer{Err: errs.NewErrorInternalError(nil)}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorBatchLimitIsExceeded(1, 0)}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorBatchLimitIsExceeded(0, 1)}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorUnexpectedCommand(fields.Command(10), fields.Handshake)}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorUnexpectedCommand(fields.Handshake, fields.Command(10))}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: false, InterruptAfterError: false, Requests: nil}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: false, InterruptAfterError: false, Requests: nil}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: true, InterruptAfterError: false, Requests: nil}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: false, InterruptAfterError: true, Requests: nil}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: false, Requests: nil}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: true, InterruptAfterError: true, Requests: nil}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: false, InterruptAfterError: true, Requests: nil}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: nil}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
		}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
			{Number: 0, Body: bodies.Read("key")},
		}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
			{Number: 1, Body: bodies.Read("key")},
		}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
			{Number: 1, Body: bodies.Delete("key")},
		}}},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer(nil)},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
		})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 0, Body: bodies.WriteAnswer{}},
		})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 0, Body: bodies.DeleteAnswer{}},
		})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 1, Body: bodies.DeleteAnswer{}},
		})},
		zstd,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 1, Body: bodies.ReadAnswer{Value: values.Int(123)}},
		})},
		zstd,
	},

	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshake("", [32]byte{}, []fields.Compression{})},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{fields.Zstd})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{fields.Zstd, fields.Zstd})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshake("test", fakeHash(), []fields.Compression{fields.Zstd, fields.Compression(23)})},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshakeAnswer([]fields.Compression{})},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.NewHandshakeAnswer([]fields.Compression{})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshakeAnswer([]fields.Compression{fields.Zstd})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshakeAnswer([]fields.Compression{fields.Zstd, fields.Zstd})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.NewHandshakeAnswer([]fields.Compression{fields.Zstd, fields.Compression(23)})},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Read("")},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Read("key")},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Read("")},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Read("key")},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Bytes{}}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Bytes{0xFF}}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Int(0)}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.Int(1)}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{}}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{values.Bytes{}}}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{values.Bytes{0xff}}}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.UntypedArray{values.Bytes{0xff}, values.Int(1)}}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ReadAnswer{Value: values.TypedArray{
			Elems:    []value.V{values.Bytes{0xff}, values.Bytes{0x00}},
			ElemType: fields.Bytes,
		}}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Write{Key: fields.Key(""), Value: values.Bytes{}}},
		s2,
	},
	{
		frame.Frame{RequestID: 2, Body: bodies.Write{Key: fields.Key("k"), Value: values.Bytes{0xFF}}},
		s2,
	},
	{
		frame.Frame{RequestID: 3, Body: bodies.Write{Key: fields.Key("ke"), Value: values.Int(0)}},
		s2,
	},
	{
		frame.Frame{RequestID: 4, Body: bodies.Write{Key: fields.Key("key"), Value: values.Int(1)}},
		s2,
	},
	{
		frame.Frame{RequestID: 5, Body: bodies.Write{Key: fields.Key("key "), Value: values.UntypedArray{}}},
		s2,
	},
	{
		frame.Frame{RequestID: 6, Body: bodies.Write{Key: fields.Key("key k"), Value: values.UntypedArray{values.Bytes{}}}},
		s2,
	},
	{
		frame.Frame{RequestID: 7, Body: bodies.Write{Key: fields.Key("key ke"), Value: values.UntypedArray{values.Bytes{0xff}}}},
		s2,
	},
	{
		frame.Frame{RequestID: 8, Body: bodies.Write{Key: fields.Key("key key"), Value: values.UntypedArray{values.Bytes{0xff}, values.Int(1)}}},
		s2,
	},
	{
		frame.Frame{RequestID: 9, Body: bodies.Write{Key: fields.Key("key yek"), Value: values.TypedArray{
			Elems:    []value.V{values.Bytes{0xff}, values.Bytes{0x00}},
			ElemType: fields.Bytes,
		}}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.WriteAnswer{}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.WriteAnswer{}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Delete("")},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Delete("key")},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Delete("")},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Delete("key")},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.DeleteAnswer{}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.DeleteAnswer{}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.Ping{}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Ping{}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.PingAnswer{}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.PingAnswer{}},
		s2,
	},
	{
		frame.Frame{RequestID: 0, Body: bodies.ErrorAnswer{Err: errs.NewErrorInternalError(nil)}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorBatchLimitIsExceeded(1, 0)}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorBatchLimitIsExceeded(0, 1)}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorUnexpectedCommand(fields.Command(10), fields.Handshake)}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: errs.NewErrorUnexpectedCommand(fields.Handshake, fields.Command(10))}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: false, InterruptAfterError: false, Requests: nil}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: false, InterruptAfterError: false, Requests: nil}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: true, InterruptAfterError: false, Requests: nil}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: false, InterruptAfterError: true, Requests: nil}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: false, Requests: nil}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: false, IsOneAnswer: true, InterruptAfterError: true, Requests: nil}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: false, InterruptAfterError: true, Requests: nil}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: nil}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
		}}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
			{Number: 0, Body: bodies.Read("key")},
		}}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
			{Number: 1, Body: bodies.Read("key")},
		}}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.Batch{IsSequentialExecution: true, IsOneAnswer: true, InterruptAfterError: true, Requests: []bodies.Request{
			{Number: 0, Body: bodies.Read("key")},
			{Number: 1, Body: bodies.Delete("key")},
		}}},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer(nil)},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
		})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 0, Body: bodies.WriteAnswer{}},
		})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 0, Body: bodies.DeleteAnswer{}},
		})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 1, Body: bodies.DeleteAnswer{}},
		})},
		s2,
	},
	{
		frame.Frame{RequestID: 1, Body: bodies.BatchAnswer([]bodies.Result{
			{Number: 0, Body: bodies.WriteAnswer{}},
			{Number: 1, Body: bodies.ReadAnswer{Value: values.Int(123)}},
		})},
		s2,
	},
}

func fakeHash() [32]byte {
	hash := make([]byte, 0, 32)

	for range 32 {
		hash = append(hash, byte(rand.Int31()))
	}

	return [32]byte(hash)
}

func compareValues(t *testing.T, v1, v2 value.V) {
	switch want := v2.(type) {
	case values.Bytes:
		got, ok := v1.(values.Bytes)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", v1, want)
			return
		}
		if !bytes.Equal(want, got) {
			t.Fatalf("Compare values.Bytes: got %v, want %v", got, want)
		}
	case values.JSON:
		got, ok := v1.(values.JSON)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", v1, want)
			return
		}
		if !bytes.Equal(want, got) {
			t.Fatalf("Compare values.JSON: got %v, want %v", got, want)
		}
	case values.Int:
		got, ok := v1.(values.Int)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", v1, want)
			return
		}
		if want != got {
			t.Fatalf("Compare values.Int: got %v, want %v", got, want)
		}
	case values.Uint:
		got, ok := v1.(values.Uint)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", v1, want)
			return
		}
		if want != got {
			t.Fatalf("Compare values.Uint: got %v, want %v", got, want)
		}
	case values.Float:
		got, ok := v1.(values.Float)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", v1, want)
			return
		}
		if want != got {
			t.Fatalf("Compare values.Float: got %v, want %v", got, want)
		}
	case values.String:
		got, ok := v1.(values.String)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", v1, want)
			return
		}
		if !strings.EqualFold(string(got), string(want)) {
			t.Fatalf("Compare values.String: got %v, want %v", got, want)
		}
	case values.TypedArray:
		got, ok := v1.(values.TypedArray)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", v1, want)
			return
		}
		if want.ElemType != got.ElemType {
			t.Fatalf("Compare values.TypedArray: mismatched elem type: got %v, want %v", got.ElemType, want.ElemType)
		}
		if len(want.Elems) != len(got.Elems) {
			t.Fatalf("Compare values.TypedArray: mismatched lens: got %v, want %v", len(got.Elems), len(want.Elems))
		}

		for i := range want.Elems {
			compareValues(t, got.Elems[i], want.Elems[i])
		}
	case values.UntypedArray:
		got, ok := v1.(values.UntypedArray)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", v1, want)
			return
		}
		if len(want) != len(got) {
			t.Fatalf("Compare values.TypedArray: mismatched lens: got %v, want %v", len(got), len(want))
		}

		for i := range want {
			compareValues(t, got[i], want[i])
		}
	}
}

func compareBodies(t *testing.T, b1, b2 body.Body) {

	switch want := b2.(type) {
	case bodies.Handshake:
		got, ok := b1.(bodies.Handshake)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
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
		got, ok := b1.(bodies.HandshakeAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}

		if !slices.Equal(got.Compressions, want.Compressions) {
			t.Errorf("Compare Compressions: got %v and want %v", got.Compressions, want.Compressions)
		}
	case bodies.ErrorAnswer:
		got, ok := b1.(bodies.ErrorAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}

		if want.Err.TracebackID() != got.Err.TracebackID() {
			t.Errorf("Compare Traceback ID: got % X and want % X", got.Err.TracebackID(), want.Err.TracebackID())
		}
	case bodies.Read:
		got, ok := b1.(bodies.Read)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}

		if !bytes.Equal([]byte(got), []byte(want)) {
			t.Errorf("Compare Keys: got % X and want % X", got, want)
		}

	case bodies.ReadAnswer:
		got, ok := b1.(bodies.ReadAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}

		compareValues(t, got.Value, want.Value)

	case bodies.Write:
		got, ok := b1.(bodies.Write)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}

		if !bytes.Equal([]byte(got.Key), []byte(want.Key)) {
			t.Errorf("Compare Keys: got % X and want % X", got.Key, want.Key)
		}

		compareValues(t, got.Value, want.Value)
	case bodies.WriteAnswer:
		_, ok := b1.(bodies.WriteAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}
	case bodies.Ping:
		_, ok := b1.(bodies.Ping)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}
	case bodies.PingAnswer:
		_, ok := b1.(bodies.PingAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}
	case bodies.Delete:
		got, ok := b1.(bodies.Delete)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}

		if !bytes.Equal([]byte(got), []byte(want)) {
			t.Errorf("Compare Keys: got % X and want % X", got, want)
		}
	case bodies.DeleteAnswer:
		_, ok := b1.(bodies.DeleteAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}

	case bodies.Batch:
		got, ok := b1.(bodies.Batch)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}

		if got.InterruptAfterError != want.InterruptAfterError {
			t.Errorf("Compare InterruptAfterError: got %v and want %v", got.InterruptAfterError, want.InterruptAfterError)
		}

		if got.IsSequentialExecution != want.IsSequentialExecution {
			t.Errorf("Compare IsSequentialExecution: got %v and want %v", got.IsSequentialExecution, want.IsSequentialExecution)
		}

		if got.IsOneAnswer != want.IsOneAnswer {
			t.Errorf("Compare IsOneAnswer: got %v and want %v", got.IsOneAnswer, want.IsOneAnswer)
		}

		if len(got.Requests) != len(want.Requests) {
			t.Fatalf("Compare Batch: mismatched lens: got %v, want %v", len(got.Requests), len(want.Requests))
		}

		for i := range got.Requests {
			if got.Requests[i].Number != want.Requests[i].Number {
				t.Errorf("Compare Batch Request Numbers: got %v, want %v", got.Requests[i].Number, want.Requests[i].Number)
			}

			compareBodies(t, got.Requests[i].Body, want.Requests[i].Body)
		}
	case bodies.BatchAnswer:
		got, ok := b1.(bodies.BatchAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %v and want %v", b1, want)
			return
		}

		if len(got) != len(want) {
			t.Fatalf("Compare BatchAnswer: mismatched lens: got %v, want %v", len(got), len(want))
		}

		for i := range got {
			if got[i].Number != want[i].Number {
				t.Errorf("Compare Result Numbers: got %v, want %v", got[i].Number, want[i].Number)
			}

			compareBodies(t, got[i].Body, want[i].Body)
		}
	}
}

func compareFrames(t *testing.T, f1, f2 frame.Frame) {
	if f1.RequestID != f2.RequestID {
		t.Fatalf("Compare Request IDs: got %v, want %v", f1, f2.RequestID)
	}

	compareBodies(t, f1.Body, f2.Body)
}

func TestDecoder_BySpecs(t *testing.T) {

	for i, tt := range testsSpecs {
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

				compareFrames(t, f, tt.decoded)

			},
		)
	}
}

func TestDecoder_Preamble(t *testing.T) {
	tests := []struct {
		p         []byte
		wantError bool
		version   fields.Version
	}{
		{[]byte{}, true, 0},
		{[]byte{0x00}, true, 0},
		{[]byte{0x0A}, true, 0},
		{[]byte{0x0A, 0xDB}, true, 0},
		{[]byte{0xDB, 0x0A}, true, 0},
		{[]byte{0x11, 0xFF}, true, 0},
		{[]byte{0x0A, 0xDB, 0x00}, false, 0},
		{[]byte{0xDB, 0x0A, 0x00}, true, 0},
		{[]byte{0x11, 0xFF, 0x00}, true, 0},
		{[]byte{0x0A, 0xDB, 0x01}, false, 1},
		{[]byte{0xDB, 0x0A, 0x01}, true, 0},
		{[]byte{0x11, 0xFF, 0x01}, true, 0},
		{[]byte{0x0A, 0xDB, 0xFF}, false, 255},
		{[]byte{0xDB, 0x0A, 0xFF}, true, 0},
		{[]byte{0x11, 0xFF, 0xFF}, true, 0},
		{[]byte{0x0A, 0xDB, 0xFF, 0x01}, false, 255},
		{[]byte{0xDB, 0x0A, 0xFF, 0x01}, true, 0},
		{[]byte{0x11, 0xFF, 0xFF, 0x01}, true, 0},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDecoder_Preamble_%v", i),
			func(t *testing.T) {
				decoder := NewDecoder(1024, nil)
				reader := bufio.NewReader(bytes.NewReader(tt.p))
				got, err := decoder.DecodePreamble(reader)

				if err == nil == tt.wantError {
					t.Fatalf("DecodePreamble: expected err %v got %v", tt.wantError, err)
					return
				}

				if got != tt.version {
					t.Errorf("DecodePreamble: got %v version got %v version", got, tt.version)
				}
			},
		)
	}
}

func TestDecoder_EncodeDecodeEncode(t *testing.T) {

	decoder := NewDecoder(1024, []compressor.Compressor{
		zstd,
		s2,
	})

	for i, tt := range testsRoundTrip {
		t.Run(
			fmt.Sprintf("TestDecoder_EncodeDecodeEncode_%v", i),
			func(t *testing.T) {
				buf := buffer.Slice{}
				buf.Preallocate(tt.f.Size())

				err := tt.f.Encode(&buf, tt.compressor)
				if err != nil {
					t.Fatalf("EncodeFrame: got err %v", err)
				}

				reader := bufio.NewReader(bytes.NewReader(buf.Bytes()))

				_, err = decoder.DecodePreamble(reader)
				if err != nil {
					t.Fatalf("DecodePreamble: got err %v", err)
					return
				}

				f1, err := decoder.DecodeFrame(reader)
				if err != nil {
					t.Fatalf("DecodeFrame: got err %v", err)
					return
				}

				compareFrames(t, f1, tt.f)

				buf.Clean()
				buf.Preallocate(f1.Size())

				err = f1.Encode(&buf, tt.compressor)
				if err != nil {
					t.Fatalf("EncodeFrame: got err %v", err)
				}

				reader = bufio.NewReader(bytes.NewReader(buf.Bytes()))

				_, err = decoder.DecodePreamble(reader)
				if err != nil {
					t.Fatalf("DecodePreamble: got err %v", err)
					return
				}

				f2, err := decoder.DecodeFrame(reader)
				if err != nil {
					t.Fatalf("DecodeFrame: got err %v", err)
					return
				}

				compareFrames(t, f2, tt.f)
			},
		)
	}
}

func TestDecoder_WithUnexpectedPayload(t *testing.T) {
	tests := []struct {
		e []byte
	}{
		{
			[]byte{
				0x0A, 0xDB, 0x01, 0x06, 0x00, 0x00, 0x00, 0x04,
				0x00, 0x00, 0x00, 0x00, 0x03, 0xC0, 0xFF, 0xEE,
				0x83, 0x2C, 0x45, 0x52,
			},
		},
		{
			[]byte{
				0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x15, 0x01, 0x02, 0x06,
				0x00, 0x00, 0x00, 0x09, 0x73, 0x6F, 0x6D, 0x65,
				0x2D, 0x64, 0x61, 0x74, 0x61, 0x76, 0x69, 0x72,
				0x75, 0x73, 0xD2, 0x19, 0xB2, 0xFB,
			},
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDecoder_WithUnexpectedPayload_%v", i),
			func(t *testing.T) {
				d := NewDecoder(1<<10, nil)
				reader := bufio.NewReader(bytes.NewReader(tt.e))

				_, er := d.DecodeFrame(reader)
				if er == nil {
					t.Fatal("expected err, got nil")
				}

				if _, ok := er.(err.TrailledError); !ok {
					t.Fatalf("expected TrailledError, got %v", er)
				}
			},
		)
	}
}

func TestDecoder_WithWrongChecksum(t *testing.T) {
	tests := []struct {
		e []byte
	}{
		{
			[]byte{
				0x0A, 0xDB, 0x01, 0x06, 0x00, 0x00, 0x00, 0x04,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x83, 0x2C, 0x45,
				0x52,
			},
		},
		{
			[]byte{
				0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x26, 0x01, 0x05, 0x00,
				0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00,
				0x00, 0x00, 0x0E, 0x01, 0x02, 0x06, 0x00, 0x00,
				0x00, 0x07, 0x6D, 0x65, 0x73, 0x73, 0x61, 0x67,
				0x65, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
				0x02, 0x01, 0x03, 0xAA, 0xAA, 0xAA, 0xAA,
			},
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDecoder_WithWrongChecksum_%v", i),
			func(t *testing.T) {
				d := NewDecoder(1<<10, nil)
				reader := bufio.NewReader(bytes.NewReader(tt.e))

				_, er := d.DecodeFrame(reader)
				if er == nil {
					t.Fatal("expected err, got nil")
				}

				if _, ok := er.(errs.ErrorMismatchedChecksum); !ok {
					t.Fatalf("expected ErrorMismatchedChecksum, got %v", er)
				}
			},
		)
	}
}

func TestDecoder_WithReaderUnexpectedEOF(t *testing.T) {
	tests := []struct {
		e []byte
	}{
		{
			[]byte{
				0x0A, 0xDB, 0x01, 0x01, 0x00, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x00, 0x00, 0x02, 0x01, 0x7B, 0x54,
				0x40, 0xCA,
			},
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestDecoder_WithReaderUnexpectedEOF_%v", i),
			func(t *testing.T) {
				d := NewDecoder(1<<10, nil)
				reader := bufio.NewReader(bytes.NewReader(tt.e))

				_, er := d.DecodeFrame(reader)
				if er == nil {
					t.Fatal("expected err, got nil")
				}

				de, ok := er.(err.DecodeError)
				if !ok {
					t.Fatalf("expected ErrorMismatchedChecksum, got %v", er)
				}

				if !errors.Is(io.ErrUnexpectedEOF, de.Unwrap()) {
					t.Errorf("expected io.ErrUnexpectedEOF, got %v", de)
				}
			},
		)
	}
}

func TestDecoder_ZipBomb(t *testing.T) {
	frames := []struct {
		f          frame.Frame
		compressor compressor.Compressor
	}{
		{
			frame.Frame{RequestID: 1, Body: bodies.ReadAnswer{Value: values.Bytes(make([]byte, 1<<13))}},
			zstd,
		},
		{
			frame.Frame{RequestID: 1, Body: bodies.ReadAnswer{Value: values.Bytes(make([]byte, 1<<13))}},
			s2,
		},
	}

	for i, tt := range frames {
		t.Run(
			fmt.Sprintf("TestDecoder_ZipBomb_%v", i),
			func(t *testing.T) {
				decoder := NewDecoder(1<<12, []compressor.Compressor{zstd, s2})

				buf := buffer.Slice{}
				buf.Preallocate(tt.f.Size())

				err := tt.f.Encode(&buf, zstd)
				if err != nil {
					t.Fatalf("EncodeFrame: got err %v", err)
				}

				reader := bufio.NewReader(bytes.NewReader(buf.Bytes()))

				_, err = decoder.DecodeFrame(reader)
				if err == nil {
					t.Fatal("DecodeFrame: didn't get err")
					return
				}
			},
		)
	}
}

func TestDecoder_WithNilReader(t *testing.T) {
	t.Run(
		"TestDecoder_WithNilReader",
		func(t *testing.T) {
			d := NewDecoder(1, nil)

			_, e := d.DecodePreamble(nil)
			if e == nil {
				t.Error("DecodePreamble: expected err, got nil")
			}

			if _, ok := e.(err.DecodeError); !ok {
				t.Errorf("DecodePreamble: expected DecodeError, got %v", e)
			}

			_, e = d.DecodeFrame(nil)
			if e == nil {
				t.Error("DecodeFrame: expected err, got nil")
			}

			if _, ok := e.(err.DecodeError); !ok {
				t.Errorf("DecodeFrame: expected DecodeError, got %v", e)
			}
		},
	)
}

func FuzzDecoder_Preamble(f *testing.F) {
	f.Add([]byte{0x0A, 0xDB, 0x01})
	f.Fuzz(func(t *testing.T, a []byte) {
		decoder := NewDecoder(1024, nil)
		_, _ = decoder.DecodePreamble(bufio.NewReader(bytes.NewBuffer(a)))
	})
}

func FuzzDecoder_BySpecs(f *testing.F) {
	for _, t := range testsSpecs {
		f.Add(t.encoded)
	}
	f.Fuzz(func(t *testing.T, a []byte) {
		decoder := NewDecoder(1024, nil)
		_, _ = decoder.DecodeFrame(bufio.NewReader(bytes.NewBuffer(a)))
	})
}

func FuzzDecoder_RoundTrip(f *testing.F) {
	for _, t := range testsRoundTrip {
		buf := buffer.Slice{}
		buf.Preallocate(t.f.Size())

		_ = t.f.Encode(&buf, t.compressor)

		f.Add(buf.Bytes())
	}
	f.Fuzz(func(t *testing.T, a []byte) {
		decoder := NewDecoder(1<<12, nil)
		_, _ = decoder.DecodeFrame(bufio.NewReader(bytes.NewBuffer(a)))
	})
}
