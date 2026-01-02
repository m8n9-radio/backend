package listener

import "hub/internal/domain/shared"

const (
	EventListenerTracked = "listener.tracked"
)

// ListenerTrackedEvent is emitted when a listener is tracked.
type ListenerTrackedEvent struct {
	shared.BaseEvent
	userID  string
	trackID string
}

// NewListenerTrackedEvent creates a new ListenerTrackedEvent.
func NewListenerTrackedEvent(userID, trackID string) *ListenerTrackedEvent {
	return &ListenerTrackedEvent{
		BaseEvent: shared.NewBaseEvent(EventListenerTracked),
		userID:    userID,
		trackID:   trackID,
	}
}

// Payload returns the event payload.
func (e *ListenerTrackedEvent) Payload() interface{} {
	return map[string]string{
		"user_id":  e.userID,
		"track_id": e.trackID,
	}
}

// UserID returns the user ID.
func (e *ListenerTrackedEvent) UserID() string {
	return e.userID
}

// TrackID returns the track ID.
func (e *ListenerTrackedEvent) TrackID() string {
	return e.trackID
}
