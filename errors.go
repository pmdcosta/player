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

// HTTP errors.
const (
	ErrHTTPFailed = Error("failed to start server")
	ErrHTTPCert   = Error("failed to load certificates")
	ErrHTTPAPI    = Error("unknown API error")
)

// General errors.
const (
	ErrInvalidJSON  = Error("invalid json")
	ErrUnauthorized = Error("unauthorized")
	ErrInternal     = Error("internal error")
)

// Error represents a Box error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
