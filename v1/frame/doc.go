// Package frame defines the protocol v1 wire frame and its encoding.
//
// A [Frame] combines a request ID with a [body.Body] payload, prefixed by magic bytes,
// version, compression metadata, and a checksum.
package frame
