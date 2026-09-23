package decoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
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

	BytesLenOffsetInValue = 0
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

func (d Decoder) decodeUint8(data []byte) uint8 {
	u := make([]byte, 0, 1)
	copy(u, data[:1])

	return uint8(u[0])
}

func (d Decoder) DecodeFrame(reader bufio.Reader) (frame.Frame, error) {
	headersBuf := make([]byte, 0, PreambleLen+frame.HeadersLen)

	if _, e := io.ReadFull(&reader, headersBuf); e != nil {
		return frame.Frame{}, d.handleReaderError(e)
	}

	command := fields.Command(d.decodeUint8(headersBuf[CommandOffset:]))
	requestID := fields.RequestID(d.decodeUint32(headersBuf[RequestIDOffest:]))
	compression := fields.Compression(d.decodeUint8(headersBuf[CompressionOffset:]))
	bodyLen := d.decodeUint32(headersBuf[BodyLenOffset:])

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
		return d.handshake(f, bodyBuf)
	case fields.Read:
		return d.read(f, bodyBuf)
	case fields.Delete:
		return d.delete(f, bodyBuf)
	case fields.Ping:
		return d.ping(f, bodyBuf)
	case fields.Write:

	}

	return f, nil
}

func (d Decoder) checkIfBodyLenIsTooSmall(l uint32, bound uint32) error {
	if l < bound {
		return errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid request: %v bytes min len, got %v bytes", HandshakeMinBodySize, l),
		)
	}

	return nil
}

func (d Decoder) checkIfBodyLenLowerThanExpected(l, expected uint32) error {
	if l > expected {
		return errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid request: body (%v bytes) is lower than expected", l),
		)
	}
	return nil
}

func (d Decoder) ping(f frame.Frame, body []byte) (frame.Frame, error) {
	f.Body = bodies.Ping{}
	return f, nil
}

func (d Decoder) handshake(f frame.Frame, body []byte) (frame.Frame, error) {
	cursor := uint32(HandshakeLoginLenOffset)
	bodyLen := uint32(len(body))

	if e := d.checkIfBodyLenIsTooSmall(bodyLen, HandshakeMinBodySize); e != nil {
		return f, e
	}

	loginLen := d.decodeUint32(body[cursor:])
	cursor += bodies.LoginLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, loginLen+cursor); e != nil {
		return f, e
	}

	login := ""

	if loginLen != 0 {
		login = d.decodeString(body[HandshakeLoginOffset : HandshakeLoginOffset+loginLen])
	}

	cursor += loginLen

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+bodies.HashFieldSize); e != nil {
		return f, e
	}

	hash := make([]byte, 0, 32)
	copy(hash, body[cursor:cursor+bodies.HashFieldSize])

	cursor += bodies.HashFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+bodies.MaxCompressionsPerOneHandshake); e != nil {
		return f, e
	}

	compressionsLen := d.decodeUint8(body[cursor+bodies.HashFieldSize:])
	cursor += bodies.CompressionLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+uint32(compressionsLen)); e != nil {
		return f, e
	}

	compressions := make([]fields.Compression, 0, compressionsLen)

	for _, c := range body[cursor : cursor+uint32(compressionsLen)] {
		compressions = append(compressions, fields.Compression(c))
	}

	f.Body = bodies.NewHandshake(login, [32]byte(hash), compressions)

	return f, nil
}

func (d Decoder) decodeBytes(body []byte) []byte {
	key := make([]byte, 0, len(body))
	copy(key, body)

	return key
}

func (d Decoder) read(f frame.Frame, body []byte) (frame.Frame, error) {
	f.Body = bodies.Read(d.decodeBytes(body))
	return f, nil
}

func (d Decoder) delete(f frame.Frame, body []byte) (frame.Frame, error) {
	f.Body = bodies.Delete(d.decodeBytes(body))
	return f, nil
}

func (d Decoder) write(f frame.Frame, body []byte) (frame.Frame, error) {
	cursor := uint32(WriteKeyLenOffset)
	bodyLen := uint32(len(body))

	if e := d.checkIfBodyLenIsTooSmall(bodyLen, cursor+fields.KeyLenFieldSize); e != nil {
		return f, e
	}

	keyLen := d.decodeUint32(body[WriteKeyLenOffset:])
	cursor += fields.KeyLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+keyLen); e != nil {
		return f, e
	}

	key := fields.Key{}

	if keyLen != 0 {
		key = fields.Key(d.decodeBytes(body[cursor:]))
	}

	cursor += keyLen

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+fields.TypeFieldSize); e != nil {
		return f, e
	}

	valueType := fields.Type(d.decodeUint8(body[cursor:]))
	cursor += fields.TypeFieldSize

	switch valueType {
	case fields.Bytes:
		return bodies.Write{Key: key}
	}

	return bodies.Write{Key: key}, nil
}

func (d Decoder) decodeBytesValue(value []byte) (values.Bytes, error) {
	cursor := uint32(BytesLenOffsetInValue)
	valueLen := uint32(len(value))

	if e := d.checkIfBodyLenLowerThanExpected(valueLen, cursor+values.BytesLenFieldSize); e != nil {
		return values.Bytes{}, e
	}

	bytesLen := d.decodeUint32(value[cursor:])
	cursor += values.BytesLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(valueLen, cursor+bytesLen); e != nil {
		return values.Bytes{}, e
	}

	return values.Bytes(d.decodeBytes(value[cursor : cursor+bytesLen])), nil
}
