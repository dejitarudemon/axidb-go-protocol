package err

import "fmt"

// BuildError wraps an error raised while building or assembling a protocol message.
type BuildError struct {
	msg    string
	source error
}

// NewBuildError returns a BuildError with message msg wrapping source.
func NewBuildError(msg string, err error) BuildError {
	return BuildError{
		msg:    msg,
		source: err,
	}
}

// Error returns the build failure message and the wrapped error.
func (d BuildError) Error() string {
	return fmt.Sprintf("%v: %v", d.msg, d.source)
}

// Unwrap returns the underlying error.
func (d BuildError) Unwrap() error {
	return d.source
}
