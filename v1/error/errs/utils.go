package errs

import "uuid"

func generateNewTracebackID() uuid.UUID {
	return uuid.New()
}
