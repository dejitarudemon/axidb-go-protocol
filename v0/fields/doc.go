// Package fields defines wire types from the AxiDB protocol v0 specification.
//
// Version 0 is the Hello frame used to agree on a working protocol version.
// Most types support [buffer.Appender] encoding via Encode and report their encoded width with Size.
package fields
