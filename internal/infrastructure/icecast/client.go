package icecast

import (
	"fmt"
	"hub/internal/config"
	"io"
	"net/http"
)

type (
	Client interface {
		MountStats() (*ResponseSourceStats, error)
		ListClients() (*ResponseClientList, error)
	}

	client struct {
		host     string
		user     string
		password string
		mount    string
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

func (c *client) request(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.host, url), nil)
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

	return body, nil
}
