package err

import "fmt"

type BuildError struct {
	msg    string
	source error
}

func NewBuildError(msg string, err error) BuildError {
	return BuildError{
		msg:    msg,
		source: err,
	}
}

func (d BuildError) Error() string {
	return fmt.Sprintf("%v: %v", d.msg, d.source)
}

func (d BuildError) Unwrap() error {
	return d.source
}
