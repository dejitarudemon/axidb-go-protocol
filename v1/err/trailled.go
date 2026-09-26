package err

import "fmt"

// TrailledError reports trailing bytes after the expected amount of frame data was consumed.
type TrailledError struct {
	expected uint32
	consumed uint32
}

// NewTrailledError returns a TrailledError for expected total bytes versus consumed bytes.
func NewTrailledError(expected, consumed uint32) TrailledError {
	return TrailledError{
		expected: expected,
		consumed: consumed,
	}
}

// Error returns a human-readable description of the trailing data mismatch.
func (t TrailledError) Error() string {
	return fmt.Sprintf("trailled data. expected %v bytes, but consumed %v", t.expected, t.consumed)
}
