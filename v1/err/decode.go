package err

import "fmt"

// DecodeError wraps an error raised while decoding a protocol message.
type DecodeError struct {
	msg    string
	source error
}

// NewDecodeError returns a DecodeError with message msg wrapping source.
func NewDecodeError(msg string, err error) DecodeError {
	return DecodeError{
		msg:    msg,
		source: err,
	}
}

// Error returns the decode failure message and the wrapped error.
func (d DecodeError) Error() string {
	return fmt.Sprintf("%v: %v", d.msg, d.source)
}

// Unwrap returns the underlying error.
func (d DecodeError) Unwrap() error {
	return d.source
}
