package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/pmdcosta/player"
	"github.com/pmdcosta/player/mpv"
	"net/http"
)

var Player player.Player
var logger log.Logger
var playlist string

func main() {
	// create player client and service.
	playerClient := mpv.NewClient()
	Player = playerClient.Service()

	// player flags.
	flags := map[string]bool{"no-resume-playback": true, "no-input-terminal": true, "quiet": true, "keep-open": true, "ytdl": true, "border": false, "loop": true, "fullscreen": true}
	for f, v := range flags {
		Player.SetFlag(f, v)
	}

	// player options.
	options := map[string]string{"hwdec": "vdpau"}
	for f, v := range options {
		Player.SetOption(f, v)
	}

	// start player.
	if err := playerClient.Open(); err != nil {
		panic(err)
	}

	var (
		assets = flag.String("assets", "./assets/", "Choose assets location.")
		pl     = flag.String("playlist", "./test/playlist", "Choose playlist location.")
	)
	flag.Parse()
	playlist = *pl

	// lazy routing and handling...
	router := mux.NewRouter()
	router.HandleFunc("/play", Play).Methods("GET")
	router.HandleFunc("/pause", Pause).Methods("GET")
	router.HandleFunc("/resume", Resume).Methods("GET")
	router.HandleFunc("/seek/{value}", Seek).Methods("GET")
	router.HandleFunc("/next", Next).Methods("GET")
	router.HandleFunc("/pre", Previous).Methods("GET")
	router.HandleFunc("/osd/{value}", Osd).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(*assets))))
	panic(http.ListenAndServe(":8000", router))

}

func Play(w http.ResponseWriter, req *http.Request) {
	if err := Player.SetCommand("loadlist", playlist); err != nil {
		logger.Error(err)
	}
}

func Pause(w http.ResponseWriter, req *http.Request) {
	if err := Player.SetProperty("pause", "yes"); err != nil {
		logger.Error(err)
	}
}

func Resume(w http.ResponseWriter, req *http.Request) {
	if err := Player.SetProperty("pause", "no"); err != nil {
		logger.Error(err)
	}
}

func Seek(w http.ResponseWriter, req *http.Request) {
	if err := Player.SetCommand("seek", mux.Vars(req)["value"]); err != nil {
		logger.Error(err)
	}
}

func Next(w http.ResponseWriter, req *http.Request) {
	if err := Player.SetCommand("playlist-next"); err != nil {
		logger.Error(err)
	}
}

func Previous(w http.ResponseWriter, req *http.Request) {
	if err := Player.SetCommand("playlist-prev"); err != nil {
		logger.Error(err)
	}
}

func Osd(w http.ResponseWriter, req *http.Request) {
	if err := Player.SetCommand("osd", mux.Vars(req)["value"]); err != nil {
		logger.Error(err)
	}
}
