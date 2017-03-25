package mpv

// Event mpv event types.
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

// Mpv represents a service for managing mpv.
type Mpv struct {
	client  *Client
	playing bool
}

// SetCommand sets the command on the player.
func (mpv *Mpv) SetCommand(cmd ...string) error {
	mpv.client.logger.Log("msg", "setCommand", "cmd", cmd[0], "value", cmd[1])
	return mpv.client.sendCommand(cmd)
}

// SetOption sets the options on the player.
func (mpv *Mpv) SetOption(key, value string) error {
	mpv.client.logger.Log("msg", "setOption", "key", key, "value", value)
	return mpv.client.setOptionString(key, value)
}

// SetProperty sets the property on the player.
func (mpv *Mpv) SetProperty(key, value string) error {
	mpv.client.logger.Log("msg", "setProperty", "key", key, "value", value)
	return mpv.client.setProperty(key, value)
}

// GetEvents handles mpv player events.
func (mpv *Mpv) GetEvents() {
	for event := range mpv.client.eventsChannel {
		switch event {
		case MpvEventSetPropertyReply:
		case MpvEventCommandReply:
		case MpvEventStartFile:
			mpv.client.logger.Log("msg", "mpv player event received", "event", "START_FILE")
			mpv.playing = true
		case MpvEventEndFile:
			mpv.client.logger.Log("msg", "mpv player event received", "event", "END_FILE")
			mpv.playing = false
		case MpvEventIdle:
		default:
			continue
		}
	}
}
