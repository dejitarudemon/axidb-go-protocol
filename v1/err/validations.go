package err

import "fmt"

var _ error = ValidationError{}

type ValidationError struct {
	msg     string
	details []any
}

func NewValidationError(msg string, args ...any) ValidationError {
	return ValidationError{
		msg:     msg,
		details: args,
	}
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%v. use Details() to get more info", v.msg)
}

func (v ValidationError) Details() []any {
	return v.details
}
