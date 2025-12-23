package icecast

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"hub/internal/config"
)

type (
	Client interface {
		GetMountStats() (*SourceStats, error)
	}
	client struct {
		host     string
		user     string
		password string
		mount    string
	}

	// Icecast XML response structure
	IcecastStats struct {
		XMLName xml.Name `xml:"icestats"`
		Sources []Source `xml:"source"`
	}

	Source struct {
		Mount           string `xml:"mount,attr"`
		ServerName      string `xml:"server_name"`
		ServerDesc      string `xml:"server_description"`
		ServerURL       string `xml:"server_url"`
		Listeners       int    `xml:"listeners"`
		ListenerPeak    int    `xml:"listener_peak"`
		Genre           string `xml:"genre"`
		Bitrate         string `xml:"bitrate"`
		Title           string `xml:"title"`
		MetadataURL     string `xml:"metadata_url"`
		AudioInfo       string `xml:"audio_info"`
	}

	SourceStats struct {
		Name         string
		Description  string
		StreamURL    string
		Listeners    int
		ListenerPeak int
		Genre        string
		Bitrate      string
		Title        string
		MetadataURL  string
		AudioInfo    string
	}
)

func NewClient(cfg config.Config) (Client, error) {
	host, user, password, mount := cfg.IcecastConnection()

	return &client{
		host:     host,
		user:     user,
		password: password,
		mount:    mount,
	}, nil
}

func (c *client) GetMountStats() (*SourceStats, error) {
	// Make HTTP request to Icecast admin stats
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

	// Parse XML
	var stats IcecastStats
	if err := xml.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	// Find the mount we're looking for
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
