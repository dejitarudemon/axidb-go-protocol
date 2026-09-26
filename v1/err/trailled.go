package err

import "fmt"

type TrailledError struct {
	expected uint32
	consumed uint32
}

func NewTrailledError(expected, consumed uint32) TrailledError {
	return TrailledError{
		expected: expected,
		consumed: consumed,
	}
}

func (t TrailledError) Error() string {
	return fmt.Sprintf("trailled data. expected %v bytes, but consumed %v", t.expected, t.consumed)
}
