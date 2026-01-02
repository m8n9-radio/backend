package track

import (
	"hub/internal/domain/shared"
)

// ErrInvalidTrackID is returned when a track ID is not a valid MD5 hash.
var ErrInvalidTrackID = shared.NewDomainError(
	shared.ErrInvalidInput,
	"track ID must be a valid MD5 hash (32 hexadecimal characters)",
)

// TrackID is a value object representing a track's unique identifier.
// It must be a valid MD5 hash (32 hexadecimal characters).
type TrackID struct {
	value string
}

// NewTrackID creates a new TrackID from a string.
// Returns an error if the string is not a valid MD5 hash.
func NewTrackID(value string) (TrackID, error) {
	if !isValidMD5(value) {
		return TrackID{}, ErrInvalidTrackID
	}
	return TrackID{value: value}, nil
}

// String returns the string representation of the TrackID.
func (id TrackID) String() string {
	return id.value
}

// IsEmpty returns true if the TrackID is empty.
func (id TrackID) IsEmpty() bool {
	return id.value == ""
}

// Equals checks if two TrackIDs are equal.
func (id TrackID) Equals(other TrackID) bool {
	return id.value == other.value
}

// isValidMD5 checks if a string is a valid MD5 hash.
func isValidMD5(s string) bool {
	if len(s) != 32 {
		return false
	}
	for _, c := range s {
		if !isHexChar(c) {
			return false
		}
	}
	return true
}

// isHexChar checks if a rune is a valid hexadecimal character.
func isHexChar(c rune) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'f') ||
		(c >= 'A' && c <= 'F')
}
