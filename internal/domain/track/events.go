package track

import (
	"hub/internal/domain/shared"
)

// TrackCreated is emitted when a new track is created.
type TrackCreated struct {
	shared.BaseEvent
	trackID TrackID
	title   Title
	cover   Cover
}

// NewTrackCreated creates a new TrackCreated event.
func NewTrackCreated(id TrackID, title Title, cover Cover) TrackCreated {
	return TrackCreated{
		BaseEvent: shared.NewBaseEvent("track.created"),
		trackID:   id,
		title:     title,
		cover:     cover,
	}
}

// Payload returns the event data.
func (e TrackCreated) Payload() interface{} {
	return map[string]interface{}{
		"track_id": e.trackID.String(),
		"title":    e.title.String(),
		"cover":    e.cover.String(),
	}
}

// TrackID returns the track ID.
func (e TrackCreated) TrackID() TrackID { return e.trackID }

// Title returns the track title.
func (e TrackCreated) Title() Title { return e.title }

// Cover returns the track cover.
func (e TrackCreated) Cover() Cover { return e.cover }

// TrackRotated is emitted when a track's rotation count is incremented.
type TrackRotated struct {
	shared.BaseEvent
	trackID   TrackID
	newRotate int
}

// NewTrackRotated creates a new TrackRotated event.
func NewTrackRotated(id TrackID, newRotate int) TrackRotated {
	return TrackRotated{
		BaseEvent: shared.NewBaseEvent("track.rotated"),
		trackID:   id,
		newRotate: newRotate,
	}
}

// Payload returns the event data.
func (e TrackRotated) Payload() interface{} {
	return map[string]interface{}{
		"track_id":   e.trackID.String(),
		"new_rotate": e.newRotate,
	}
}

// TrackID returns the track ID.
func (e TrackRotated) TrackID() TrackID { return e.trackID }

// NewRotate returns the new rotation count.
func (e TrackRotated) NewRotate() int { return e.newRotate }

// CoverUpdated is emitted when a track's cover is updated.
type CoverUpdated struct {
	shared.BaseEvent
	trackID  TrackID
	oldCover Cover
	newCover Cover
}

// NewCoverUpdated creates a new CoverUpdated event.
func NewCoverUpdated(id TrackID, oldCover, newCover Cover) CoverUpdated {
	return CoverUpdated{
		BaseEvent: shared.NewBaseEvent("track.cover_updated"),
		trackID:   id,
		oldCover:  oldCover,
		newCover:  newCover,
	}
}

// Payload returns the event data.
func (e CoverUpdated) Payload() interface{} {
	return map[string]interface{}{
		"track_id":  e.trackID.String(),
		"old_cover": e.oldCover.String(),
		"new_cover": e.newCover.String(),
	}
}

// TrackID returns the track ID.
func (e CoverUpdated) TrackID() TrackID { return e.trackID }

// OldCover returns the old cover.
func (e CoverUpdated) OldCover() Cover { return e.oldCover }

// NewCover returns the new cover.
func (e CoverUpdated) NewCover() Cover { return e.newCover }
