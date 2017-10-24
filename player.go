package player

type Player interface {
	SetCommand(cmd ...string) error
	SetOption(key, value string) error
	SetProperty(key, value string) error
	PlayFile(path string) error
	SetFlag(key string, value bool) error
	Start() error
	Stop() error
}

type Subscriber interface {
	OnMessage(topic string, payload []byte) error
}
