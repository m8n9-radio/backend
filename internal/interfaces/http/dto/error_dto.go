package dto

// ErrorResponse represents an HTTP error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new error response.
func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error:   code,
		Message: message,
	}
}

// Common error responses
var (
	ErrBadRequest = func(msg string) ErrorResponse {
		return NewErrorResponse("bad_request", msg)
	}
	ErrNotFound = func(msg string) ErrorResponse {
		return NewErrorResponse("not_found", msg)
	}
	ErrConflict = func(msg string) ErrorResponse {
		return NewErrorResponse("conflict", msg)
	}
	ErrInternalServer = NewErrorResponse("internal_error", "Internal server error")
)
