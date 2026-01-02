package track

import (
	"hub/internal/domain/shared"
)

// Domain errors for track operations.
var (
	ErrTrackNotFound = shared.NewDomainError(
		shared.ErrNotFound,
		"track not found",
	)
)
