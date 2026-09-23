package decoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
)

const (
	MagicBytesLen = 2
	PreambleLen   = MagicBytesLen + fields.VersionFieldSize

	CommandOffset     = 3
	RequestIDOffest   = 4
	CompressionOffset = 8
	BodyLenOffset     = 9

	HandshakeLoginLenOffset = 0
	HandshakeLoginOffset    = 4
	HandshakeMinBodySize    = 37

	WriteKeyLenOffset = 0
	WriteKeyOffset    = 4
	WriteMinBodySize  = 5
)

type Decoder struct {
	limit       fields.BodyLimit
	compressors map[fields.Compression]compressor.Compressor
}

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

func (d Decoder) handleReaderError(e error) error {
	switch e {
	case io.EOF:
		return err.NewDecodeError("reader eof", e)
	case io.ErrUnexpectedEOF:
		return err.NewDecodeError("unexpected reader eof", e)
	}
	return err.NewDecodeError("internal reader error", e)
}

func (d Decoder) DecodePreamble(reader bufio.Reader) (fields.Version, error) {
	buf, e := reader.Peek(MagicBytesLen)
	if e != nil {
		return 0, d.handleReaderError(e)
	}

	if !bytes.Equal(frame.MagicBytes, buf[:MagicBytesLen]) {
		return 0, err.NewBrokenFrameError(buf[:MagicBytesLen])
	}

	return fields.Version(buf[MagicBytesLen]), nil
}

func (d Decoder) decodeUint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

func (d Decoder) decodeString(data []byte) string {
	return string(data)
}

func (d Decoder) DecodeFrame(reader bufio.Reader) (frame.Frame, error) {
	headersBuf := make([]byte, 0, PreambleLen+frame.HeadersLen)

	if _, e := io.ReadFull(&reader, headersBuf); e != nil {
		return frame.Frame{}, d.handleReaderError(e)
	}

	command := fields.Command(headersBuf[CommandOffset])
	requestID := fields.RequestID(d.decodeUint32(headersBuf[RequestIDOffest : RequestIDOffest+fields.RequestIDFieldSize]))
	compression := fields.Compression(headersBuf[CompressionOffset])
	bodyLen := d.decodeUint32(headersBuf[BodyLenOffset : BodyLenOffset+body.BodyLenFieldSize])

	if bodyLen > uint32(d.limit) {
		return frame.Frame{}, errs.NewErrorBodyLimitIsExceeded(bodyLen, d.limit)
	}

	compressor, ok := d.compressors[compression]
	if !ok {
		return frame.Frame{}, errs.NewErrorUnsupportedCompression(compression)
	}

	bodyBuf := make([]byte, 0, int(bodyLen))
	if _, e := io.ReadFull(&reader, bodyBuf); e != nil {
		return frame.Frame{}, d.handleReaderError(e)
	}

	checksumBuf := make([]byte, 0, fields.ChecksumFieldSize)
	if _, e := io.ReadFull(&reader, checksumBuf); e != nil {
		return frame.Frame{}, d.handleReaderError(e)
	}

	checksumWanted := fields.Checksum(d.decodeUint32(checksumBuf))
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

	switch command {
	case fields.Handshake:
		body, err := d.handshake(bodyBuf)
		if err != nil {
			return frame.Frame{}, err
		}

		f.Body = body
	case fields.Read:
		body, err := d.read(bodyBuf)
		if err != nil {
			return frame.Frame{}, err
		}

		f.Body = body
	case fields.Delete:
		body, err := d.delete(bodyBuf)
		if err != nil {
			return frame.Frame{}, err
		}
		f.Body = body
	case fields.Ping:
		f.Body = bodies.Ping{}
	case fields.Write:

	}

	return f, nil
}

func (d Decoder) handshake(buf []byte) (bodies.Handshake, error) {
	if len(buf) < HandshakeMinBodySize {
		return bodies.Handshake{}, errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid handshake-request: %v bytes min len, got %v bytes", HandshakeMinBodySize, len(buf)),
		)
	}

	bufLen := uint32(len(buf))

	loginLen := d.decodeUint32(buf[HandshakeLoginLenOffset : HandshakeLoginLenOffset+bodies.LoginLenFieldSize])

	if loginLen >= bufLen {
		return bodies.Handshake{}, errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid handshake-request: login (%v bytes) is greather than body (%v)", loginLen, len(buf)),
		)
	}

	login := ""
	off := uint32(HandshakeLoginOffset)

	if loginLen != 0 {
		login = d.decodeString(buf[HandshakeLoginOffset : HandshakeLoginOffset+loginLen])
	}

	off += loginLen
	if off+bodies.HashFieldSize > bufLen {
		return bodies.Handshake{}, errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid handshake-request: hash (%v bytes) missed", bodies.HashFieldSize),
		)
	}

	hash := buf[off : off+bodies.HashFieldSize]
	off += bodies.HashFieldSize + bodies.CompressionLenFieldSize

	if off > bufLen {
		return bodies.Handshake{}, errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid handshake-request: supported compression len (%v bytes) missed", bodies.CompressionLenFieldSize),
		)
	}

	compressionsLen := uint8(buf[off])

	if off+uint32(compressionsLen) > bufLen {
		return bodies.Handshake{}, errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid handshake-request: expected %v compressions", compressionsLen),
		)
	}

	compressions := make([]fields.Compression, 0, compressionsLen)

	for _, c := range buf[off : off+uint32(compressionsLen)] {
		compressions = append(compressions, fields.Compression(c))
	}

	return bodies.NewHandshake(login, [32]byte(hash), compressions), nil
}

func (d Decoder) read(buf []byte) (bodies.Read, error) {
	key := make([]byte, 0, len(buf))
	copy(key, buf)

	return bodies.Read(key), nil
}

func (d Decoder) delete(buf []byte) (bodies.Delete, error) {
	key := make([]byte, 0, len(buf))
	copy(key, buf)

	return bodies.Delete(key), nil
}

func (d Decoder) write(buf []byte) (bodies.Write, error) {
	if len(buf) < WriteMinBodySize {
		return bodies.Write{}, errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid write-request: %v bytes min len, got %v bytes", WriteMinBodySize, len(buf)),
		)
	}

	bufLen := uint32(len(buf))
	keyLen := d.decodeUint32(buf[WriteKeyLenOffset : WriteKeyLenOffset+fields.KeyLenFieldSize])

	if keyLen > bufLen {
		return bodies.Write{}, errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid write-request: key (%v bytes) is greather than body (%v)", keyLen, len(buf)),
		)
	}

	key := fields.Key{}
	off := uint32(WriteKeyOffset)

	if keyLen != 0 {
		key = fields.Key(buf[WriteKeyOffset : WriteKeyOffset+keyLen])
	}

	off += keyLen

	if off+fields.TypeFieldSize > bufLen {
		return bodies.Write{}, errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid write-request: value type (%v bytes) missed", fields.TypeFieldSize),
		)
	}

	return bodies.Write{Key: key}, nil
}
