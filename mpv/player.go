package mpv

// Event ps event types.
type Event int

const (
	// MpvEventSetPropertyReply is a reply to a set_property request.
	MpvEventSetPropertyReply Event = 4
	// MpvEventCommandReply is areply to a command request.
	MpvEventCommandReply Event = 5
	// MpvEventStartFile is a notification before playback start of a file.
	MpvEventStartFile Event = 6
	// MpvEventEndFile is a notification after playback end.
	MpvEventEndFile Event = 7
	// MpvEventIdle signifies an idle mode was entered.
	MpvEventIdle Event = 11
)

// PlayerService represents a service for managing ps.
type PlayerService struct {
	client  *Client
	playing bool
}

// SetCommand sets the command on the player.
func (ps *PlayerService) SetCommand(cmd ...string) error {
	ps.client.logger.Log("msg", "setCommand", "cmd", cmd[0], "value", cmd[1])
	return ps.client.sendCommand(cmd)
}

// SetOption sets the options on the player.
func (ps *PlayerService) SetOption(key, value string) error {
	ps.client.logger.Log("msg", "setOption", "key", key, "value", value)
	return ps.client.setOptionString(key, value)
}

// SetProperty sets the property on the player.
func (ps *PlayerService) SetProperty(key, value string) error {
	ps.client.logger.Log("msg", "setProperty", "key", key, "value", value)
	return ps.client.setProperty(key, value)
}

// GetEvents handles ps player events.
func (ps *PlayerService) GetEvents() {
	for event := range ps.client.eventsChannel {
		switch event {
		case MpvEventSetPropertyReply:
		case MpvEventCommandReply:
		case MpvEventStartFile:
			ps.client.logger.Log("msg", "player event received", "event", "START_FILE")
			ps.playing = true
		case MpvEventEndFile:
			ps.client.logger.Log("msg", "player event received", "event", "END_FILE")
			ps.playing = false
		case MpvEventIdle:
		default:
			continue
		}
	}
}
