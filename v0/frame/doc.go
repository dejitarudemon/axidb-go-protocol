// Package frame defines the protocol v0 Hello wire frame and its encoding.
//
// A [Frame] carries a [body.Body] listing supported versions, prefixed by magic
// bytes, version 0, a version count, and a CRC-32/XFER checksum.
package frame
