/*
package err предназначен для представления кодов
ошибок (Error Code) согласно спецификации протокола v1.

Использование:

	с = Code(1)
*/

package err

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
FieldSize представляет размер в байтах,
отведенный для хранения кода команды в сообщении.
*/
const FieldSize = 2

/*
type Code предназначен для хранения
кода ошибки, его валидации и кодирования в сообщение.
*/
type Code uint16

/*
Константы, представляющие ошибки,
используемые в спецификации протокола v1.
*/
const (
	NoHello Code = iota
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
	UnknownErrorCode
)

/*
func Encode предназначена для кодирования кода ошибки
в сообщении.

Принимааемые параметры:
  - buf buffer.Appender - буфер для хранения закодированного значения.
*/
func (c Code) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(c))
}

/*
func String предназначена для вывода человекочитаемого названия
ошибки, представленного конкретным кодом.
*/
func (c Code) String() string {
	switch c {
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
	case UnknownErrorCode:
		return "Unknown Error Code"
	}

	return fmt.Sprintf("Unknown (%d)", c)
}

/*
func Size возвращает размер кода ошибок в байтах.
*/
func (c Code) Size() int {
	return FieldSize
}

/*
func IsValid возвращает true, если код ошибки валиден.
В противном случае false.
*/
func (c Code) IsValid() bool {
	return c <= UnknownErrorCode
}
