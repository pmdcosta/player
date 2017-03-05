package player

type Client interface {
	Open() error
	Close() error
}

type Player interface {
	SetCommand(cmd []string) error
	SetOption(value []string) error
	SetProperty(value []string) error
	GetEvents()
}

type Subscriber interface {
	OnMessage(topic string, payload []byte) error
}
