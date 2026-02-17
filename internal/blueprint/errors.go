package blueprint

// StructuredError is returned for invalid pattern, root unreadable, zone not found, etc.
type StructuredError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *StructuredError) Error() string {
	return e.Code + ": " + e.Message
}
