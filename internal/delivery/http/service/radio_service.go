package service

import (
	"context"
	"errors"
	"fmt"

	"hub/internal/delivery/http/entity"
	"hub/internal/infrastructure/icecast"
)

var (
	ErrInvalidResponse = errors.New("invalid response from icecast server")
	ErrNoActiveStream  = errors.New("no active stream available")
)

type (
	RadioService interface {
		GetRadioInfo(ctx context.Context) (*entity.RadioEntity, error)
	}
	radioService struct {
		icecastClient icecast.Client
	}
)

func NewRadioService(icecastClient icecast.Client) RadioService {
	return &radioService{
		icecastClient: icecastClient,
	}
}

func (s *radioService) GetRadioInfo(ctx context.Context) (*entity.RadioEntity, error) {
	source, err := s.icecastClient.GetMountStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get icecast stats: %w", err)
	}

	return &entity.RadioEntity{
		Name:        source.Name,
		Description: source.Description,
		StreamUrl:   source.StreamURL,
		Listener: entity.ListenerEntity{
			Current: source.Listeners,
			Peak:    source.ListenerPeak,
		},
	}, nil
}
