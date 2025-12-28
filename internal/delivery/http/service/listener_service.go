package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"hub/internal/delivery/http/repository"
	"hub/internal/infrastructure/icecast"
)

type (
	ListenerService interface {
		TrackCurrentListeners(ctx context.Context) error
	}

	listenerService struct {
		icecastClient icecast.Client
		listenerRepo  repository.ListenerRepository
		trackRepo     repository.TrackRepository
	}
)

func NewListenerService(
	icecastClient icecast.Client,
	listenerRepo repository.ListenerRepository,
	trackRepo repository.TrackRepository,
) ListenerService {
	return &listenerService{
		icecastClient: icecastClient,
		listenerRepo:  listenerRepo,
		trackRepo:     trackRepo,
	}
}

// TrackCurrentListeners tracks all active listeners for the current track
func (s *listenerService) TrackCurrentListeners(ctx context.Context) error {
	// 1. Get current track from Icecast metadata
	stats, err := s.icecastClient.MountStats()
	if err != nil {
		return fmt.Errorf("failed to get mount stats: %w", err)
	}

	// 2. Extract track ID from title (format: "Artist - Title [MD5]")
	trackID := icecast.ExtractTrackID(stats.Title)
	if trackID == "" {
		// No track ID in metadata, skip tracking
		return nil
	}

	// 2.5. Verify track exists in database (skip if not)
	exists, err := s.trackRepo.ExistsByID(ctx, trackID)
	if err != nil {
		return fmt.Errorf("failed to check if track exists: %w", err)
	}
	if !exists {
		// Track doesn't exist yet, skip listener tracking
		// Tracks should be created by another service (e.g., from API or liquidsoap)
		return nil
	}

	// 3. Get all active listeners
	clientList, err := s.icecastClient.ListClients()
	if err != nil {
		return fmt.Errorf("failed to get client list: %w", err)
	}

	// Debug: Log how many listeners Icecast reported
	fmt.Printf("Icecast reports %d listeners, got %d listener records\n",
		clientList.Count, len(clientList.Listeners))

	// 4. Track each listener
	for _, listener := range clientList.Listeners {
		userID := generateUserID(listener.IP, listener.UserAgent, listener.ID)

		// Debug: Log each listener being tracked
		fmt.Printf("Tracking listener: ID=%d IP=%s UserAgent=%s => user_id=%s\n",
			listener.ID, listener.IP, listener.UserAgent, userID)

		// UPSERT: Insert new or increment listen_count
		err := s.listenerRepo.TrackListener(ctx, userID, trackID)
		if err != nil {
			// Log error but continue with other listeners
			fmt.Printf("failed to track listener %s: %v\n", userID, err)
			continue
		}
	}

	// 5. Update listener count in tracks table
	count, err := s.listenerRepo.GetUniqueListenerCount(ctx, trackID)
	if err != nil {
		return fmt.Errorf("failed to get listener count: %w", err)
	}

	err = s.trackRepo.UpdateListenerCount(ctx, trackID, count)
	if err != nil {
		return fmt.Errorf("failed to update track listener count: %w", err)
	}

	return nil
}

// generateUserID creates a unique user ID from IP, UserAgent, and Icecast client ID
func generateUserID(ip, userAgent string, icecastID int) string {
	data := fmt.Sprintf("%s%s%d", ip, userAgent, icecastID)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
