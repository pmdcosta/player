package player

// Client is a generic client interface for managing packages.
type Client interface {
	Open() error
	Close() error
}

// PlayerService represents a service for managing the media player.
type PlayerService interface {
	SetCommand(cmd ...string) error
	SetOption(key, value string) error
	SetProperty(key, value string) error
	GetEvents()
}
