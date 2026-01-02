package listener

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// ListenerID uniquely identifies a listener (user + track combination).
type ListenerID struct {
	value   string
	userID  string
	trackID string
}

// NewListenerID creates a new ListenerID from userID and trackID.
func NewListenerID(userID, trackID string) (ListenerID, error) {
	if userID == "" {
		return ListenerID{}, ErrInvalidUserID
	}
	if trackID == "" {
		return ListenerID{}, ErrInvalidTrackID
	}

	hash := generateHash(userID, trackID)
	return ListenerID{
		value:   hash,
		userID:  userID,
		trackID: trackID,
	}, nil
}

// String returns the string representation.
func (id ListenerID) String() string {
	return id.value
}

// UserID returns the user ID component.
func (id ListenerID) UserID() string {
	return id.userID
}

// TrackID returns the track ID component.
func (id ListenerID) TrackID() string {
	return id.trackID
}

// IsEmpty returns true if the ID is empty.
func (id ListenerID) IsEmpty() bool {
	return id.value == ""
}

// Equals compares two ListenerIDs.
func (id ListenerID) Equals(other ListenerID) bool {
	return id.value == other.value
}

func generateHash(userID, trackID string) string {
	data := fmt.Sprintf("%s:%s", userID, trackID)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
