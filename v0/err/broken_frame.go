package err

import "fmt"

// BrokenFrameError reports a frame that could not be parsed, including the raw bytes received.
type BrokenFrameError struct {
	got []byte
}

// NewBrokenFrameError returns a BrokenFrameError for the invalid frame bytes in got.
func NewBrokenFrameError(got []byte) BrokenFrameError {
	return BrokenFrameError{got: got}
}

// Error returns a hex dump prefix of the invalid frame bytes.
func (b BrokenFrameError) Error() string {
	return fmt.Sprintf("got: % X", b.got)
}
