package statistics

import (
	"context"
)

// TrackStats represents track statistics.
type TrackStats struct {
	Title     string
	Cover     string
	Rotate    int
	Likes     int
	Dislikes  int
	Listeners int
}

// Category represents a statistics category.
type Category struct {
	Key         string
	Description string
	Icon        string
	Tracks      []*TrackStats
}

// Repository defines the statistics repository interface.
type Repository interface {
	GetHistory(ctx context.Context) ([]*TrackStats, error)
	GetTopListened(ctx context.Context) ([]*TrackStats, error)
	GetTopRotate(ctx context.Context) ([]*TrackStats, error)
	GetTopLikes(ctx context.Context) ([]*TrackStats, error)
	GetTopDislikes(ctx context.Context) ([]*TrackStats, error)
}

// Service defines the statistics service interface.
type Service interface {
	GetStatistics(ctx context.Context) ([]*Category, error)
}

type service struct {
	repo Repository
}

// NewService creates a new statistics service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetStatistics(ctx context.Context) ([]*Category, error) {
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

	return []*Category{
		{Description: "Last V tracks", Key: "history", Icon: "HistoryIcon", Tracks: history},
		{Description: "Top V listend tracks", Key: "listen", Icon: "ListenIcon", Tracks: topListened},
		{Description: "Top V rotated tracks", Key: "rotate", Icon: "RotateIcon", Tracks: topRotate},
		{Description: "Top V liked tracks", Key: "likes", Icon: "LikeIcon", Tracks: topLikes},
		{Description: "Top V disliked tracks", Key: "dislikes", Icon: "DislikeIcon", Tracks: topDislikes},
	}, nil
}
