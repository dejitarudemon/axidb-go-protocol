package errs

import (
	"github.com/dejitarudemon/ignicula-wire/v1/fields"
	"github.com/google/uuid"
)

// generateNewTracebackID returns a new random traceback ID.
func generateNewTracebackID() fields.TracebackID {
	return fields.TracebackID(uuid.New())
}
