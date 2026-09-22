package fields

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
ErrorFieldSize представляет размер в байтах,
отведенный для хранения кода команды в сообщении.
*/
const ErrorFieldSize = 2

/*
type Error предназначен для хранения
кода ошибки, его валидации и кодирования в сообщение.
*/
type Error uint16

/*
Константы, представляющие ошибки,
используемые в спецификации протокола v1.
*/
const (
	NoHello Error = iota
	UnsupportedVersion
	UnexpectedCommand
	UnsupportedCommand
	RequestsConflict
	UnsupportedCompression
	BodyLimitIsExceeded
	MismatchedChecksum
	InternalError
	MalformedValue
	NotFound
	ProhibitedCompression
	RequestInterrupted
	BatchLimitIsExceeded
	UnexpectedCommandInBatch
	InvalidRequestID
	Unauthorized
	RestrictedRequest
	UnknownErrorError
)

/*
func Encode предназначена для кодирования кода ошибки
в сообщении.

Принимааемые параметры:
  - buf buffer.Appender - буфер для хранения закодированного значения.
*/
func (e Error) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e))
}

/*
func String предназначена для вывода человекочитаемого названия
ошибки, представленного конкретным кодом.
*/
func (e Error) String() string {
	switch e {
	case NoHello:
		return "No Hello"
	case UnsupportedVersion:
		return "Unsupported Version"
	case UnexpectedCommand:
		return "Unexpected Command"
	case UnsupportedCommand:
		return "Unsupported Command"
	case RequestsConflict:
		return "Requests Conflict"
	case UnsupportedCompression:
		return "Unsupported Compression"
	case BodyLimitIsExceeded:
		return "Body Limit Is Exceeded"
	case MismatchedChecksum:
		return "Mismatched Checksum"
	case InternalError:
		return "Internal Error"
	case MalformedValue:
		return "Malformed Value"
	case NotFound:
		return "Not Found"
	case ProhibitedCompression:
		return "Prohibited Compression"
	case RequestInterrupted:
		return "Request Interrupted"
	case BatchLimitIsExceeded:
		return "Batch Limit Is Exceeded"
	case UnexpectedCommandInBatch:
		return "Unexpected Command In Batch"
	case InvalidRequestID:
		return "Invalid Request ID"
	case Unauthorized:
		return "Unauthorized"
	case RestrictedRequest:
		return "Restricted Request"
	case UnknownErrorError:
		return "Unknown Error Error"
	}

	return fmt.Sprintf("Unknown (%d)", e)
}

/*
func Size возвращает размер кода ошибок в байтах.
*/
func (e Error) Size() int {
	return ErrorFieldSize
}

/*
func IsValid возвращает true, если код ошибки валиден.
В противном случае false.
*/
func (e Error) IsValid() bool {
	return e <= UnknownErrorError
}
