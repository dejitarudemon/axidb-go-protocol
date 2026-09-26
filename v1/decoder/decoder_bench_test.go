package decoder

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/specs"
)

func BenchmarkDecoder_DecodeFrame(b *testing.B) {
	d := NewDecoder(1<<20, nil)

	for _, tt := range specs.Frames {
		b.Run(tt.Name, func(b *testing.B) {
			b.SetBytes(int64(len(tt.Encoded)))
			b.ReportAllocs()

			for b.Loop() {
				if _, e := d.DecodeFrame(bufio.NewReader(bytes.NewReader(tt.Encoded))); e != nil {
					b.Fatal(e)
				}
			}
		})
	}
}

func BenchmarkDecoder_DecodePreamble(b *testing.B) {
	d := NewDecoder(1<<20, nil)
	data := specs.Frames[0].Encoded

	b.SetBytes(3)
	b.ReportAllocs()

	for b.Loop() {
		if _, e := d.DecodePreamble(bufio.NewReader(bytes.NewReader(data))); e != nil {
			b.Fatal(e)
		}
	}
}

func BenchmarkDecoder_RoundTrip(b *testing.B) {
	d := NewDecoder(1<<12, []compressor.Compressor{zstd, s2})
	cases := []struct {
		name       string
		compressor compressor.Compressor
	}{
		{"none", nil},
		{"zstd", zstd},
		{"s2", s2},
	}

	for _, c := range cases {
		f := roundTripFrames[0].f
		encoded := encodeFrame(b, f, c.compressor)

		b.Run(c.name, func(b *testing.B) {
			b.SetBytes(int64(len(encoded)))
			b.ReportAllocs()

			for b.Loop() {
				if _, e := d.DecodeFrame(bufio.NewReader(bytes.NewReader(encoded))); e != nil {
					b.Fatal(e)
				}
			}
		})
	}
}

func BenchmarkFrame_Encode(b *testing.B) {
	for _, tt := range specs.Frames {
		b.Run(tt.Name, func(b *testing.B) {
			b.SetBytes(int64(tt.Frame.Size()))
			b.ReportAllocs()

			var buf buffer.Slice
			for b.Loop() {
				buf.Clean()
				buf.Preallocate(tt.Frame.Size())
				if e := tt.Frame.Encode(&buf, nil); e != nil {
					b.Fatal(e)
				}
			}
		})
	}
}
