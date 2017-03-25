package player

// Client is a generic client interface for managing packages.
type Client interface {
	Open() error
	Close() error
}

// Mpv represents a service for managing the media player.
type Mpv interface {
	SetCommand(cmd ...string) error
	SetOption(key, value string) error
	SetProperty(key, value string) error
	GetEvents()
}
