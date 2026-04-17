package main

import (
	"github.com/pankaj28843/confluence-cli/internal/client"
)

func newClient() (*client.Client, error) {
	cfg, err := client.FromEnv()
	if err != nil {
		return nil, newConfigError(err)
	}
	cfg.Debug = debug
	c, err := client.New(cfg)
	if err != nil {
		return nil, newConfigError(err)
	}
	return c, nil
}
