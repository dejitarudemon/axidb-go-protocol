package err

import "fmt"

var _ error = ValidationError{}

// ValidationError reports inconsistent error payload fields, typically from IsValid checks.
type ValidationError struct {
	msg     string
	details []any
}

// NewValidationError returns a ValidationError with msg and optional detail key-value pairs in args.
func NewValidationError(msg string, args ...any) ValidationError {
	return ValidationError{
		msg:     msg,
		details: args,
	}
}

// Error returns a summary message; call Details for structured fields.
func (v ValidationError) Error() string {
	return fmt.Sprintf("%v. use Details() to get more info", v.msg)
}

// Details returns structured detail fields attached to the validation failure.
func (v ValidationError) Details() []any {
	return v.details
}
