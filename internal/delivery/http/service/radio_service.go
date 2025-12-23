package service

import (
	"context"

	"hub/internal/config"
	"hub/internal/delivery/http/entity"
)

type RadioService interface {
	GetRadioInfo(ctx context.Context) (*entity.RadioEntity, error)
}

type radioService struct {
	cfg           config.Config
	icecastClient IcecastClient
}

func NewRadioService(cfg config.Config, icecastClient IcecastClient) RadioService {
	return &radioService{
		cfg:           cfg,
		icecastClient: icecastClient,
	}
}

func (s *radioService) GetRadioInfo(ctx context.Context) (*entity.RadioEntity, error) {
	source, err := s.icecastClient.FetchStatus(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.RadioEntity{
		Name:        source.ServerName,
		Description: source.ServerDescription,
		StreamUrl:   s.cfg.IcecastStreamURL(),
		Listener: entity.ListenerEntity{
			Current: source.Listeners,
			Peak:    source.ListenerPeak,
		},
	}, nil
}
