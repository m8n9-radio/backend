package track

import "strings"

// Cover is a value object representing a track's cover image URL.
// Cover is optional and can be empty.
type Cover struct {
	value string
}

// NewCover creates a new Cover from a string.
// Empty values are allowed as cover is optional.
func NewCover(value string) Cover {
	return Cover{value: strings.TrimSpace(value)}
}

// String returns the string representation of the Cover.
func (c Cover) String() string {
	return c.value
}

// IsEmpty returns true if the Cover is empty.
func (c Cover) IsEmpty() bool {
	return c.value == ""
}

// Equals checks if two Covers are equal.
func (c Cover) Equals(other Cover) bool {
	return c.value == other.value
}
