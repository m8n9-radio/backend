package reaction

import (
	"strings"

	"hub/internal/domain/shared"
)

// ErrInvalidUserID is returned when a user ID is empty.
var ErrInvalidUserID = shared.NewDomainError(
	shared.ErrInvalidInput,
	"user ID cannot be empty",
)

// UserID is a value object representing a user's unique identifier.
type UserID struct {
	value string
}

// NewUserID creates a new UserID from a string.
// Returns an error if the string is empty.
func NewUserID(value string) (UserID, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return UserID{}, ErrInvalidUserID
	}
	return UserID{value: trimmed}, nil
}

// String returns the string representation of the UserID.
func (id UserID) String() string {
	return id.value
}

// IsEmpty returns true if the UserID is empty.
func (id UserID) IsEmpty() bool {
	return id.value == ""
}

// Equals checks if two UserIDs are equal.
func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}
