package http

import (
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/player"
	"net/http"
	"strings"
)

// Handler is a collection of all the service handlers.
type Handler struct {
	PlayerHandler *PlayerHandler
}

// ServeHTTP delegates a request to the appropriate sub-handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// route request.
	if strings.HasPrefix(r.URL.Path, h.PlayerHandler.Path) {
		h.PlayerHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// Error writes an API error message to the response and logger.
func Error(w http.ResponseWriter, err error, code int, logger log.Logger) {
	// log error.
	logger.Log("err", player.ErrHTTPAPI, "code", code, "msg", err.Error())

	// hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = player.ErrInternal
	}

	// write generic error response.
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}

// errorResponse is a generic response for sending a error.
type errorResponse struct {
	Err string `json:"err,omitempty"`
}

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, logger log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}
