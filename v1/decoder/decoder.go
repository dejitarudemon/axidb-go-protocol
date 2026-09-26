package decoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
)

const (
	// MagicBytesLen is the size in bytes of the frame magic prefix.
	MagicBytesLen = 2
	// PreambleLen is the size in bytes of the magic prefix and version.
	PreambleLen = MagicBytesLen + fields.VersionFieldSize

	// CommandOffset is the header index of the command byte.
	CommandOffset = 3
	// RequestIDOffset is the header index of the request ID.
	RequestIDOffset = 4
	// CompressionOffset is the header index of the compression code.
	CompressionOffset = 8
	// BodyLenOffset is the header index of the body length.
	BodyLenOffset = 9

	// HandshakeMinBodySize is the minimum decoded size in bytes of a handshake body.
	HandshakeMinBodySize = 37
	// WriteMinBodySize is the minimum decoded size in bytes of a write body.
	WriteMinBodySize = 5
	// AnswerMinBodySize is the minimum decoded size in bytes of an answer body.
	AnswerMinBodySize = 2
	// BatchMinBodySize is the minimum decoded size in bytes of a batch body.
	BatchMinBodySize = 5
	// BatchRequestMinBodySize is the minimum decoded size in bytes of one batch request.
	BatchRequestMinBodySize = 9
	// BatchAnswerMinBodySize is the minimum decoded size in bytes of a batch answer payload.
	BatchAnswerMinBodySize = 4
	// BatchResultMinBodySize is the minimum decoded size in bytes of one batch result.
	BatchResultMinBodySize = 8
)

// Decoder reads protocol v1 frames and rejects bodies larger than a configured limit.
type Decoder struct {
	limit       fields.BodyLimit
	compressors map[fields.Compression]compressor.Compressor
}

// NewDecoder returns a Decoder that rejects bodies larger than limit.
// compressors are indexed by their protocol code; the first compressor for a code is kept.
func NewDecoder(limit fields.BodyLimit, compressors []compressor.Compressor) Decoder {
	c := make(map[fields.Compression]compressor.Compressor, len(compressors))

	for _, compressor := range compressors {
		if _, ok := c[compressor.Code()]; !ok {
			c[compressor.Code()] = compressor
		}
	}

	return Decoder{
		limit:       limit,
		compressors: c,
	}
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

// DecodeFrame reads one frame from reader.
// It checks the checksum and body limit, decompresses the body when needed, and decodes the command payload.
// reader must be non-nil.
func (d Decoder) DecodeFrame(reader *bufio.Reader) (frame.Frame, error) {
	if reader == nil {
		return frame.Frame{}, err.NewDecodeError("got nil reader", nil)
	}

	headersBuf := make([]byte, PreambleLen+frame.HeadersLen)
	if _, e := io.ReadFull(reader, headersBuf); e != nil {
		return frame.Frame{}, d.handleReaderError(e)
	}

	command := fields.Command(headersBuf[CommandOffset])
	requestID := fields.RequestID(binary.BigEndian.Uint32(headersBuf[RequestIDOffset:]))
	compression := fields.Compression(headersBuf[CompressionOffset])
	bodyLen := binary.BigEndian.Uint32(headersBuf[BodyLenOffset:])

	if bodyLen > uint32(d.limit) {
		return frame.Frame{}, errs.NewErrorBodyLimitIsExceeded(bodyLen, d.limit)
	}

	compressor, ok := d.compressors[compression]
	if !ok && compression != fields.None {
		return frame.Frame{}, errs.NewErrorUnsupportedCompression(compression)
	}

	bodyBuf := make([]byte, int(bodyLen)+fields.ChecksumFieldSize)
	if _, e := io.ReadFull(reader, bodyBuf); e != nil {
		return frame.Frame{}, d.handleReaderError(e)
	}

	bodyBuf, checksumBuf := bodyBuf[:bodyLen], bodyBuf[bodyLen:]

	checksumWanted := fields.Checksum(binary.BigEndian.Uint32(checksumBuf))
	checksumReal := fields.NewChecksumWithParts(headersBuf, bodyBuf)

	if !checksumReal.Equal(checksumWanted) {
		return frame.Frame{}, errs.NewErrorMismatchedChecksum(checksumReal, checksumWanted)
	}

	f := frame.Frame{
		RequestID: requestID,
	}

	if compressor != nil {
		decompressed, e := compressor.Decompress(bodyBuf)
		if e != nil {
			return frame.Frame{}, err.NewDecodeError("internal compressor error", e)
		}

		bodyBuf = decompressed
	}

	c := newCursor(bodyBuf)

	b, e := d.decodeBody(c, command)
	if e != nil {
		return f, e
	}

	if c.remaining() != 0 {
		return f, err.NewTrailledError(uint32(len(bodyBuf)), uint32(c.pos))
	}

	f.Body = b
	return f, nil
}

// decodeBody decodes a command payload from c.
// It does not check that c is fully consumed.
func (d Decoder) decodeBody(c *cursor, command fields.Command) (body.Body, error) {
	switch command {
	case fields.Handshake:
		return d.handshake(c)
	case fields.Read:
		return d.read(c)
	case fields.Delete:
		return d.delete(c)
	case fields.Ping:
		return d.ping(c)
	case fields.Write:
		return d.write(c)
	case fields.Answer:
		return d.answer(c)
	case fields.Batch:
		return d.batch(c)
	}

	return nil, errs.NewErrorUnsupportedCommand(command)
}
