package reaction

import (
	"hub/internal/domain/shared"
)

// ErrInvalidReactionType is returned when a reaction type is invalid.
var ErrInvalidReactionType = shared.NewDomainError(
	shared.ErrInvalidInput,
	"reaction type must be 'like' or 'dislike'",
)

// ReactionType is a value object representing the type of reaction.
type ReactionType struct {
	value string
}

// Predefined reaction types.
var (
	Like    = ReactionType{value: "like"}
	Dislike = ReactionType{value: "dislike"}
)

// NewReactionType creates a new ReactionType from a string.
// Returns an error if the string is not "like" or "dislike".
func NewReactionType(value string) (ReactionType, error) {
	switch value {
	case "like":
		return Like, nil
	case "dislike":
		return Dislike, nil
	default:
		return ReactionType{}, ErrInvalidReactionType
	}
}

// String returns the string representation of the ReactionType.
func (rt ReactionType) String() string {
	return rt.value
}

// IsLike returns true if the reaction is a like.
func (rt ReactionType) IsLike() bool {
	return rt.value == "like"
}

// IsDislike returns true if the reaction is a dislike.
func (rt ReactionType) IsDislike() bool {
	return rt.value == "dislike"
}

// Equals checks if two ReactionTypes are equal.
func (rt ReactionType) Equals(other ReactionType) bool {
	return rt.value == other.value
}
