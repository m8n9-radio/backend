package radio

import (
	"context"
	"errors"
	"fmt"

	"hub/internal/infrastructure/icecast"
)

var (
	ErrInvalidResponse = errors.New("invalid response from icecast server")
	ErrNoActiveStream  = errors.New("no active stream available")
)

// RadioInfo represents radio stream information.
type RadioInfo struct {
	Name         string
	Description  string
	StreamUrl    string
	Listeners    int
	ListenerPeak int
}

// ListenerInfo represents listener count information.
type ListenerInfo struct {
	Current int
	Peak    int
}

// Service defines the radio service interface.
type Service interface {
	GetRadioInfo(ctx context.Context) (*RadioInfo, error)
	GetListeners(ctx context.Context) (*ListenerInfo, error)
}

type service struct {
	icecastClient icecast.Client
}

// NewService creates a new radio service.
func NewService(icecastClient icecast.Client) Service {
	return &service{icecastClient: icecastClient}
}

func (s *service) GetRadioInfo(ctx context.Context) (*RadioInfo, error) {
	source, err := s.icecastClient.MountStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get icecast stats: %w", err)
	}

	return &RadioInfo{
		Name:         source.Name,
		Description:  source.Description,
		StreamUrl:    source.StreamURL,
		Listeners:    source.Listeners,
		ListenerPeak: source.ListenerPeak,
	}, nil
}

func (s *service) GetListeners(ctx context.Context) (*ListenerInfo, error) {
	source, err := s.icecastClient.MountStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get icecast stats: %w", err)
	}

	return &ListenerInfo{
		Current: source.Listeners,
		Peak:    source.ListenerPeak,
	}, nil
}
