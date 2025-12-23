package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"hub/internal/config"
)

var (
	ErrIcecastUnavailable = errors.New("icecast server unavailable")
	ErrInvalidResponse    = errors.New("invalid response from icecast")
	ErrNoActiveStream     = errors.New("no active stream found")
)

type IcecastSource struct {
	ServerName        string `json:"server_name"`
	ServerDescription string `json:"server_description"`
	Listeners         int    `json:"listeners"`
	ListenerPeak      int    `json:"listener_peak"`
}

type IcecastStats struct {
	Source json.RawMessage `json:"source"`
}

type IcecastStatus struct {
	Icestats IcecastStats `json:"icestats"`
}

type IcecastClient interface {
	FetchStatus(ctx context.Context) (*IcecastSource, error)
}

type icecastClient struct {
	cfg        config.Config
	httpClient *http.Client
}

func NewIcecastClient(cfg config.Config) IcecastClient {
	return &icecastClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *icecastClient) FetchStatus(ctx context.Context) (*IcecastSource, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.cfg.IcecastStatusURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIcecastUnavailable, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIcecastUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code %d", ErrIcecastUnavailable, resp.StatusCode)
	}

	var status IcecastStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidResponse, err)
	}

	return c.parseSource(status.Icestats.Source)
}

func (c *icecastClient) parseSource(raw json.RawMessage) (*IcecastSource, error) {
	if len(raw) == 0 {
		return nil, ErrNoActiveStream
	}

	// Try parsing as single source object
	var source IcecastSource
	if err := json.Unmarshal(raw, &source); err == nil {
		if source.ServerName != "" {
			return &source, nil
		}
	}

	// Try parsing as array of sources
	var sources []IcecastSource
	if err := json.Unmarshal(raw, &sources); err == nil {
		if len(sources) > 0 {
			return &sources[0], nil
		}
	}

	return nil, ErrNoActiveStream
}
