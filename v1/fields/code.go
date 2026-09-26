package fields

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// ErrorFieldSize is the encoded size in bytes of an [Error] code.
const ErrorFieldSize = 2

// Error is a protocol v1 error code carried in error answers.
type Error uint16

// Error codes defined by protocol v1.
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
)

// Encode writes the wire encoding of the error code into buf.
func (e Error) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e))
}

// String returns a human-readable error name.
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
	}

	return fmt.Sprintf("Unknown (%d)", e)
}

// Size returns the encoded size in bytes.
func (e Error) Size() int {
	return ErrorFieldSize
}

// IsValid reports whether e is a known protocol error code (0 through [RestrictedRequest]).
func (e Error) IsValid() bool {
	return e <= RestrictedRequest
}
