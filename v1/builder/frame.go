package builder

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

const (
	handshakeRequestID = fields.RequestID(0)
)

type FrameBuilder struct {
	limit int
}

func NewFrameBuilder(limit int) FrameBuilder {
	return FrameBuilder{
		limit: limit,
	}
}

func (fb FrameBuilder) requestIDIsNotZero(requestID fields.RequestID) error {
	if requestID == 0 {
		return err.NewBuildError("request ID is 0", nil)
	}

	return nil
}

func (fb FrameBuilder) frameSizeLowerLimit(size int) error {
	if size > fb.limit {
		return err.NewBuildError(fmt.Sprintf("frame size is %v, but limit is %v", size, fb.limit), nil)
	}

	return nil
}

func (fb FrameBuilder) NewHandshake(login string, hash [32]byte, compressions []fields.Compression) (frame.Frame, error) {
	f := frame.Frame{
		RequestID: handshakeRequestID,
		Body:      bodies.NewHandshake(login, hash, compressions),
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewHandshakeAnswer(compressions []fields.Compression) (frame.Frame, error) {
	f := frame.Frame{
		RequestID: handshakeRequestID,
		Body:      bodies.NewHandshakeAnswer(compressions),
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewRead(requestID fields.RequestID, key fields.Key) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.Read(key),
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewReadAnswer(requestID fields.RequestID, value value.V) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.ReadAnswer{Value: value},
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewWrite(requestID fields.RequestID, key fields.Key, value value.V) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.Write{Key: key, Value: value},
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewWriteAnswer(requestID fields.RequestID) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.WriteAnswer{},
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewDelete(requestID fields.RequestID, key fields.Key) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.Delete(key),
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewDeleteAnswer(requestID fields.RequestID) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.DeleteAnswer{},
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewPing(requestID fields.RequestID) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.Ping{},
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewPingAnswer(requestID fields.RequestID) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.Ping{},
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewErrAnswer(requestID fields.RequestID, pe err.ProtocolError) (frame.Frame, error) {
	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.ErrorAnswer{Err: pe},
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewBatch(requestID fields.RequestID, batch BatchRequestsBuilder) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	requests := make([]bodies.Request, len(batch.requests))
	copy(requests, batch.requests)

	f := frame.Frame{
		RequestID: requestID,
		Body: bodies.Batch{
			IsSequentialExecution: batch.isSequentialExecution,
			IsOneAnswer:           batch.isOneAnswer,
			InterruptAfterError:   batch.interruptAfterError,

			Requests: requests,
		},
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}

func (fb FrameBuilder) NewBatchAnswer(requestID fields.RequestID, batch BatchResultsBuilder) (frame.Frame, error) {
	if err := fb.requestIDIsNotZero(requestID); err != nil {
		return frame.Frame{}, err
	}

	results := make([]bodies.Result, len(batch.results))
	copy(results, batch.results)

	f := frame.Frame{
		RequestID: requestID,
		Body:      bodies.BatchAnswer(results),
	}

	if err := f.IsValid(); err != nil {
		return frame.Frame{}, err
	}

	return f, fb.frameSizeLowerLimit(f.Size())
}
