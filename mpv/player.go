package mpv

// Event mpv event types.
type Event int

const (
	MPV_EVENT_SET_PROPERTY_REPLY Event = 4  // Reply to a set_property request
	MPV_EVENT_COMMAND_REPLY      Event = 5  // Reply to a command request
	MPV_EVENT_START_FILE         Event = 6  // Notification before playback start of a file
	MPV_EVENT_END_FILE           Event = 7  // Notification after playback end
	MPV_EVENT_IDLE               Event = 11 // Idle mode was entered.
)

// Player represents a service for managing mpv.
type Player struct {
	client  *Client
	playing bool
}

// SetCommand sets the command on the player.
func (mpv *Player) SetCommand(cmd []string) error {
	mpv.client.logger.Log("msg", "setCommand", "cmd", cmd[0])
	return mpv.client.sendCommand(cmd)
}

// SetOptions sets the options on the player.
func (mpv *Player) SetOption(values []string) error {
	mpv.client.logger.Log("msg", "setOption", "key", values[0], "value", values[1])
	return mpv.client.setOptionString(values[0], values[1])
}

// SetProperty sets the options on the player.
func (mpv *Player) SetProperty(values []string) error {
	mpv.client.logger.Log("msg", "setOption", "key", values[0], "value", values[1])
	return mpv.client.setProperty(values[0], values[1])
}

// GetEvents handles mpv player events.
func (mpv *Player) GetEvents() {
	for event := range mpv.client.eventsChannel {
		switch event {
		case MPV_EVENT_SET_PROPERTY_REPLY:
		case MPV_EVENT_COMMAND_REPLY:
		case MPV_EVENT_START_FILE:
			mpv.client.logger.Log("msg", "mpv player event received", "event", "START_FILE")
			mpv.playing = true
		case MPV_EVENT_END_FILE:
			mpv.client.logger.Log("msg", "mpv player event received", "event", "END_FILE")
			mpv.playing = false
		case MPV_EVENT_IDLE:
		default:
			continue
		}
	}
}
