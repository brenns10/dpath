package main

type chainedError struct {
	Message string
	Cause   error
}

func (e *chainedError) Error() string {
	return "chained error: " + e.Message + "\ncause: " + e.Cause.Error()
}

func ChainedError(cause error, message string) *chainedError {
	return &chainedError{Message: message, Cause: cause}
}
