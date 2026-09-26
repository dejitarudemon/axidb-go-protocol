package decoder

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"testing"
	"testing/iotest"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
)

// rawFrame builds a frame with a valid checksum around an arbitrary body.
func rawFrame(command fields.Command, compression fields.Compression, body []byte) []byte {
	headers := cat(
		frame.MagicBytes,
		[]byte{0x01, byte(command)},
		u32(1),
		[]byte{byte(compression)},
		u32(uint32(len(body))),
	)

	return cat(headers, body, u32(uint32(fields.NewChecksumWithParts(headers, body))))
}

type fakeCompressor struct {
	code       fields.Compression
	decompress func([]byte) ([]byte, error)
}

func (f fakeCompressor) Code() fields.Compression { return f.code }

func (f fakeCompressor) Compress(buf []byte) ([]byte, error) { return buf, nil }

func (f fakeCompressor) Decompress(buf []byte) ([]byte, error) { return f.decompress(buf) }

func decodeFrameBytes(d Decoder, data []byte) (frame.Frame, error) {
	return d.DecodeFrame(bufio.NewReader(bytes.NewReader(data)))
}

func TestDecoder_DecodeFrame_ReaderError(t *testing.T) {
	sentinel := errors.New("broken pipe")
	d := NewDecoder(1024, nil)

	_, e := d.DecodeFrame(bufio.NewReader(iotest.ErrReader(sentinel)))

	var de err.DecodeError
	if !errors.As(e, &de) {
		t.Fatalf("expected DecodeError, got %v", e)
	}

	if !errors.Is(e, sentinel) {
		t.Errorf("expected wrapped %v, got %v", sentinel, e)
	}
}

func TestDecoder_DecodePreamble_ReaderError(t *testing.T) {
	sentinel := errors.New("broken pipe")
	d := NewDecoder(1024, nil)

	_, e := d.DecodePreamble(bufio.NewReader(iotest.ErrReader(sentinel)))
	if !errors.Is(e, sentinel) {
		t.Errorf("expected wrapped %v, got %v", sentinel, e)
	}
}

func TestDecoder_DecodeFrame_TruncatedChecksum(t *testing.T) {
	d := NewDecoder(1024, nil)
	data := rawFrame(fields.Ping, fields.None, nil)

	_, e := decodeFrameBytes(d, data[:len(data)-1])
	if !errors.Is(e, io.ErrUnexpectedEOF) {
		t.Fatalf("expected unexpected EOF, got %v", e)
	}
}

func TestDecoder_DecodeFrame_UnsupportedCompression(t *testing.T) {
	d := NewDecoder(1024, nil)

	_, e := decodeFrameBytes(d, rawFrame(fields.Ping, fields.Zstd, nil))

	var want errs.ErrorUnsupportedCompression
	if !errors.As(e, &want) {
		t.Fatalf("expected ErrorUnsupportedCompression, got %v", e)
	}
}

func TestDecoder_DecodeFrame_CompressorError(t *testing.T) {
	sentinel := errors.New("corrupted stream")
	d := NewDecoder(1024, []compressor.Compressor{
		fakeCompressor{code: fields.Zstd, decompress: func([]byte) ([]byte, error) { return nil, sentinel }},
	})

	_, e := decodeFrameBytes(d, rawFrame(fields.Ping, fields.Zstd, []byte{0x01}))

	var de err.DecodeError
	if !errors.As(e, &de) || !errors.Is(e, sentinel) {
		t.Fatalf("expected DecodeError wrapping %v, got %v", sentinel, e)
	}
}

func TestDecoder_DecodeFrame_TrailledAfterDecompression(t *testing.T) {
	d := NewDecoder(1024, []compressor.Compressor{
		fakeCompressor{code: fields.S2, decompress: func([]byte) ([]byte, error) { return []byte{0x00, 0x01, 0x02}, nil }},
	})

	_, e := decodeFrameBytes(d, rawFrame(fields.Ping, fields.S2, []byte{0xFF}))

	want := err.NewTrailledError(3, 0)
	if !errors.Is(e, want) {
		t.Fatalf("expected %v, got %v", want, e)
	}
}

func TestDecoder_DecodeFrame_BodyErrorKeepsRequestID(t *testing.T) {
	d := NewDecoder(1024, nil)

	f, e := decodeFrameBytes(d, rawFrame(fields.Write, fields.None, u32(10)))
	assertMalformed(t, e)

	if f.RequestID != 1 {
		t.Errorf("RequestID: got %v, want 1", f.RequestID)
	}

	if f.Body != nil {
		t.Errorf("Body: got %v, want nil", f.Body)
	}
}

func TestDecoder_DecodeFrame_UnsupportedCommand(t *testing.T) {
	d := NewDecoder(1024, nil)

	_, e := decodeFrameBytes(d, rawFrame(fields.Command(0x7F), fields.None, nil))

	var want errs.ErrorUnsupportedCommand
	if !errors.As(e, &want) {
		t.Fatalf("expected ErrorUnsupportedCommand, got %v", e)
	}
}

func TestDecoder_DecodeFrame_MalformedValueAnswer(t *testing.T) {
	d := NewDecoder(1024, nil)
	pe := errs.NewErrorMalformedValueWithTracebackID("bad value", testTracebackID)

	f, e := decodeFrameBytes(d, rawFrame(fields.Answer, fields.None, encodeBody(bodies.ErrorAnswer{Err: pe})))
	if e != nil {
		t.Fatalf("DecodeFrame: got err %v", e)
	}

	compareFrames(t, f, frame.Frame{RequestID: 1, Body: bodies.ErrorAnswer{Err: pe}})
}
