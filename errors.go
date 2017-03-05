package player

// mpv errors.
const (
	ErrAlreadyRunning = Error("already running")
	ErrAlreadyClosed  = Error("already closed")
	ErrMpvApiError    = Error("mpv api error")
	ErrCalloc         = Error("failed to allocate memory")
)

// Error represents a Box error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
