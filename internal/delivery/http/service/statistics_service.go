package service

import (
	"context"

	"hub/internal/delivery/http/dto"
	"hub/internal/delivery/http/repository"
)

type (
	StatisticsService interface {
		GetStatistics(ctx context.Context) (*dto.StatisticsResponse, error)
	}
	statisticsService struct {
		repo repository.StatisticsRepository
	}
)

func NewStatisticsService(repo repository.StatisticsRepository) StatisticsService {
	return &statisticsService{repo: repo}
}

func (s *statisticsService) GetStatistics(ctx context.Context) (*dto.StatisticsResponse, error) {
	history, err := s.repo.GetHistory(ctx)
	if err != nil {
		return nil, err
	}

	topListened, err := s.repo.GetTopListened(ctx)
	if err != nil {
		return nil, err
	}

	topRotate, err := s.repo.GetTopRotate(ctx)
	if err != nil {
		return nil, err
	}

	topLikes, err := s.repo.GetTopLikes(ctx)
	if err != nil {
		return nil, err
	}

	topDislikes, err := s.repo.GetTopDislikes(ctx)
	if err != nil {
		return nil, err
	}

	response := &dto.StatisticsResponse{
		Statistics: []*dto.StatisticCategory{
			{
				Key:         "history",
				Description: "",
				Icon:        "HistoryIcon",
				Tracks:      history,
			},
			{
				Key:         "listen",
				Description: "",
				Icon:        "ListenIcon",
				Tracks:      topListened,
			},
			{
				Key:         "rotate",
				Description: "",
				Icon:        "RotateIcon",
				Tracks:      topRotate,
			},
			{
				Key:         "likes",
				Description: "",
				Icon:        "LikeIcon",
				Tracks:      topLikes,
			},
			{
				Key:         "dislikes",
				Description: "",
				Icon:        "DislikeIcon",
				Tracks:      topDislikes,
			},
		},
	}

	return response, nil
}
