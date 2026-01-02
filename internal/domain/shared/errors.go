package shared

import "errors"

// Base domain errors
var (
	ErrNotFound        = errors.New("entity not found")
	ErrAlreadyExists   = errors.New("entity already exists")
	ErrInvalidInput    = errors.New("invalid input")
	ErrOperationFailed = errors.New("operation failed")
)

// DomainError wraps a domain error with additional context.
type DomainError struct {
	Err     error
	Message string
	Field   string
}

// Error implements the error interface.
func (e *DomainError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

// Unwrap returns the underlying error.
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error.
func NewDomainError(err error, message string) *DomainError {
	return &DomainError{
		Err:     err,
		Message: message,
	}
}

// NewFieldError creates a domain error for a specific field.
func NewFieldError(err error, field, message string) *DomainError {
	return &DomainError{
		Err:     err,
		Field:   field,
		Message: message,
	}
}
