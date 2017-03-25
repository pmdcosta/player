package http

import (
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/julienschmidt/httprouter"
	"github.com/pmdcosta/player"
	"net/http"
)

// PlayerHandler represents an HTTP API handler for the player.
type PlayerHandler struct {
	Logger log.Logger
	*httprouter.Router

	// injected services.
	playerService player.PlayerService

	// handler path.
	Path string
}

// PlayerPath is the default path to handler.
const PlayerPath = "/api/player/"

// NewPlayerHandler returns a new instance of PlayerHandler.
func NewPlayerHandler(ps player.PlayerService, logger log.Logger) *PlayerHandler {
	// create handler.
	h := &PlayerHandler{
		Router:        httprouter.New(),
		playerService: ps,
		Logger:        logger,
		Path:          PlayerPath,
	}

	// start listening.
	h.POST(h.Path+"play", h.handlePostPlay)
	return h
}

// handlePostPlay handles requests to play a video.
func (h *PlayerHandler) handlePostPlay(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode request.
	var req postPlayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, player.ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	// start video user.
	switch err := h.playerService.SetCommand("loadfile", req.Video); err {
	case nil:
		encodeJSON(w, &postPlayResponse{}, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type postPlayRequest struct {
	Video string `json:"video"`
}

type postPlayResponse struct {
	Err string `json:"err,omitempty"`
}
