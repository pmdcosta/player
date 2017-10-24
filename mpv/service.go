package mpv

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

// Service represents a service for interacting with the mpv media player.
type Service struct {
	client       *Client
	playing      bool
	playingMutex sync.Mutex
}

// SetCommand sets the command on the player.
func (s *Service) SetCommand(cmd ...string) error {
	return s.client.sendCommand(cmd)
}

// SetOption sets the options on the player.
func (s *Service) SetOption(key, value string) error {
	return s.client.setOptionString(key, value)
}

// SetFlag sets the option flag on the player.
func (s *Service) SetFlag(key string, value bool) error {
	return s.client.setOptionFlag(key, value)
}

// SetProperty sets the options on the player.
func (s *Service) SetProperty(key, value string) error {
	return s.client.setProperty(key, value)
}

// Stop closes the player client.
func (s *Service) Stop() error {
	return s.client.Close()
}

// Start opens the player client.
func (s *Service) Start() error {
	return s.client.Open()
}

// PlayFile starts playing the supplied file.
func (s *Service) PlayFile(path string) error {
	// check if file exists if it is not a remote video.
	if _, err := os.Stat(path); os.IsNotExist(err) && !strings.Contains(path, "http") {
		s.client.logger.WithFields(log.Fields{"path": path}).Info(ErrFileNotFound)
		return ErrFileNotFound
	}

	// start playing file.
	if err := s.SetProperty("pause", "no"); err != nil {
		return err
	}
	if err := s.SetCommand("loadfile", path); err != nil {
		return err
	}
	s.client.logger.WithFields(log.Fields{"video": path}).Info("playing video")
	return nil
}

// GetEvents handles mpv player events.
func (s *Service) getEvents() {
	for event := range s.client.events {
		switch event {
		case SetPropertyReply:
		case CommandReply:
		case StartFile:
			s.client.logger.WithFields(log.Fields{"event": "START_FILE"}).Debug("event received")
			s.playingMutex.Lock()
			s.playing = true
			s.playingMutex.Unlock()
		case EndFile:
			s.client.logger.WithFields(log.Fields{"event": "END_FILE"}).Debug("event received")
			s.playingMutex.Lock()
			s.playing = false
			s.playingMutex.Unlock()
			s.SetProperty("pause", "no")
		case Idle:
		default:
			continue
		}
	}
}
