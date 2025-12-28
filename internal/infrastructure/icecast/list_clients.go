package icecast

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type (
	icecastClientListStats struct {
		XMLName xml.Name                  `xml:"icestats"`
		Sources []icecastClientListSource `xml:"source"`
	}
	icecastClientListSource struct {
		Mount     string     `xml:"mount,attr"`
		Listeners int        `xml:"listeners"`
		Listener  []listener `xml:"listener"`
	}
	listener struct {
		ID        int    `xml:"ID"`
		IP        string `xml:"IP"`
		UserAgent string `xml:"UserAgent"`
		Lag       int    `xml:"lag"`
		Connected int    `xml:"Connected"`
	}

	ResponseClientList struct {
		Count     int
		Listeners []ResponseListener
	}

	ResponseListener struct {
		ID        int
		IP        string
		UserAgent string
		Connected int
	}
)

func (c *client) ListClients() (*ResponseClientList, error) {
	byte, err := c.request(fmt.Sprintf("/admin/listclients?mount=%s", c.mount))
	if err != nil {
		return nil, err
	}

	var stats icecastClientListStats
	if err := xml.Unmarshal(byte, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	for _, source := range stats.Sources {
		if source.Mount == c.mount {

			resp := &ResponseClientList{
				Count:     source.Listeners,
				Listeners: make([]ResponseListener, 0, len(source.Listener)),
			}

			for _, l := range source.Listener {
				resp.Listeners = append(resp.Listeners, ResponseListener{
					ID:        l.ID,
					IP:        strings.TrimSpace(l.IP),
					UserAgent: strings.TrimSpace(l.UserAgent),
					Connected: l.Connected,
				})
			}

			return resp, nil
		}
	}

	return nil, fmt.Errorf("mount %s not found", c.mount)
}
