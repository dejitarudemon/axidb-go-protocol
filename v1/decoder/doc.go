// Package decoder reads protocol v1 frames from a buffered byte stream.
//
// [Decoder] checks magic bytes, checksums, body limits, and optional compression,
// then decodes the command payload into a [frame.Frame].
package decoder
