package track

import (
	"strings"

	"hub/internal/domain/shared"
)

// ErrInvalidTitle is returned when a title is empty or invalid.
var ErrInvalidTitle = shared.NewDomainError(
	shared.ErrInvalidInput,
	"title cannot be empty",
)

// Title is a value object representing a track's title.
type Title struct {
	value string
}

// NewTitle creates a new Title from a string.
// Returns an error if the string is empty or whitespace only.
func NewTitle(value string) (Title, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Title{}, ErrInvalidTitle
	}
	return Title{value: trimmed}, nil
}

// String returns the string representation of the Title.
func (t Title) String() string {
	return t.value
}

// IsEmpty returns true if the Title is empty.
func (t Title) IsEmpty() bool {
	return t.value == ""
}

// Equals checks if two Titles are equal.
func (t Title) Equals(other Title) bool {
	return t.value == other.value
}
