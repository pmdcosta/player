package mpv_test

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pmdcosta/player/mpv"
)

// Client is a test wrapper.
type Client struct {
	*mpv.Client
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	// create client wrapper.
	c := &Client{
		Client: mpv.NewClient(),
	}
	return c
}

// MustOpenClient returns an new, open instance of Client.
func MustOpenClient() *Client {
	log.SetLevel(log.DebugLevel)
	c := NewClient()
	if err := c.Client.Open(); err != nil {
		panic(err)
	}
	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	return c.Client.Close()
}
