package mpv

import (
	"github.com/pmdcosta/player"
)

// MPV errors.
const (
	ErrAlreadyRunning = player.Error("player is already running")
	ErrNotRunning     = player.Error("player is not running")
	ErrInitPlayer     = player.Error("failed to initiate mpv player")
	ErrSetOption      = player.Error("failed to set mpv option")
	ErrSetCommand     = player.Error("failed to set mpv command")
	ErrSetProperty    = player.Error("failed to set mpv property")
	ErrMalloc         = player.Error("failed to allocate memory")
	ErrEventError     = player.Error("invalid mpv event received")
	ErrFileNotFound   = player.Error("file not found")
)
