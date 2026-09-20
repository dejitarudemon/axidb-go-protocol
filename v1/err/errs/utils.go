package errs

import "github.com/google/uuid"

func generateNewTracebackID() uuid.UUID {
	return uuid.New()
}
