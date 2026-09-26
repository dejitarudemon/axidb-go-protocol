// Package errs defines concrete [err.ProtocolError] values for each protocol v1 error code.
//
// Each type can be validated with IsValid, encoded on the wire, and carries a traceback ID
// for correlation. Constructors exist with a generated traceback ID or with an explicit one.
package errs
