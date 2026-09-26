// Package buffer defines buffers used to encode protocol v1 frames.
//
// [Appender] writes bytes into a buffer. [Buffer] adds allocation and readback.
// [Slice] is a byte-slice implementation. It is not safe for concurrent use.
package buffer
