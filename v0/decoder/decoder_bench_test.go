package decoder

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v0/internal/specs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func BenchmarkDecoder_DecodeFrame(b *testing.B) {
	d := NewDecoder()

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
	d := NewDecoder()
	data := specs.Frames[0].Encoded

	b.SetBytes(3)
	b.ReportAllocs()

	for b.Loop() {
		if _, e := d.DecodePreamble(bufio.NewReader(bytes.NewReader(data))); e != nil {
			b.Fatal(e)
		}
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
				if e := tt.Frame.Encode(&buf); e != nil {
					b.Fatal(e)
				}
			}
		})
	}
}
