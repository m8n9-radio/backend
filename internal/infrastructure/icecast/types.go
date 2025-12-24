package icecast

import "encoding/xml"

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

	IcecastStats struct {
		XMLName xml.Name `xml:"icestats"`
		Sources []Source `xml:"source"`
	}

	Source struct {
		Mount        string `xml:"mount,attr"`
		ServerName   string `xml:"server_name"`
		ServerDesc   string `xml:"server_description"`
		ServerURL    string `xml:"server_url"`
		Listeners    int    `xml:"listeners"`
		ListenerPeak int    `xml:"listener_peak"`
		Genre        string `xml:"genre"`
		Bitrate      string `xml:"bitrate"`
		Title        string `xml:"title"`
		MetadataURL  string `xml:"metadata_url"`
		AudioInfo    string `xml:"audio_info"`
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
