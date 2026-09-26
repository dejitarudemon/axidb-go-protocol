// Package err provides protocol v1 error types and helpers for AxiDB.
//
// Wire-format errors implement [ProtocolError] and can be encoded into messages.
// Additional wrapper types ([BuildError], [DecodeError], and others) represent
// local encoding, decoding, and framing failures.
package err
