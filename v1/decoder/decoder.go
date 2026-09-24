package decoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

const (
	MagicBytesLen = 2
	PreambleLen   = MagicBytesLen + fields.VersionFieldSize

	CommandOffset     = 3
	RequestIDOffest   = 4
	CompressionOffset = 8
	BodyLenOffset     = 9

	HandshakeMinBodySize    = 37
	WriteMinBodySize        = 5
	AnswerMinBodySize       = 2
	BatchMinBodySize        = 5
	BatchRequestMinBodySize = 9
	BatchAnswerMinBodySize  = 4
	BatchResultMinBodySize  = 8

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

func (d Decoder) DecodePreamble(reader *bufio.Reader) (fields.Version, error) {
	if reader == nil {
		return 0, err.NewDecodeError("got nil reader", nil)
	}

	buf, e := reader.Peek(MagicBytesLen + fields.VersionFieldSize)
	if e != nil {
		return 0, d.handleReaderError(e)
	}

	if !bytes.Equal(frame.MagicBytes, buf[:MagicBytesLen]) {
		return 0, err.NewBrokenFrameError(buf[:MagicBytesLen])
	}

	return fields.Version(buf[MagicBytesLen]), nil
}

func (d Decoder) DecodeFrame(reader *bufio.Reader) (frame.Frame, error) {
	if reader == nil {
		return frame.Frame{}, err.NewDecodeError("got nil reader", nil)
	}

	headersBuf := make([]byte, PreambleLen+frame.HeadersLen)

	if _, e := io.ReadFull(reader, headersBuf); e != nil {
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
	if !ok && compression != fields.None {
		return frame.Frame{}, errs.NewErrorUnsupportedCompression(compression)
	}

	bodyBuf := make([]byte, int(bodyLen))
	if _, e := io.ReadFull(reader, bodyBuf); e != nil {
		return frame.Frame{}, d.handleReaderError(e)
	}

	checksumBuf := make([]byte, fields.ChecksumFieldSize)
	if _, e := io.ReadFull(reader, checksumBuf); e != nil {
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

	body, e := d.decodeBody(bodyBuf, command)
	if e != nil {
		return f, e
	}

	f.Body = body
	return f, nil
}

func (d Decoder) decodeBody(body []byte, command fields.Command) (body.Body, error) {
	switch command {
	case fields.Handshake:
		return d.handshake(body)
	case fields.Read:
		return d.read(body)
	case fields.Delete:
		return d.delete(body)
	case fields.Ping:
		return d.ping()
	case fields.Write:
		return d.write(body)
	case fields.Answer:
		return d.answer(body)
	case fields.Batch:
		return d.batch(body)
	}

	return nil, err.NewDecodeError(fmt.Sprintf("unkown command: %v", command), nil)
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
	if l < expected {
		return errs.NewErrorMalformedValue(
			fmt.Sprintf("invalid request: body (%v bytes) is lower than expected (%v bytes)", l, expected),
		)
	}
	return nil
}

func (d Decoder) answer(body []byte) (body.Answer, error) {
	cursor := uint32(0)
	bodyLen := uint32(len(body))

	if e := d.checkIfBodyLenIsTooSmall(bodyLen, AnswerMinBodySize); e != nil {
		return nil, e
	}

	result := d.decodeUint8(body)
	cursor += fields.CommandFieldSize

	if result > 1 {
		return nil, errs.NewErrorMalformedValue(
			fmt.Sprintf("unkown result: %v", result),
		)
	}

	if result == 0 {
		return d.answerErr(body[cursor:])
	}

	originalCommand := fields.Command(d.decodeUint8(body[cursor:]))
	cursor += fields.CommandFieldSize

	switch originalCommand {
	case fields.Answer:
		return nil, errs.NewErrorMalformedValue("got answer for answer")
	case fields.Ping:
		return d.answerPing()
	case fields.Write:
		return d.answerWrite()
	case fields.Delete:
		return d.answerDelete()
	case fields.Read:
		return d.answerRead(body[cursor:])
	case fields.Handshake:
		return d.answerHandshake(body[cursor:])
	case fields.Batch:
		return d.answerBatch(body[cursor:])
	}

	return nil, err.NewDecodeError(fmt.Sprintf("unknown command in answer: %v", originalCommand), nil)
}

func (d Decoder) answerHandshake(body []byte) (bodies.HandshakeAnswer, error) {
	bodyLen := uint32(len(body))
	cursor := uint32(0)

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+bodies.CompressionLenFieldSize); e != nil {
		return bodies.HandshakeAnswer{}, e
	}

	compressionsLen := d.decodeUint8(body[cursor:])
	cursor += bodies.CompressionLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+uint32(compressionsLen)); e != nil {
		return bodies.HandshakeAnswer{}, e
	}

	compressions := make([]fields.Compression, 0, compressionsLen)
	for _, c := range body[cursor : cursor+uint32(compressionsLen)] {
		compressions = append(compressions, fields.Compression(c))
	}

	return bodies.NewHandshakeAnswer(compressions), nil
}

func (d Decoder) answerRead(body []byte) (bodies.ReadAnswer, error) {
	bodyLen := uint32(len(body))
	cursor := uint32(0)

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+fields.TypeFieldSize); e != nil {
		return bodies.ReadAnswer{}, e
	}

	valueType := fields.Type(d.decodeUint8(body[cursor:]))
	cursor += fields.TypeFieldSize

	value, _, e := d.decodeValue(body[cursor:], valueType)
	if e != nil {
		return bodies.ReadAnswer{}, e
	}

	return bodies.ReadAnswer{Value: value}, nil
}

func (d Decoder) answerDelete() (bodies.DeleteAnswer, error) {
	return bodies.DeleteAnswer{}, nil
}

func (d Decoder) answerWrite() (bodies.WriteAnswer, error) {
	return bodies.WriteAnswer{}, nil
}

func (d Decoder) answerPing() (bodies.PingAnswer, error) {
	return bodies.PingAnswer{}, nil
}

func (d Decoder) answerErr(body []byte) (bodies.ErrorAnswer, error) {
	cursor := uint32(0)
	bodyLen := uint32(len(body))

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, fields.ErrorFieldSize); e != nil {
		return bodies.ErrorAnswer{}, e
	}

	code := fields.Error(d.decodeUint16(body[cursor:]))
	cursor += fields.ErrorFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+fields.TracebackIDFieldSize); e != nil {
		return bodies.ErrorAnswer{}, e
	}

	tracebackID := fields.TracebackID(d.decodeBytes(body[cursor : cursor+fields.TracebackIDFieldSize]))
	cursor += fields.TracebackIDFieldSize

	pe := err.ProtocolError(nil)

	switch code {
	case fields.NoHello:
		pe = errs.NewErrorNoHelloWithTracebackID(tracebackID)
	case fields.UnsupportedVersion:
		pe = errs.NewErrorUnsupportedVersionWithTracebackID(0, tracebackID)
	case fields.UnsupportedCommand:
		pe = errs.NewErrorUnsupportedCommandWithTracebackID(0, tracebackID)
	case fields.RequestsConflict:
		pe = errs.NewErrorRequestInterruptedWithTracebackID(0, tracebackID)
	case fields.UnsupportedCompression:
		pe = errs.NewErrorUnsupportedCommandWithTracebackID(0, tracebackID)
	case fields.MismatchedChecksum:
		pe = errs.NewErrorMismatchedChecksumWithTracebackID(0, 0, tracebackID)
	case fields.InternalError:
		pe = errs.NewErrorInternalErrorWithTracebackID(errors.New(""), tracebackID)
	case fields.NotFound:
		pe = errs.NewErrorNotFoundWithTracebackID(fields.Key(" "), tracebackID)
	case fields.ProhibitedCompression:
		pe = errs.NewErrorProhibitedCompressionWithTracebackID(0, 0, tracebackID)
	case fields.RequestInterrupted:
		pe = errs.NewErrorRequestInterruptedWithTracebackID(0, tracebackID)
	case fields.InvalidRequestID:
		pe = errs.NewErrorInvalidRequestIDWithTracebackID(0, 0, tracebackID)
	case fields.Unauthorized:
		pe = errs.NewErrorUnauthorizedWithTracebackID([]byte{}, tracebackID)
	case fields.RestrictedRequest:
		pe = errs.NewErrorRestrictedRequestWithTracebackID(fields.Key{}, []byte{}, []byte{}, 0, 0, tracebackID)
	case fields.UnexpectedCommand:
		if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+fields.CommandFieldSize); e != nil {
			return bodies.ErrorAnswer{}, e
		}
		pe = errs.NewErrorUnexpectedCommandWithTracebackID(0, fields.Command(d.decodeUint8(body[cursor:])), tracebackID)
	case fields.BodyLimitIsExceeded:
		if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+fields.BodyLimitFieldSize); e != nil {
			return bodies.ErrorAnswer{}, e
		}
		pe = errs.NewErrorBodyLimitIsExceededWithTracebackID(0, fields.BodyLimit(d.decodeUint32(body[cursor:])), tracebackID)
	case fields.MalformedValue:
		pe = errs.NewErrorMalformedValueWithTracebackID(d.decodeString(body[cursor:]), tracebackID)
	case fields.BatchLimitIsExceeded:
		if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+fields.BatchLimitFieldSize); e != nil {
			return bodies.ErrorAnswer{}, e
		}
		pe = errs.NewErrorBatchLimitIsExceededWithTracebackID(0, fields.BatchLimit(d.decodeUint32(body[cursor:])), tracebackID)
	case fields.UnexpectedCommandInBatch:
		if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+fields.RequestNumberFieldSize); e != nil {
			return bodies.ErrorAnswer{}, e
		}
		pe = errs.NewErrorUnexpectedCommandInBatchWithTracebackID(0, fields.RequestNumber(d.decodeUint32(body[cursor:])), tracebackID)
	default:
		return bodies.ErrorAnswer{}, err.NewDecodeError(
			fmt.Sprintf("unknown err code: %v", code), nil,
		)
	}

	return bodies.ErrorAnswer{Err: pe}, nil
}

func (d Decoder) answerBatch(body []byte) (bodies.BatchAnswer, error) {
	bodyLen := uint32(len(body))
	cursor := uint32(0)

	if e := d.checkIfBodyLenIsTooSmall(bodyLen, BatchAnswerMinBodySize); e != nil {
		return bodies.BatchAnswer{}, e
	}

	resultsLen := d.decodeUint32(body[cursor:])
	cursor += bodies.RequestsLenFieldSize

	results := make([]bodies.Result, 0, resultsLen)

	for range resultsLen {
		result, i, e := d.batchResult(body[cursor:])
		if e != nil {
			return bodies.BatchAnswer{}, e
		}

		results = append(results, result)
		cursor += uint32(i)
	}

	return bodies.BatchAnswer(results), nil
}

func (d Decoder) batchResult(body []byte) (bodies.Result, int, error) {
	cursor := uint32(0)
	bodyLen := uint32(len(body))

	if e := d.checkIfBodyLenIsTooSmall(bodyLen, BatchResultMinBodySize); e != nil {
		return bodies.Result{}, 0, e
	}

	requestNumber := fields.RequestNumber(d.decodeUint32(body[cursor:]))
	cursor += fields.RequestNumberFieldSize

	resultBodyLen := d.decodeUint32(body[cursor:])
	cursor += bodies.RequestBodyLenFieldSize

	b, e := d.answer(body[cursor : cursor+resultBodyLen])
	if e != nil {
		return bodies.Result{}, 0, e
	}

	cursor += resultBodyLen

	return bodies.Result{
		Number: requestNumber,
		Body:   b,
	}, int(cursor), nil

}

func (d Decoder) batch(body []byte) (bodies.Batch, error) {
	bodyLen := uint32(len(body))
	cursor := uint32(0)

	if e := d.checkIfBodyLenIsTooSmall(bodyLen, BatchMinBodySize); e != nil {
		return bodies.Batch{}, e
	}

	flags := d.decodeUint8(body[cursor:])
	cursor += bodies.FlagsFieldSize

	requestsLen := d.decodeUint32(body[cursor:])
	cursor += bodies.RequestsLenFieldSize

	requests := make([]bodies.Request, 0, requestsLen)

	for range requestsLen {
		request, i, e := d.batchRequest(body[cursor:])
		if e != nil {
			return bodies.Batch{}, e
		}

		requests = append(requests, request)
		cursor += uint32(i)
	}

	return bodies.Batch{
		Requests:              requests,
		IsSequentialExecution: flags&bodies.IsSequentialExecution != 0,
		InterruptAfterError:   flags&bodies.InterruptAfterError != 0,
		IsOneAnswer:           flags&bodies.IsOneAnswer != 0,
	}, nil
}

func (d Decoder) batchRequest(body []byte) (bodies.Request, int, error) {
	cursor := uint32(0)
	bodyLen := uint32(len(body))

	if e := d.checkIfBodyLenIsTooSmall(bodyLen, BatchRequestMinBodySize); e != nil {
		return bodies.Request{}, 0, e
	}

	requestNumber := fields.RequestNumber(d.decodeUint32(body[cursor:]))
	cursor += fields.RequestNumberFieldSize

	command := fields.Command(d.decodeUint8(body[cursor:]))
	cursor += fields.CommandFieldSize

	if command == fields.Batch {
		return bodies.Request{}, 0, errs.NewErrorUnexpectedCommandInBatch(command, requestNumber)
	}

	requestBodyLen := d.decodeUint32(body[cursor:])
	cursor += bodies.RequestBodyLenFieldSize

	b, e := d.decodeBody(body[cursor:cursor+requestBodyLen], command)
	if e != nil {
		return bodies.Request{}, 0, e
	}

	cursor += requestBodyLen

	return bodies.Request{
		Number: requestNumber,
		Body:   b,
	}, int(cursor), nil

}

func (d Decoder) ping() (bodies.Ping, error) {
	return bodies.Ping{}, nil
}

func (d Decoder) handshake(body []byte) (bodies.Handshake, error) {
	cursor := uint32(0)
	bodyLen := uint32(len(body))

	if e := d.checkIfBodyLenIsTooSmall(bodyLen, HandshakeMinBodySize); e != nil {
		return bodies.Handshake{}, e
	}

	loginLen := d.decodeUint32(body[cursor:])
	cursor += bodies.LoginLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, loginLen+cursor); e != nil {
		return bodies.Handshake{}, e
	}

	login := ""

	if loginLen != 0 {
		login = d.decodeString(body[cursor : cursor+loginLen])
	}

	cursor += loginLen

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+bodies.HashFieldSize); e != nil {
		return bodies.Handshake{}, e
	}

	hash := body[cursor : cursor+bodies.HashFieldSize]

	cursor += bodies.HashFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+bodies.CompressionLenFieldSize); e != nil {
		return bodies.Handshake{}, e
	}

	compressionsLen := d.decodeUint8(body[cursor:])
	cursor += bodies.CompressionLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+uint32(compressionsLen)); e != nil {
		return bodies.Handshake{}, e
	}

	compressions := make([]fields.Compression, 0, compressionsLen)

	for _, c := range body[cursor : cursor+uint32(compressionsLen)] {
		compressions = append(compressions, fields.Compression(c))
	}

	return bodies.NewHandshake(login, [32]byte(hash), compressions), nil
}

func (d Decoder) read(body []byte) (bodies.Read, error) {
	return bodies.Read(d.decodeBytes(body)), nil
}

func (d Decoder) delete(body []byte) (bodies.Delete, error) {
	return bodies.Delete(d.decodeBytes(body)), nil
}

func (d Decoder) write(body []byte) (bodies.Write, error) {
	cursor := uint32(0)
	bodyLen := uint32(len(body))

	if e := d.checkIfBodyLenIsTooSmall(bodyLen, cursor+fields.KeyLenFieldSize); e != nil {
		return bodies.Write{}, e
	}

	keyLen := d.decodeUint32(body[cursor:])
	cursor += fields.KeyLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+keyLen); e != nil {
		return bodies.Write{}, e
	}

	key := fields.Key{}

	if keyLen != 0 {
		key = fields.Key(d.decodeBytes(body[cursor : cursor+keyLen]))
	}

	cursor += keyLen

	if e := d.checkIfBodyLenLowerThanExpected(bodyLen, cursor+fields.TypeFieldSize); e != nil {
		return bodies.Write{}, e
	}

	valueType := fields.Type(d.decodeUint8(body[cursor:]))

	if !valueType.IsValid() {
		return bodies.Write{}, errs.NewErrorMalformedValue(
			fmt.Sprintf("unknown value type: %v", valueType),
		)
	}

	cursor += fields.TypeFieldSize

	value := value.V(nil)
	e := error(nil)

	switch valueType {
	case fields.Bytes:
		value, _, e = d.decodeBytesValue(body[cursor:])
	case fields.TypedArray:
		value, _, e = d.decodeTypedArray(body[cursor:])
	case fields.UntypedArray:
		value, _, e = d.decodeUntypedArray(body[cursor:])
	case fields.Int:
		value, _, e = d.decodeIntValue(body[cursor:])
	case fields.Uint:
		value, _, e = d.decodeUintValue(body[cursor:])
	case fields.Float:
		value, _, e = d.decodeFloatValue(body[cursor:])
	case fields.String:
		value, _, e = d.decodeStringValue(body[cursor:])
	case fields.JSON:
		value, _, e = d.decodeJSONValue(body[cursor:])
	}

	if e != nil {
		return bodies.Write{}, e
	}

	return bodies.Write{
		Key:   key,
		Value: value,
	}, nil
}

func (d Decoder) decodeBytesValue(value []byte) (values.Bytes, int, error) {
	cursor := uint32(0)
	valueLen := uint32(len(value))

	if e := d.checkIfBodyLenLowerThanExpected(valueLen, cursor+values.BytesLenFieldSize); e != nil {
		return values.Bytes{}, 0, e
	}

	bytesLen := d.decodeUint32(value[cursor:])
	cursor += values.BytesLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(valueLen, cursor+bytesLen); e != nil {
		return values.Bytes{}, 0, e
	}

	return values.Bytes(d.decodeBytes(value[cursor : cursor+bytesLen])), int(cursor + bytesLen), nil
}

func (d Decoder) decodeIntValue(value []byte) (values.Int, int, error) {
	if e := d.checkIfBodyLenLowerThanExpected(uint32(len(value)), values.IntValueFieldSize); e != nil {
		return values.Int(0), 0, e
	}

	return values.Int(d.decodeInt64(value)), values.IntValueFieldSize, nil
}

func (d Decoder) decodeUintValue(value []byte) (values.Uint, int, error) {
	if e := d.checkIfBodyLenLowerThanExpected(uint32(len(value)), values.UintValueFieldSize); e != nil {
		return values.Uint(0), 0, e
	}

	return values.Uint(d.decodeUint64(value)), values.UintValueFieldSize, nil
}

func (d Decoder) decodeFloatValue(value []byte) (values.Float, int, error) {
	if e := d.checkIfBodyLenLowerThanExpected(uint32(len(value)), values.FloatValueFieldSize); e != nil {
		return values.Float(0), 0, e
	}

	return values.Float(d.decodeFloat64(value)), values.FloatValueFieldSize, nil
}

func (d Decoder) decodeUntypedArray(data []byte) (values.UntypedArray, int, error) {
	cursor := uint32(0)
	lenValue := uint32(len(data))

	if e := d.checkIfBodyLenLowerThanExpected(lenValue, cursor+values.TypedArrayLenFieldSize); e != nil {
		return values.UntypedArray{}, 0, e
	}

	arrayLen := d.decodeUint32(data[cursor:])
	cursor += values.TypedArrayLenFieldSize

	elems := make([]value.V, 0, arrayLen)

	for range arrayLen {
		if e := d.checkIfBodyLenLowerThanExpected(lenValue, cursor+fields.TypeFieldSize); e != nil {
			return values.UntypedArray{}, 0, e
		}

		elemType := fields.Type(d.decodeUint8(data[cursor:]))
		cursor += fields.TypeFieldSize

		elem, i, e := d.decodeValue(data[cursor:], elemType)

		if e != nil {
			return values.UntypedArray{}, 0, e
		}

		elems = append(elems, elem)
		cursor += uint32(i)
	}

	return values.UntypedArray(elems), int(cursor), nil
}

func (d Decoder) decodeTypedArray(data []byte) (values.TypedArray, int, error) {
	cursor := uint32(0)
	lenValue := uint32(len(data))

	if e := d.checkIfBodyLenLowerThanExpected(lenValue, cursor+values.TypedArrayLenFieldSize); e != nil {
		return values.TypedArray{}, 0, e
	}

	arrayLen := d.decodeUint32(data[cursor:])
	cursor += values.TypedArrayLenFieldSize

	if e := d.checkIfBodyLenLowerThanExpected(lenValue, cursor+fields.TypeFieldSize); e != nil {
		return values.TypedArray{}, 0, e
	}

	elemType := fields.Type(d.decodeUint8(data[cursor:]))
	elems := make([]value.V, 0, arrayLen)

	cursor += fields.TypeFieldSize

	for range arrayLen {
		elem, i, e := d.decodeValue(data[cursor:], elemType)

		if e != nil {
			return values.TypedArray{}, 0, e
		}

		elems = append(elems, elem)
		cursor += uint32(i)
	}

	return values.TypedArray{ElemType: elemType, Elems: elems}, int(cursor), nil
}

func (d Decoder) decodeValue(value []byte, valueType fields.Type) (value.V, int, error) {
	switch valueType {
	case fields.Bytes:
		return d.decodeBytesValue(value)
	case fields.TypedArray:
		return d.decodeTypedArray(value)
	case fields.UntypedArray:
		return d.decodeUntypedArray(value)
	case fields.Int:
		return d.decodeIntValue(value)
	case fields.Uint:
		return d.decodeUintValue(value)
	case fields.Float:
		return d.decodeFloatValue(value)
	case fields.String:
		return d.decodeStringValue(value)
	case fields.JSON:
		return d.decodeJSONValue(value)
	}

	return nil, 0, err.NewDecodeError(fmt.Sprintf("unkown type: %v", valueType), nil)
}

func (d Decoder) decodeStringValue(value []byte) (values.String, int, error) {
	b, i, e := d.decodeBytesValue(value)
	return values.String(b), i, e
}

func (d Decoder) decodeJSONValue(value []byte) (values.JSON, int, error) {
	b, i, e := d.decodeBytesValue(value)
	return values.JSON(b), i, e
}

func (d Decoder) decodeFloat64(data []byte) float64 {
	return math.Float64frombits(d.decodeUint64(data))
}

func (d Decoder) decodeUint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

func (d Decoder) decodeString(data []byte) string {
	return string(data)
}

func (d Decoder) decodeUint8(data []byte) uint8 {
	return uint8(data[0])
}

func (d Decoder) decodeBytes(body []byte) []byte {
	return body
}

func (d Decoder) decodeUint64(data []byte) uint64 {
	return binary.BigEndian.Uint64(data)
}

func (d Decoder) decodeInt64(data []byte) int64 {
	return int64(d.decodeUint64(data))
}

func (d Decoder) decodeUint16(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}
