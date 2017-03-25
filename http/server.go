package http

import (
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/player"
	"net"
	"net/http"
	"net/url"
)

// Server represents an HTTP server.
type Server struct {
	Logger log.Logger
	ln     net.Listener

	// handler to serve.
	Handler *Handler

	// server host address.
	Host string
}

// NewServer returns a new instance of Server.
func NewServer(l log.Logger, h *Handler, host string) *Server {
	return &Server{
		Logger:  l,
		Handler: h,
		Host:    host,
	}
}

// Open opens a socket and serves the http server.
func (s *Server) Open() error {
	// open socket.
	ln, err := net.Listen("tcp", s.Host)
	if err != nil {
		s.Logger.Log("err", player.ErrHTTPFailed, "msg", err.Error())
		return err
	}
	s.ln = ln

	// start http server.
	go func() { http.Serve(s.ln, s.Handler) }()

	s.Logger.Log("msg", "server started", "host", s.Host)
	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

// Port returns the port that the server is open on.
func (s *Server) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}

// Client represents a client to connect to the http server.
type Client struct {
	URL           url.URL
	playerService PlayerService
}

// NewClient returns a new instance of client.
func NewClient() *Client {
	c := &Client{}
	c.playerService.URL = &c.URL
	return c
}

// PlayerService returns the service for managing the player.
func (c *Client) PlayerService() player.PlayerService {
	return &c.playerService
}
