package frame

import (
	"github.com/dejitarudemon/axidb-go-protocol/v0/body"
	"github.com/dejitarudemon/axidb-go-protocol/v0/err"
	"github.com/dejitarudemon/axidb-go-protocol/v0/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// MagicBytes is the two-byte frame preamble that identifies an AxiDB protocol frame.
var MagicBytes = []byte{0x0A, 0xDB}

const (
	// HeadersLen is the size in bytes of the fixed header after the magic bytes and version:
	// the version count.
	HeadersLen = fields.VersionLenFieldSize
)

// Frame is a protocol v0 Hello frame.
type Frame struct {
	// Body holds the advertised versions; it must be non-nil to encode or validate.
	Body body.Body
}

// Size returns the encoded frame size in bytes, including magic, headers, body, and checksum.
func (f Frame) Size() int {
	size := len(MagicBytes) + f.Version().Size() + fields.VersionLenFieldSize + fields.ChecksumFieldSize
	if f.Body == nil {
		return size
	}

	return size + f.Body.Size()
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

// Encode writes the wire encoding of the frame into buf.
// Body must be non-nil.
func (f Frame) Encode(buf buffer.Buffer) error {
	if f.Body == nil {
		return err.NewValidationError(
			"nil Body",
			"target", "Frame",
		)
	}

	buf.Append(MagicBytes)
	f.Version().Encode(buf)
	fields.VersionLen(f.Body.Size()).Encode(buf)
	f.Body.Encode(buf)
	fields.NewChecksum(buf.Raw()).Encode(buf)

	return nil
}
