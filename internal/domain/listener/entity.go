package listener

import (
	"time"

	"hub/internal/domain/shared"
)

// Listener represents a user listening to a track.
type Listener struct {
	shared.AggregateRoot
	id        ListenerID
	userID    string
	trackID   string
	createdAt time.Time
}

// NewListener creates a new Listener entity.
func NewListener(userID, trackID string) (*Listener, error) {
	id, err := NewListenerID(userID, trackID)
	if err != nil {
		return nil, err
	}

	l := &Listener{
		id:        id,
		userID:    userID,
		trackID:   trackID,
		createdAt: time.Now(),
	}

	l.AddEvent(NewListenerTrackedEvent(userID, trackID))
	return l, nil
}

// ReconstructListener rebuilds a Listener from persistence.
func ReconstructListener(userID, trackID string, createdAt time.Time) (*Listener, error) {
	id, err := NewListenerID(userID, trackID)
	if err != nil {
		return nil, err
	}

	return &Listener{
		id:        id,
		userID:    userID,
		trackID:   trackID,
		createdAt: createdAt,
	}, nil
}

// ID returns the listener ID.
func (l *Listener) ID() ListenerID {
	return l.id
}

// UserID returns the user ID.
func (l *Listener) UserID() string {
	return l.userID
}

// TrackID returns the track ID.
func (l *Listener) TrackID() string {
	return l.trackID
}

// CreatedAt returns when the listener was tracked.
func (l *Listener) CreatedAt() time.Time {
	return l.createdAt
}
