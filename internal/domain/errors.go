package domain

// StructuredError is returned for invalid pattern, root unreadable, zone not found, etc.
// It is the domain error type used across application and adapters.
type StructuredError struct {
	Code    string
	Message string
}

func (e *StructuredError) Error() string {
	return e.Code + ": " + e.Message
}
