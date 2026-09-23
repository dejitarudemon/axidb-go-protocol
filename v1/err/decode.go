package err

import "fmt"

type DecodeError struct {
	msg    string
	source error
}

func NewDecodeError(msg string, err error) DecodeError {
	return DecodeError{
		msg:    msg,
		source: err,
	}
}

func (d DecodeError) Error() string {
	return fmt.Sprintf("%v: %v", d.msg, d.source)
}

func (d DecodeError) Unwrap() error {
	return d.source
}
