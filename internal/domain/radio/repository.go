package radio

import "context"

// Repository defines the radio repository interface.
type Repository interface {
	// GetCurrentInfo returns the current radio stream information.
	GetCurrentInfo(ctx context.Context) (*RadioInfo, error)
}
