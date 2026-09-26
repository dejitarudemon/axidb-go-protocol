package err

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v0/fields"
)

// MismatchedChecksum reports that a Hello frame checksum did not match the computed value.
type MismatchedChecksum struct {
	got      fields.Checksum
	expected fields.Checksum
}

// NewMismatchedChecksum returns a MismatchedChecksum for the received and computed values.
func NewMismatchedChecksum(got, expected fields.Checksum) MismatchedChecksum {
	return MismatchedChecksum{
		got:      got,
		expected: expected,
	}
}

// Error returns a human-readable summary of the checksum mismatch.
func (e MismatchedChecksum) Error() string {
	return fmt.Sprintf("mismatched checksum: got %08X, expected %08X", e.got, e.expected)
}
