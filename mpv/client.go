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
	"github.com/pmdcosta/player"
	"sync"
)

// Client represents a client to the underlying video player.
type Client struct {
	logger log.Logger

	// mpv handler.
	handle *C.mpv_handle

	// player state.
	running      bool
	runningMutex sync.Mutex

	// synchronize closing the client.
	mainloopExit chan struct{}

	// player service.
	player Mpv

	// events channel.
	eventsChannel chan Event
}

// NewClient returns a new instance of Client.
func NewClient(log log.Logger) *Client {
	c := &Client{
		logger: log,
	}
	c.player.client = c
	c.player.playing = false
	return c
}

// Mpv returns the player service associated with the client.
func (mpv *Client) Mpv() player.Mpv { return &mpv.player }

// Open starts the player.
func (mpv *Client) Open(flags map[string]bool, options map[string]string) error {
	// check if player is already running.
	if mpv.handle != nil || mpv.running {
		mpv.logger.Log("err", player.ErrMpvRunning)
		return player.ErrMpvRunning
	}

	// creates mpv player.
	mpv.runningMutex.Lock()
	mpv.running = true
	mpv.runningMutex.Unlock()
	mpv.mainloopExit = make(chan struct{})
	mpv.handle = C.mpv_create()

	// set mpv startup flags.
	for k, v := range flags {
		if err := mpv.setOptionFlag(k, v); err {
			mpv.logger.Log("err", player.ErrMpvSetOption, "msg", err.Error())
			return err
		}
	}

	// initializes mpv player.
	if int(C.mpv_initialize(mpv.handle)) != 0 {
		mpv.logger.Log("err", player.ErrMpvInit)
		return player.ErrMpvInit
	}

	// set mpv startup options.
	for k, v := range options {
		if err := mpv.setOptionString(k, v); err {
			mpv.logger.Log("err", player.ErrMpvSetOption, "msg", err.Error())
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
		mpv.logger.Log("err", player.ErrMpvClosed)
		return player.ErrMpvClosed
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
		mpv.logger.Log("err", player.ErrMpvSetOption, "status", C.GoString(C.mpv_error_string(err)))
		return player.ErrMpvSetOption
	}
	return nil
}

// sendCommand sends a command to the libmpv player.
func (mpv *Client) sendCommand(command []string) error {
	cArray := C.makeCharArray(C.int(len(command) + 1))
	if cArray == nil {
		mpv.logger.Log("err", player.ErrMpvCalloc)
		return player.ErrMpvCalloc
	}
	defer C.free(unsafe.Pointer(cArray))

	for i, s := range command {
		cStr := C.CString(s)
		C.setArrayString(cArray, C.int(i), cStr)
		defer C.free(unsafe.Pointer(cStr))
	}

	err := C.mpv_command_async(mpv.handle, 0, cArray)
	if int(err) != 0 {
		mpv.logger.Log("err", player.ErrMpvSetOption, "status", C.GoString(C.mpv_error_string(err)))
		return player.ErrMpvSetOption
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
		mpv.logger.Log("err", player.ErrMpvSetOption, "status", C.GoString(C.mpv_error_string(err)))
		return player.ErrMpvSetOption
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
			mpv.logger.Log("err", player.ErrMpvPlayer, "event", event)
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
