package reaction

import (
	"time"

	"hub/internal/domain/track"
)

// Reaction represents a user's reaction to a track.
type Reaction struct {
	id           string
	userID       UserID
	trackID      track.TrackID
	reactionType ReactionType
	createdAt    time.Time
}

// NewReaction creates a new Reaction entity.
func NewReaction(userID UserID, trackID track.TrackID, reactionType ReactionType) *Reaction {
	return &Reaction{
		userID:       userID,
		trackID:      trackID,
		reactionType: reactionType,
		createdAt:    time.Now(),
	}
}

// ReconstructReaction rebuilds a Reaction from persistence data.
func ReconstructReaction(
	id string,
	userID, trackID, reactionType string,
	createdAt time.Time,
) (*Reaction, error) {
	uid, err := NewUserID(userID)
	if err != nil {
		return nil, err
	}

	tid, err := track.NewTrackID(trackID)
	if err != nil {
		return nil, err
	}

	rt, err := NewReactionType(reactionType)
	if err != nil {
		return nil, err
	}

	return &Reaction{
		id:           id,
		userID:       uid,
		trackID:      tid,
		reactionType: rt,
		createdAt:    createdAt,
	}, nil
}

// Getters

// ID returns the reaction's database ID.
func (r *Reaction) ID() string { return r.id }

// UserID returns the user's ID.
func (r *Reaction) UserID() UserID { return r.userID }

// TrackID returns the track's ID.
func (r *Reaction) TrackID() track.TrackID { return r.trackID }

// ReactionType returns the type of reaction.
func (r *Reaction) ReactionType() ReactionType { return r.reactionType }

// CreatedAt returns when the reaction was created.
func (r *Reaction) CreatedAt() time.Time { return r.createdAt }

// IsLike returns true if this is a like reaction.
func (r *Reaction) IsLike() bool { return r.reactionType.IsLike() }

// IsDislike returns true if this is a dislike reaction.
func (r *Reaction) IsDislike() bool { return r.reactionType.IsDislike() }
