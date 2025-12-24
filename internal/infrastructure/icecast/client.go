package icecast

import (
	"hub/internal/config"
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
