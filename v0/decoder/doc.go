// Package decoder reads protocol v0 Hello frames from a buffered byte stream.
//
// [Decoder] checks magic bytes and the CRC-32/XFER checksum, then decodes the
// advertised versions into a [frame.Frame].
package decoder
