package builder

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v0/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v0/err"
	"github.com/dejitarudemon/axidb-go-protocol/v0/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v0/frame"
)

// FrameBuilder builds Hello frames and rejects those larger than a size limit.
type FrameBuilder struct {
	limit int
}

// NewFrameBuilder returns a FrameBuilder that rejects frames whose encoded size exceeds limit bytes.
func NewFrameBuilder(limit int) FrameBuilder {
	return FrameBuilder{
		limit: limit,
	}
}

// frameSizeLowerLimit reports an error when size exceeds the builder limit.
func (fb FrameBuilder) frameSizeLowerLimit(size int) error {
	if size > fb.limit {
		return err.NewBuildError(fmt.Sprintf("frame size is %v, but limit is %v", size, fb.limit), nil)
	}

	return nil
}

// NewHello returns a Hello frame advertising versions.
// Version 0, values that do not fit in one byte, and duplicates are dropped.
// It returns an error when the body is invalid or the encoded frame exceeds the size limit.
func (fb FrameBuilder) NewHello(versions []fields.Version) (frame.Frame, error) {
	f := frame.Frame{
		Body: bodies.NewHello(versions),
	}

	if e := f.IsValid(); e != nil {
		return frame.Frame{}, e
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}
