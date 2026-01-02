package listener

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"hub/internal/infrastructure/icecast"
	"hub/internal/logger"
)

// Repository defines the listener repository interface.
type Repository interface {
	TrackListener(ctx context.Context, userID, trackID string) error
	GetUniqueListenerCount(ctx context.Context, trackID string) (int, error)
}

// TrackRepository defines the track repository interface for listener service.
type TrackRepository interface {
	ExistsByID(ctx context.Context, id string) (bool, error)
	UpdateListenerCount(ctx context.Context, trackID string, count int) error
}

// Service defines the listener service interface.
type Service interface {
	TrackCurrentListeners(ctx context.Context) error
}

type service struct {
	icecastClient icecast.Client
	listenerRepo  Repository
	trackRepo     TrackRepository
	logger        *logger.Logger
}

// NewService creates a new listener service.
func NewService(
	icecastClient icecast.Client,
	listenerRepo Repository,
	trackRepo TrackRepository,
	log *logger.Logger,
) Service {
	return &service{
		icecastClient: icecastClient,
		listenerRepo:  listenerRepo,
		trackRepo:     trackRepo,
		logger:        log,
	}
}

func (s *service) TrackCurrentListeners(ctx context.Context) error {
	log := s.logger.WithContext("listener", "track_current")

	stats, err := s.icecastClient.MountStats()
	if err != nil {
		log.WithError(err).Error("failed to get mount stats")
		return fmt.Errorf("failed to get mount stats: %w", err)
	}

	trackID := icecast.ExtractTrackID(stats.Title)
	if trackID == "" {
		log.Debug("no track ID in stream title")
		return nil
	}

	exists, err := s.trackRepo.ExistsByID(ctx, trackID)
	if err != nil {
		log.WithError(err).Error("failed to check if track exists")
		return fmt.Errorf("failed to check if track exists: %w", err)
	}
	if !exists {
		log.WithField("track_id", trackID).Debug("track not found, skipping")
		return nil
	}

	clientList, err := s.icecastClient.ListClients()
	if err != nil {
		log.WithError(err).Error("failed to get client list")
		return fmt.Errorf("failed to get client list: %w", err)
	}

	for _, l := range clientList.Listeners {
		userID := generateUserID(l.IP, l.UserAgent, l.ID)
		if err := s.listenerRepo.TrackListener(ctx, userID, trackID); err != nil {
			log.WithError(err).WithField("user_id", userID).Warn("failed to track listener")
		}
	}

	count, err := s.listenerRepo.GetUniqueListenerCount(ctx, trackID)
	if err != nil {
		log.WithError(err).Error("failed to get listener count")
		return fmt.Errorf("failed to get listener count: %w", err)
	}

	log.WithFields(map[string]interface{}{
		"track_id":       trackID,
		"listener_count": count,
	}).Debug("updated listener count")

	return s.trackRepo.UpdateListenerCount(ctx, trackID, count)
}

func generateUserID(ip, userAgent string, icecastID int) string {
	data := fmt.Sprintf("%s%s%d", ip, userAgent, icecastID)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
