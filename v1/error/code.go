package err

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

const CodeFieldSize = 2

type Code uint16

const (
	NoHello Code = iota
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

func (c Code) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(c))
}

func (c Code) String() string {
	switch c {
	case NoHello:
		return "No Hello"
	case UnexpectedCommand:
		return "UnexpectedCommand"
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
		return "Unkown Error Code"
	}

	return fmt.Sprintf("Unkown (%d)", c)
}
