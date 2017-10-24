package mpv

// #include <mpv/client.h>
// #include <stdlib.h>
// #cgo LDFLAGS: -lmpv
//
// /* some helper functions for string arrays */
// char** makeCharArray(int size) {
//     return calloc(sizeof(char*), size);
// }
// void setArrayString(char** a, int i, char* s) {
//     a[i] = s;
// }
import "C"
import (
	"unsafe"
)

/*
	Interacting with MPV (the video player) requires the libmpv lib.
	Documentation @ https://mpv.io/manual/stable/
*/

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pmdcosta/player"
	"sync"
	"time"
)

// Client represents a client to for playing videos.
type Client struct {
	// package logger.
	logger *log.Entry

	// underlying mpv client.
	handle *C.mpv_handle

	// video player state.
	running bool
	lock    sync.Mutex

	// mpv events stream.
	events chan Event

	// gracefully shutdown.
	quitChan chan struct{}

	// service for interacting with the video player.
	service Service
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	c := &Client{
		logger:   log.WithFields(log.Fields{"package": "player"}),
		events:   make(chan Event, 5),
		quitChan: make(chan struct{}),
	}
	c.service.client = c

	// start listening for player events.
	go c.service.getEvents()
	return c
}

// Open starts client handler.
func (c *Client) Open() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// check if player is already running.
	if c.handle != nil || c.running {
		c.logger.Error(ErrAlreadyRunning)
		return ErrAlreadyRunning
	}
	c.running = true

	// creates the mpv player.
	c.handle = C.mpv_create()

	// set mpv configuration flag.
	if err := c.setOptionFlag("config", true); err != nil {
		return err
	}

	// start the mpv player.
	status := C.mpv_initialize(c.handle)
	if int(status) != 0 {
		c.logger.WithFields(log.Fields{"status": status}).Error(ErrInitPlayer)
		return ErrInitPlayer
	}

	// starts listening for mpv events.
	go c.eventHandler()

	c.logger.Info("mpv player started")
	return nil
}

// Close terminates player.
func (c *Client) Close() error {
	// checks if the player is still running.
	c.lock.Lock()
	if !c.running {
		c.logger.Info(ErrNotRunning)
		c.lock.Unlock()
		return ErrNotRunning
	}
	c.running = false
	c.lock.Unlock()

	// waits until the main loop has exited.
	C.mpv_wakeup(c.handle)
	select {
	case _, ok := <-c.quitChan:
		if !ok {
			c.quitChan = nil
		}
	}

	// terminates the mpv player and waits until the player has been fully brought down.
	C.mpv_terminate_destroy(c.handle)
	c.handle = nil

	time.Sleep(1 * time.Second)
	c.logger.Info("mpv player stopped")
	return nil
}

// setOptionFlag passes a boolean flag to mpv.
func (c *Client) setOptionFlag(key string, value bool) error {
	cValue := C.int(0)
	if value {
		cValue = 1
	}
	c.logger.WithFields(log.Fields{"key": key, "value": value}).Debug("setting option")
	return c.setOption(key, C.MPV_FORMAT_FLAG, unsafe.Pointer(&cValue))
}

// setOptionInt passes an integer option to mpv.
func (c *Client) setOptionInt(key string, value int) error {
	cValue := C.int64_t(value)
	c.logger.WithFields(log.Fields{"key": key, "value": value}).Debug("setting option")
	return c.setOption(key, C.MPV_FORMAT_INT64, unsafe.Pointer(&cValue))
}

// setOptionString passes a string option to mpv.
func (c *Client) setOptionString(key, value string) error {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	c.logger.WithFields(log.Fields{"key": key, "value": value}).Debug("setting option")
	return c.setOption(key, C.MPV_FORMAT_STRING, unsafe.Pointer(&cValue))
}

// setOption is a generic function to pass options to mpv.
func (c *Client) setOption(key string, format C.mpv_format, value unsafe.Pointer) error {
	if !c.running {
		c.logger.WithFields(log.Fields{"err": ErrNotRunning}).Debug(ErrSetOption)
		return ErrNotRunning
	}

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	if err := C.mpv_set_option(c.handle, cKey, format, value); int(err) != 0 {
		c.logger.WithFields(log.Fields{"key": key, "status": C.GoString(C.mpv_error_string(err))}).Error(ErrSetOption)
		return ErrSetOption
	}
	return nil
}

// sendCommand sends a command to the libmpv player.
func (c *Client) sendCommand(command []string) error {
	if !c.running {
		c.logger.WithFields(log.Fields{"err": ErrNotRunning}).Debug(ErrSetCommand)
		return ErrNotRunning
	}

	cArray := C.makeCharArray(C.int(len(command) + 1))
	if cArray == nil {
		c.logger.WithFields(log.Fields{"commands": command}).Error(ErrMalloc)
		return ErrMalloc
	}
	defer C.free(unsafe.Pointer(cArray))

	for i, s := range command {
		cStr := C.CString(s)
		C.setArrayString(cArray, C.int(i), cStr)
		defer C.free(unsafe.Pointer(cStr))
	}

	if err := C.mpv_command_async(c.handle, 0, cArray); int(err) != 0 {
		c.logger.WithFields(log.Fields{"commands": command, "status": C.GoString(C.mpv_error_string(err))}).Error(ErrSetCommand)
		return ErrSetCommand
	}

	c.logger.WithFields(log.Fields{"command": command}).Debug("setting command")
	return nil
}

// setProperty sets mpv property option.
func (c *Client) setProperty(name, value string) error {
	if !c.running {
		c.logger.WithFields(log.Fields{"err": ErrNotRunning}).Debug(ErrSetProperty)
		return ErrNotRunning
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	if err := C.mpv_set_property_async(c.handle, 1, cName, C.MPV_FORMAT_STRING, unsafe.Pointer(&cValue)); int(err) != 0 {
		c.logger.WithFields(log.Fields{"key": name, "value": value, "status": C.GoString(C.mpv_error_string(err))}).Error(ErrSetProperty)
		return ErrSetProperty
	}
	c.logger.WithFields(log.Fields{"key": name, "value": value}).Debug("setting property")
	return nil
}

// eventHandler waits for mpv events and sends them to the player.
func (c *Client) eventHandler() {
	for {
		// wait for an mpv event.
		event := C.mpv_wait_event(c.handle, -1)
		if event.error != 0 {
			c.logger.WithFields(log.Fields{"err": event.error}).Error(ErrEventError)
			continue
		}

		// check if player still running and close otherwise.
		c.lock.Lock()
		running := c.running
		c.lock.Unlock()

		if !running {
			c.quitChan <- struct{}{}
			return
		}

		// send event to player.
		c.events <- Event(event.event_id)
	}
}

// Event mpv player event types.
type Event int

// Event types
const (
	SetPropertyReply Event = 4  // Reply to a set_property request.
	CommandReply           = 5  // Reply to a command request.
	StartFile              = 6  // Notification before playback start of a file.
	EndFile                = 7  // Notification after playback end.
	Idle                   = 11 // Idle mode was entered.
)

// Service returns the service used to interact with the video player.
func (c *Client) Service() player.Player { return &c.service }
