package decoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
	"testing/iotest"

	"github.com/dejitarudemon/axidb-go-protocol/v0/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v0/err"
	"github.com/dejitarudemon/axidb-go-protocol/v0/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v0/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v0/internal/specs"
	"github.com/dejitarudemon/axidb-go-protocol/v0/internal/testutil"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

var roundTripFrames = []struct {
	name string
	f    frame.Frame
}{
	{"empty", frame.Frame{Body: bodies.NewHello(nil)}},
	{"one version", frame.Frame{Body: bodies.NewHello([]fields.Version{1})}},
	{"client spec", frame.Frame{Body: bodies.NewHello([]fields.Version{1, 2, 3})}},
	{"server spec", frame.Frame{Body: bodies.NewHello([]fields.Version{1, 4, 7, 11})}},
	{"duplicates on wire", frame.Frame{Body: bodies.Hello{Versions: []fields.Version{1, 1, 2}}}},
	{"max versions", frame.Frame{Body: bodies.NewHello(allWorkingVersions())}},
}

func allWorkingVersions() []fields.Version {
	vs := make([]fields.Version, bodies.MaxVersionsPerOneHello)
	for i := range vs {
		vs[i] = fields.Version(i + 1)
	}
	return vs
}

func encodeFrame(t testing.TB, f frame.Frame) []byte {
	t.Helper()

	buf := buffer.Slice{}
	buf.Preallocate(f.Size())

	if e := f.Encode(&buf); e != nil {
		t.Fatalf("Encode() = %v", e)
	}

	return buf.Bytes()
}

func decodeFrameBytes(d Decoder, data []byte) (frame.Frame, error) {
	return d.DecodeFrame(bufio.NewReader(bytes.NewReader(data)))
}

func assertSameFrame(t *testing.T, got, want frame.Frame) {
	t.Helper()

	testutil.AssertSameEncoding(t, got.Body, want.Body)
}

func rawHello(versions []byte) []byte {
	headers := []byte{0x0A, 0xDB, 0x00, byte(len(versions))}
	sum := fields.NewChecksumWithParts(headers, versions)
	return append(append(headers, versions...), binary.BigEndian.AppendUint32(nil, uint32(sum))...)
}

func isA[T error]() func(error) bool {
	return func(e error) bool {
		var target T
		return errors.As(e, &target)
	}
}

func isEOF(target error) func(error) bool {
	return func(e error) bool {
		return isA[err.DecodeError]()(e) && errors.Is(e, target)
	}
}

func TestDecoder_Specs(t *testing.T) {
	for _, tt := range specs.Frames {
		t.Run(tt.Name, func(t *testing.T) {
			d := NewDecoder()
			reader := bufio.NewReader(bytes.NewReader(tt.Encoded))

			version, e := d.DecodePreamble(reader)
			if e != nil {
				t.Fatalf("DecodePreamble() = %v", e)
			}

			if version != fields.CurrentVersion {
				t.Errorf("DecodePreamble() = %v, want %v", version, fields.CurrentVersion)
			}

			got, e := d.DecodeFrame(reader)
			if e != nil {
				t.Fatalf("DecodeFrame() = %v", e)
			}

			assertSameFrame(t, got, tt.Frame)
		})
	}
}

func TestDecoder_RoundTrip(t *testing.T) {
	d := NewDecoder()

	for _, tt := range roundTripFrames {
		t.Run(tt.name, func(t *testing.T) {
			got, e := decodeFrameBytes(d, encodeFrame(t, tt.f))
			if e != nil {
				t.Fatalf("DecodeFrame() = %v", e)
			}

			assertSameFrame(t, got, tt.f)
		})
	}
}

func TestDecoder_DecodePreamble(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    fields.Version
		wantErr bool
	}{
		{"empty", []byte{}, 0, true},
		{"one byte", []byte{0x0A}, 0, true},
		{"magic without version", []byte{0x0A, 0xDB}, 0, true},
		{"swapped magic", []byte{0xDB, 0x0A, 0x00}, 0, true},
		{"wrong magic", []byte{0x11, 0xFF, 0x00}, 0, true},
		{"version 0", []byte{0x0A, 0xDB, 0x00}, 0, false},
		{"version 1", []byte{0x0A, 0xDB, 0x01}, 1, false},
		{"version 255", []byte{0x0A, 0xDB, 0xFF}, 255, false},
		{"followed by data", []byte{0x0A, 0xDB, 0x00, 0x01}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e := NewDecoder().DecodePreamble(bufio.NewReader(bytes.NewReader(tt.data)))
			testutil.AssertErr(t, e, tt.wantErr)

			if got != tt.want {
				t.Errorf("DecodePreamble() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_DecodePreamble_DoesNotConsume(t *testing.T) {
	d := NewDecoder()
	data := specs.Frames[0].Encoded
	reader := bufio.NewReader(bytes.NewReader(data))

	if _, e := d.DecodePreamble(reader); e != nil {
		t.Fatalf("DecodePreamble() = %v", e)
	}

	if _, e := d.DecodeFrame(reader); e != nil {
		t.Errorf("DecodeFrame() after DecodePreamble() = %v", e)
	}
}

func TestDecoder_DecodeFrameErrs(t *testing.T) {
	valid := rawHello([]byte{0x01, 0x02, 0x03})
	corrupted := bytes.Clone(valid)
	corrupted[len(corrupted)-1] ^= 0xFF

	tests := []struct {
		name  string
		data  []byte
		match func(error) bool
	}{
		{"empty reader", []byte{}, isEOF(io.EOF)},
		{"truncated headers", valid[:2], isEOF(io.ErrUnexpectedEOF)},
		{"truncated versions", valid[:6], isEOF(io.ErrUnexpectedEOF)},
		{"truncated checksum", valid[:len(valid)-1], isEOF(io.ErrUnexpectedEOF)},
		{"wrong magic", append([]byte{0x11, 0xFF, 0x00, 0x00}, make([]byte, 4)...), isA[err.BrokenFrameError]()},
		{"unexpected version", rawVersion1(), isA[err.DecodeError]()},
		{"mismatched checksum", corrupted, isA[err.MismatchedChecksum]()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := decodeFrameBytes(NewDecoder(), tt.data)
			if !tt.match(e) {
				t.Errorf("DecodeFrame() = %v", e)
			}
		})
	}
}

func rawVersion1() []byte {
	headers := []byte{0x0A, 0xDB, 0x01, 0x00}
	sum := fields.NewChecksum(headers)
	return append(headers, binary.BigEndian.AppendUint32(nil, uint32(sum))...)
}

func TestDecoder_LeavesTrailingBytes(t *testing.T) {
	hello := rawHello([]byte{0x01})
	extra := []byte{0x0A, 0xDB, 0x01}
	reader := bufio.NewReader(bytes.NewReader(append(hello, extra...)))

	if _, e := NewDecoder().DecodeFrame(reader); e != nil {
		t.Fatalf("DecodeFrame() = %v", e)
	}

	leftover := make([]byte, len(extra))
	if _, e := io.ReadFull(reader, leftover); e != nil {
		t.Fatalf("leftover read = %v", e)
	}

	if !bytes.Equal(leftover, extra) {
		t.Errorf("leftover = % X, want % X", leftover, extra)
	}
}

func TestDecoder_ReaderErrs(t *testing.T) {
	sentinel := errors.New("broken pipe")

	tests := []struct {
		name   string
		decode func(d Decoder, r *bufio.Reader) error
	}{
		{"DecodePreamble", func(d Decoder, r *bufio.Reader) error { _, e := d.DecodePreamble(r); return e }},
		{"DecodeFrame", func(d Decoder, r *bufio.Reader) error { _, e := d.DecodeFrame(r); return e }},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/nil reader", func(t *testing.T) {
			if e := tt.decode(NewDecoder(), nil); !isA[err.DecodeError]()(e) {
				t.Errorf("got %v, want DecodeError", e)
			}
		})

		t.Run(tt.name+"/failing reader", func(t *testing.T) {
			e := tt.decode(NewDecoder(), bufio.NewReader(iotest.ErrReader(sentinel)))
			if !isA[err.DecodeError]()(e) || !errors.Is(e, sentinel) {
				t.Errorf("got %v, want DecodeError wrapping %v", e, sentinel)
			}
		})
	}
}

func FuzzDecoder_Preamble(f *testing.F) {
	f.Add([]byte{0x0A, 0xDB, 0x00})
	f.Fuzz(func(t *testing.T, a []byte) {
		_, _ = NewDecoder().DecodePreamble(bufio.NewReader(bytes.NewBuffer(a)))
	})
}

func FuzzDecoder_BySpecs(f *testing.F) {
	for _, tt := range specs.Frames {
		f.Add(tt.Encoded)
	}
	f.Fuzz(func(t *testing.T, a []byte) {
		_, _ = NewDecoder().DecodeFrame(bufio.NewReader(bytes.NewBuffer(a)))
	})
}

func FuzzDecoder_RoundTrip(f *testing.F) {
	for _, tt := range roundTripFrames {
		f.Add(encodeFrame(f, tt.f))
	}
	f.Fuzz(func(t *testing.T, a []byte) {
		_, _ = NewDecoder().DecodeFrame(bufio.NewReader(bytes.NewBuffer(a)))
	})
}
