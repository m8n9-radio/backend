package track

import (
	"time"

	"hub/internal/domain/shared"
)

// Track is the aggregate root for track-related operations.
type Track struct {
	shared.AggregateRoot

	id        TrackID
	title     Title
	cover     Cover
	rotate    int
	likes     int
	dislikes  int
	listeners int
	createdAt time.Time
	updatedAt time.Time
}

// NewTrack creates a new Track aggregate.
// This is the only way to create a new track.
func NewTrack(id TrackID, title Title, cover Cover) *Track {
	now := time.Now()
	t := &Track{
		id:        id,
		title:     title,
		cover:     cover,
		rotate:    1,
		likes:     0,
		dislikes:  0,
		listeners: 0,
		createdAt: now,
		updatedAt: now,
	}

	t.AddEvent(NewTrackCreated(id, title, cover))
	return t
}

// ReconstructTrack rebuilds a Track from persistence data.
// No events are emitted during reconstruction.
func ReconstructTrack(
	id, title, cover string,
	rotate, likes, dislikes, listeners int,
	createdAt, updatedAt time.Time,
) (*Track, error) {
	trackID, err := NewTrackID(id)
	if err != nil {
		return nil, err
	}

	trackTitle, err := NewTitle(title)
	if err != nil {
		return nil, err
	}

	return &Track{
		id:        trackID,
		title:     trackTitle,
		cover:     NewCover(cover),
		rotate:    rotate,
		likes:     likes,
		dislikes:  dislikes,
		listeners: listeners,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// IncrementRotation increases the rotation count by 1.
func (t *Track) IncrementRotation() {
	t.rotate++
	t.updatedAt = time.Now()
	t.AddEvent(NewTrackRotated(t.id, t.rotate))
}

// UpdateCover updates the cover if the current cover is empty.
// Returns true if the cover was updated.
func (t *Track) UpdateCover(newCover Cover) bool {
	if !t.cover.IsEmpty() {
		return false
	}
	if newCover.IsEmpty() {
		return false
	}

	oldCover := t.cover
	t.cover = newCover
	t.updatedAt = time.Now()
	t.AddEvent(NewCoverUpdated(t.id, oldCover, newCover))
	return true
}

// RecordLike increments the like count.
func (t *Track) RecordLike() {
	t.likes++
	t.updatedAt = time.Now()
}

// RecordDislike increments the dislike count.
func (t *Track) RecordDislike() {
	t.dislikes++
	t.updatedAt = time.Now()
}

// UpdateListenerCount sets the current listener count.
func (t *Track) UpdateListenerCount(count int) {
	t.listeners = count
	t.updatedAt = time.Now()
}

// Getters - no setters to enforce encapsulation

// ID returns the track's unique identifier.
func (t *Track) ID() TrackID { return t.id }

// Title returns the track's title.
func (t *Track) Title() Title { return t.title }

// Cover returns the track's cover URL.
func (t *Track) Cover() Cover { return t.cover }

// Rotate returns the rotation count.
func (t *Track) Rotate() int { return t.rotate }

// Likes returns the like count.
func (t *Track) Likes() int { return t.likes }

// Dislikes returns the dislike count.
func (t *Track) Dislikes() int { return t.dislikes }

// Listeners returns the current listener count.
func (t *Track) Listeners() int { return t.listeners }

// CreatedAt returns when the track was created.
func (t *Track) CreatedAt() time.Time { return t.createdAt }

// UpdatedAt returns when the track was last updated.
func (t *Track) UpdatedAt() time.Time { return t.updatedAt }
