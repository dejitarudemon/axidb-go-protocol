package decoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/dejitarudemon/axidb-go-protocol/v0/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v0/err"
	"github.com/dejitarudemon/axidb-go-protocol/v0/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v0/frame"
)

const (
	// MagicBytesLen is the size in bytes of the frame magic prefix.
	MagicBytesLen = 2
	// PreambleLen is the size in bytes of the magic prefix and version.
	PreambleLen = MagicBytesLen + fields.VersionFieldSize
)

// Decoder reads protocol v0 Hello frames.
type Decoder struct{}

// NewDecoder returns a Decoder for Hello frames.
func NewDecoder() Decoder {
	return Decoder{}
}

// handleReaderError wraps a reader failure as a decode error.
func (d Decoder) handleReaderError(e error) error {
	switch {
	case errors.Is(e, io.EOF):
		return err.NewDecodeError("reader eof", e)
	case errors.Is(e, io.ErrUnexpectedEOF):
		return err.NewDecodeError("unexpected reader eof", e)
	}

	return err.NewDecodeError("internal reader error", e)
}

// DecodePreamble peeks the magic bytes and version without consuming them.
// It returns an error when reader is nil, the read fails, or the magic bytes do not match.
func (d Decoder) DecodePreamble(reader *bufio.Reader) (fields.Version, error) {
	if reader == nil {
		return 0, err.NewDecodeError("got nil reader", nil)
	}

	buf, e := reader.Peek(PreambleLen)
	if e != nil {
		return 0, d.handleReaderError(e)
	}

	if !bytes.Equal(frame.MagicBytes, buf[:MagicBytesLen]) {
		return 0, err.NewBrokenFrameError(buf[:MagicBytesLen])
	}

	return fields.Version(buf[MagicBytesLen]), nil
}

// DecodeFrame reads one Hello frame from reader.
// It checks the magic bytes, Hello version, and checksum.
// Extra bytes after the frame are left in reader for the next message.
// reader must be non-nil.
func (d Decoder) DecodeFrame(reader *bufio.Reader) (frame.Frame, error) {
	if reader == nil {
		return frame.Frame{}, err.NewDecodeError("got nil reader", nil)
	}

	var headersBuf [PreambleLen + frame.HeadersLen]byte
	if _, e := io.ReadFull(reader, headersBuf[:]); e != nil {
		return frame.Frame{}, d.handleReaderError(e)
	}

	if !bytes.Equal(frame.MagicBytes, headersBuf[:MagicBytesLen]) {
		return frame.Frame{}, err.NewBrokenFrameError(headersBuf[:MagicBytesLen])
	}

	version := fields.Version(headersBuf[MagicBytesLen])
	if version != fields.CurrentVersion {
		return frame.Frame{}, err.NewDecodeError("unexpected protocol version", nil)
	}

	versionLen := int(headersBuf[PreambleLen])

	rest := make([]byte, versionLen+fields.ChecksumFieldSize)
	if _, e := io.ReadFull(reader, rest); e != nil {
		return frame.Frame{}, d.handleReaderError(e)
	}

	versionsBuf, checksumBuf := rest[:versionLen], rest[versionLen:]

	checksumWanted := fields.Checksum(binary.BigEndian.Uint32(checksumBuf))
	checksumReal := fields.NewChecksumWithParts(headersBuf[:], versionsBuf)

	if !checksumReal.Equal(checksumWanted) {
		return frame.Frame{}, err.NewMismatchedChecksum(checksumReal, checksumWanted)
	}

	versions := make([]fields.Version, versionLen)
	for i, b := range versionsBuf {
		versions[i] = fields.Version(b)
	}

	return frame.Frame{
		Body: bodies.Hello{Versions: versions},
	}, nil
}
