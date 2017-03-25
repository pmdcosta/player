package http

import (
	"github.com/pmdcosta/player"
	"net/url"
)

// Ensure service implements interface.
var _ player.PlayerService = &PlayerService{}

// PlayerService represents an HTTP implementation of player.PlayerService.
type PlayerService struct {
	URL *url.URL
}

func (s *PlayerService) SetCommand(cmd ...string) error {
	return nil
}

func (s *PlayerService) SetOption(key, value string) error {
	return nil
}

func (s *PlayerService) SetProperty(key, value string) error {
	return nil
}

func (s *PlayerService) GetEvents() {
}
