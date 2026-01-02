package reaction

import (
	"hub/internal/domain/shared"
)

// Domain errors for reaction operations.
var (
	ErrReactionExists = shared.NewDomainError(
		shared.ErrAlreadyExists,
		"user has already reacted to this track",
	)

	ErrReactionNotFound = shared.NewDomainError(
		shared.ErrNotFound,
		"reaction not found",
	)

	ErrTrackNotFound = shared.NewDomainError(
		shared.ErrNotFound,
		"cannot react to non-existent track",
	)
)
