package frame

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

// Buffer is the buffer type used when compressing a frame body before encode.
type Buffer = buffer.Slice

// MagicBytes is the two-byte frame preamble that identifies an AxiDB protocol v1 frame.
var MagicBytes = []byte{0x0A, 0xDB}

const (
	// HeadersLen is the size in bytes of the fixed header after the magic bytes and version:
	// command, request ID, compression, and body length.
	HeadersLen = 10
)

// Frame is a protocol v1 message frame with a request ID and body payload.
type Frame struct {
	// RequestID identifies the request associated with the frame.
	RequestID fields.RequestID
	// Body holds the command-specific payload; it must be non-nil to encode or validate.
	Body body.Body
}

// Size returns the encoded frame size in bytes, including magic, headers, body, and checksum.
func (f Frame) Size() int {
	size := f.RequestID.Size() + body.BodyLenFieldSize + len(MagicBytes) + f.Version().Size() + fields.CompressionFieldSize + fields.ChecksumFieldSize
	if f.Body == nil {
		return fields.CommandFieldSize + size
	}

	return f.Body.Command().Size() + f.Body.Size() + size
}

// IsValid reports whether the frame is well-formed.
// Body must be non-nil and must pass [body.Body.IsValid].
func (f Frame) IsValid() error {
	if f.Body == nil {
		return err.NewValidationError(
			"nil Body",
			"target", "Frame",
		)
	}

	return f.Body.IsValid()
}

// Version returns the protocol version written when the frame is encoded.
func (f Frame) Version() fields.Version {
	return fields.CurrentVersion
}

// encodeWithCompression writes the frame body using the given compressor.
func (f Frame) encodeWithCompression(buf buffer.Buffer, compressor compressor.Compressor) error {
	temp := Buffer{}
	temp.Preallocate(f.Body.Size())

	f.Body.Encode(&temp)

	compressed, err := compressor.Compress(temp.Raw())
	if err != nil {
		return err
	}

	f.Body.Command().Encode(buf)
	f.RequestID.Encode(buf)
	compressor.Code().Encode(buf)
	buf.AppendUint32(uint32(len(compressed)))
	buf.Append(compressed)

	fields.NewChecksum(buf.Raw()).Encode(buf)

	return nil
}

// encodeWithoutCompression writes the frame body without compression.
func (f Frame) encodeWithoutCompression(buf buffer.Buffer) {
	f.Body.Command().Encode(buf)

	f.RequestID.Encode(buf)
	fields.None.Encode(buf)

	buf.AppendUint32(uint32(f.Body.Size()))
	f.Body.Encode(buf)

	fields.NewChecksum(buf.Raw()).Encode(buf)
}

// Encode writes the wire encoding of the frame into buf.
// If compressor is nil, the body is written uncompressed; otherwise it is compressed with compressor.
// Body must be non-nil.
func (f Frame) Encode(buf buffer.Buffer, compressor compressor.Compressor) error {
	if f.Body == nil {
		return err.NewValidationError(
			"nil Body",
			"target", "Frame",
		)
	}

	buf.Append(MagicBytes)
	f.Version().Encode(buf)

	if compressor == nil {
		f.encodeWithoutCompression(buf)
		return nil
	}

	return f.encodeWithCompression(buf, compressor)
}
