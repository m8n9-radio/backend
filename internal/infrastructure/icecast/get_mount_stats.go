package icecast

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func (c *client) GetMountStats() (*SourceStats, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/stats", c.host), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(c.user, c.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Can't send request to Icecast API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Icecast API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var stats IcecastStats
	if err := xml.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	for _, source := range stats.Sources {
		if source.Mount == c.mount {
			return &SourceStats{
				Name:         source.ServerName,
				Description:  source.ServerDesc,
				StreamURL:    source.ServerURL,
				Listeners:    source.Listeners,
				ListenerPeak: source.ListenerPeak,
				Genre:        source.Genre,
				Bitrate:      source.Bitrate,
				Title:        source.Title,
				MetadataURL:  source.MetadataURL,
				AudioInfo:    source.AudioInfo,
			}, nil
		}
	}

	return nil, fmt.Errorf("mount %s not found", c.mount)
}
