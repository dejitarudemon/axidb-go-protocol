package err

import "fmt"

type BrokenFrameError struct {
	got []byte
}

func NewBrokenFrameError(got []byte) BrokenFrameError {
	return BrokenFrameError{got: got}
}

func (b BrokenFrameError) Error() string {
	return fmt.Sprintf("got: % X", b.got)
}
