package cache

import (
	"context"
	"errors"
	"time"

	"hub/internal/domain/radio"
)

const (
	radioInfoKey = "radio:info"
	defaultTTL   = 10 * time.Second
)

// RadioCache provides caching for radio information.
type RadioCache struct {
	cache      Cache
	repository radio.Repository
	ttl        time.Duration
}

// NewRadioCache creates a new RadioCache.
func NewRadioCache(cache Cache, repository radio.Repository, ttl time.Duration) *RadioCache {
	if ttl == 0 {
		ttl = defaultTTL
	}
	return &RadioCache{
		cache:      cache,
		repository: repository,
		ttl:        ttl,
	}
}

// radioInfoDTO is used for JSON serialization.
type radioInfoDTO struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	StreamURL    string `json:"stream_url"`
	Listeners    int    `json:"listeners"`
	ListenerPeak int    `json:"listener_peak"`
}

// GetCurrentInfo returns cached radio info or fetches from repository.
func (c *RadioCache) GetCurrentInfo(ctx context.Context) (*radio.RadioInfo, error) {
	var dto radioInfoDTO

	err := c.cache.Get(ctx, radioInfoKey, &dto)
	if err == nil {
		info := radio.NewRadioInfo(dto.Name, dto.Description, dto.StreamURL, dto.Listeners, dto.ListenerPeak)
		return &info, nil
	}

	if !errors.Is(err, ErrCacheMiss) {
		return nil, err
	}

	info, err := c.repository.GetCurrentInfo(ctx)
	if err != nil {
		return nil, err
	}

	dto = radioInfoDTO{
		Name:         info.Name(),
		Description:  info.Description(),
		StreamURL:    info.StreamURL(),
		Listeners:    info.Listeners(),
		ListenerPeak: info.ListenerPeak(),
	}

	_ = c.cache.Set(ctx, radioInfoKey, dto, c.ttl)

	return info, nil
}

// Invalidate clears the cached radio info.
func (c *RadioCache) Invalidate(ctx context.Context) error {
	return c.cache.Delete(ctx, radioInfoKey)
}
