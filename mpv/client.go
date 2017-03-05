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
import "unsafe"

import (
	"github.com/go-kit/kit/log"
	box "github.com/pmdcosta/player"
	"sync"
)

// Client represents a client to the underlying video player.
type Client struct {
	// mpv handler
	handle *C.mpv_handle

	// player state
	running      bool
	runningMutex sync.Mutex

	// synchronize close
	mainloopExit chan struct{}

	// Player Service
	player Player

	// Events channel
	eventsChannel chan Event

	// logger
	logger log.Logger
}

// NewClient returns a new instance of Server.
func NewClient(log log.Logger) *Client {
	c := &Client{
		logger: log,
	}
	c.player.client = c
	c.player.playing = false
	return c
}

// Player returns the player service associated with the client.
func (mpv *Client) Player() box.Player { return &mpv.player }

// Open starts player.
func (mpv *Client) Open(flags map[string]bool, options map[string]string) error {
	// check if player is already running.
	if mpv.handle != nil || mpv.running {
		mpv.logger.Log("err", "mpv player already running")
		return box.ErrAlreadyRunning
	}

	// creates mpv player.
	mpv.running = true
	mpv.mainloopExit = make(chan struct{})
	mpv.handle = C.mpv_create()

	// set mpv startup flags.
	for k, v := range flags {
		err := mpv.setOptionFlag(k, v)
		if err != nil {
			return err
		}
	}

	// initializes mpv player.
	status := C.mpv_initialize(mpv.handle)
	if int(status) != 0 {
		mpv.logger.Log("err", "failed to initialize mpv player")
		return box.ErrMpvApiError
	}

	// set mpv startup options.
	for k, v := range options {
		err := mpv.setOptionString(k, v)
		if err != nil {
			return err
		}
	}

	// starts listening for events.
	mpv.eventsChannel = make(chan Event)
	go mpv.eventHandler()

	mpv.logger.Log("msg", "mpv player started")
	return nil
}

// Close terminates player.
func (mpv *Client) Close() error {
	// checks if the player is still running.
	mpv.runningMutex.Lock()
	if !mpv.running {
		mpv.logger.Log("err", "mpv player already terminated")
		return box.ErrAlreadyClosed
	}
	mpv.running = false
	mpv.runningMutex.Unlock()

	// waits until the main loop has exited.
	C.mpv_wakeup(mpv.handle)
	<-mpv.mainloopExit

	// terminates the mpv player.
	// blocks until the player has been fully brought down.
	handle := mpv.handle
	mpv.handle = nil
	C.mpv_terminate_destroy(handle)

	mpv.logger.Log("msg", "mpv player terminated")
	return nil
}

// setOptionFlag passes a boolean flag to mpv.
func (mpv *Client) setOptionFlag(key string, value bool) error {
	cValue := C.int(0)
	if value {
		cValue = 1
	}
	return mpv.setOption(key, C.MPV_FORMAT_FLAG, unsafe.Pointer(&cValue))
}

// setOptionInt passes an integer option to mpv.
func (mpv *Client) setOptionInt(key string, value int) error {
	cValue := C.int64_t(value)
	return mpv.setOption(key, C.MPV_FORMAT_INT64, unsafe.Pointer(&cValue))
}

// setOptionString passes a string option to mpv.
func (mpv *Client) setOptionString(key, value string) error {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	return mpv.setOption(key, C.MPV_FORMAT_STRING, unsafe.Pointer(&cValue))
}

// setOption is a generic function to pass options to mpv.
func (mpv *Client) setOption(key string, format C.mpv_format, value unsafe.Pointer) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	err := C.mpv_set_option(mpv.handle, cKey, format, value)
	if int(err) != 0 {
		mpv.logger.Log("err", "failed to set option", "status", C.GoString(C.mpv_error_string(err)))
		return box.ErrMpvApiError
	}
	return nil
}

// sendCommand sends a command to the libmpv player.
func (mpv *Client) sendCommand(command []string) error {
	cArray := C.makeCharArray(C.int(len(command) + 1))
	if cArray == nil {
		mpv.logger.Log("err", "failed to allocate memory")
		return box.ErrCalloc
	}
	defer C.free(unsafe.Pointer(cArray))

	for i, s := range command {
		cStr := C.CString(s)
		C.setArrayString(cArray, C.int(i), cStr)
		defer C.free(unsafe.Pointer(cStr))
	}

	err := C.mpv_command_async(mpv.handle, 0, cArray)
	if int(err) != 0 {
		mpv.logger.Log("err", "failed to send command", "status", C.GoString(C.mpv_error_string(err)))
		return box.ErrMpvApiError
	}
	return nil
}

// setProperty sets mpv property option.
func (mpv *Client) setProperty(name, value string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	err := C.mpv_set_property_async(mpv.handle, 1, cName, C.MPV_FORMAT_STRING, unsafe.Pointer(&cValue))
	if int(err) != 0 {
		mpv.logger.Log("err", "failed to set property", "status", C.GoString(C.mpv_error_string(err)))
		return box.ErrMpvApiError
	}
	return nil
}

// eventHandler waits for mpv events and sends to the player.
func (mpv *Client) eventHandler() {
	for {
		// wait until there is an event.
		// negative timeout means infinite timeout.
		event := C.mpv_wait_event(mpv.handle, -1)
		if event.error != 0 {
			mpv.logger.Log("err", "received error from mpv player")
			panic(event)
		}

		// terminate event channel.
		mpv.runningMutex.Lock()
		running := mpv.running
		mpv.runningMutex.Unlock()
		if !running {
			close(mpv.eventsChannel)
			mpv.mainloopExit <- struct{}{}
			return
		}

		// send event to player.
		mpv.eventsChannel <- Event(event.event_id)
	}
}
