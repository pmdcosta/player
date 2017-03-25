package player

// mpv errors.
const (
	ErrMpvRunning   = Error("player is already running")
	ErrMpvClosed    = Error("player is already closed")
	ErrMpvSetOption = Error("failed to set player option")
	ErrMpvInit      = Error("failed to initialize player")
	ErrMpvCalloc    = Error("failed to allocate memory")
	ErrMpvPlayer    = Error("received mpv player error")
)

// Error represents a Box error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
