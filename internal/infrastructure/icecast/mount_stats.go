package icecast

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

type (
	icecastStats struct {
		XMLName xml.Name                 `xml:"icestats"`
		Sources []icecastMountSourceStat `xml:"source"`
	}

	icecastMountSourceStat struct {
		Mount        string `xml:"mount,attr"`
		ServerName   string `xml:"server_name"`
		ServerDesc   string `xml:"server_description"`
		ServerURL    string `xml:"server_url"`
		Listeners    int    `xml:"listeners"`
		ListenerPeak int    `xml:"listener_peak"`
		Genre        string `xml:"genre"`
		Bitrate      string `xml:"bitrate"`
		Title        string `xml:"title"`
		StreamStart  string `xml:"stream_start"`
		MetadataURL  string `xml:"metadata_url"`
		AudioInfo    string `xml:"audio_info"`
	}

	ResponseSourceStats struct {
		Name         string
		Description  string
		StreamURL    string
		Listeners    int
		ListenerPeak int
		Genre        string
		Bitrate      string
		Title        string
		StreamStart  string
		MetadataURL  string
		AudioInfo    string
	}
)

func (c *client) MountStats() (*ResponseSourceStats, error) {
	body, err := c.request("/admin/stats")
	if err != nil {
		return nil, err
	}

	var stats icecastStats
	if err := xml.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	for _, source := range stats.Sources {
		if source.Mount == c.mount {
			return &ResponseSourceStats{
				Name:         source.ServerName,
				Description:  source.ServerDesc,
				StreamURL:    source.ServerURL,
				Listeners:    source.Listeners,
				ListenerPeak: source.ListenerPeak,
				Genre:        source.Genre,
				Bitrate:      source.Bitrate,
				Title:        source.Title,
				StreamStart:  source.StreamStart,
				MetadataURL:  source.MetadataURL,
				AudioInfo:    source.AudioInfo,
			}, nil
		}
	}

	return nil, fmt.Errorf("mount %s not found", c.mount)
}

// ExtractTrackID extracts MD5 track ID from title format: "Artist - Title [MD5]"
func ExtractTrackID(title string) string {
	re := regexp.MustCompile(`\[([a-f0-9]{32})\]$`)
	matches := re.FindStringSubmatch(title)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
