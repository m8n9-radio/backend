package service

import (
	"context"

	"hub/internal/delivery/http/dto"
	"hub/internal/delivery/http/entity"
	"hub/internal/delivery/http/repository"
)

type (
	TrackService interface {
		Upsert(ctx context.Context, req *dto.CreateTrackRequest) (*entity.Track, error)
	}
	trackService struct {
		repo repository.TrackRepository
	}
)

func NewTrackService(repo repository.TrackRepository) TrackService {
	return &trackService{repo: repo}
}

func (s *trackService) Upsert(ctx context.Context, req *dto.CreateTrackRequest) (*entity.Track, error) {
	track := &entity.Track{
		ID:     req.Md5,
		Rotate: 1,
		Title:  req.StreamTitle,
		Cover:  req.StreamUrl,
	}

	return s.repo.Upsert(ctx, track)
}
