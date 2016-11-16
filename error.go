/*
error.go contains the chainedError type.
*/

package main

/*
ChainedError holds a custom message and an "original" error, so you can report
additional information on an error as you catch it.
*/
type chainedError struct {
	Message string
	Cause   error
}

func (e *chainedError) Error() string {
	return "chained error: " + e.Message + "\ncause: " + e.Cause.Error()
}

/*
Return an error object which has a custom error message and is chained to an old
one.
*/
func ChainedError(cause error, message string) *chainedError {
	return &chainedError{Message: message, Cause: cause}
}
