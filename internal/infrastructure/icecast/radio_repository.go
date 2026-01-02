package icecast

import (
	"context"

	"hub/internal/domain/radio"
)

// RadioRepository implements radio.Repository using Icecast.
type RadioRepository struct {
	client Client
}

// NewRadioRepository creates a new RadioRepository.
func NewRadioRepository(client Client) *RadioRepository {
	return &RadioRepository{client: client}
}

// Ensure RadioRepository implements radio.Repository
var _ radio.Repository = (*RadioRepository)(nil)

// GetCurrentInfo returns the current radio stream information.
func (r *RadioRepository) GetCurrentInfo(ctx context.Context) (*radio.RadioInfo, error) {
	stats, err := r.client.MountStats()
	if err != nil {
		return nil, err
	}

	info := radio.NewRadioInfo(
		stats.Name,
		stats.Description,
		stats.StreamURL,
		stats.Listeners,
		stats.ListenerPeak,
	)

	return &info, nil
}
