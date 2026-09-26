// Package buffer defines buffers used to encode protocol v1 frames.
//
// [Appender] writes bytes into a buffer. [Buffer] adds allocation and readback.
// [Slice] is a mutex-protected byte-slice implementation.
package buffer
