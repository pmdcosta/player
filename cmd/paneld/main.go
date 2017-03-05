package main

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	panel "github.com/pmdcosta/player"
	"github.com/pmdcosta/player/mpv"
	"net/http"
	"os"
)

var Player panel.Player
var logger log.Logger

func main() {
	// create logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// create client
	flags := map[string]bool{"no-resume-playback": true, "no-input-terminal": true, "quiet": true, "keep-open": true, "ytdl": true, "border": false, "loop": true}
	options := map[string]string{"hwdec": "vdpau", "screen": "2"}
	client := mpv.NewClient(logger)
	err := client.Open(flags, options)
	if err != nil {
		panic(err.Error())
	}

	// create player
	Player = client.Player()
	go Player.GetEvents()
	logger.Log("msg", "starting player")

	// routes.
	router := mux.NewRouter()
	router.HandleFunc("/play", Play).Methods("GET")
	router.HandleFunc("/pause", Pause).Methods("GET")
	router.HandleFunc("/resume", Resume).Methods("GET")
	router.HandleFunc("/seek/{value}", Seek).Methods("GET")
	router.HandleFunc("/next", Next).Methods("GET")
	router.HandleFunc("/pre", Previous).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("assets/"))))

	panic(http.ListenAndServe(":8000", router))

}

func Play(w http.ResponseWriter, req *http.Request) {
	err := Player.SetCommand([]string{"loadlist", "/home/pmdcosta/Videos/list"})
	if err != nil {
		logger.Log("err", err)
		return
	}
}

func Pause(w http.ResponseWriter, req *http.Request) {
	err := Player.SetProperty([]string{"pause", "yes"})
	if err != nil {
		logger.Log("err", err)
		return
	}
}

func Resume(w http.ResponseWriter, req *http.Request) {
	err := Player.SetProperty([]string{"pause", "no"})
	if err != nil {
		logger.Log("err", err)
		return
	}
}

func Seek(w http.ResponseWriter, req *http.Request) {
	value := mux.Vars(req)["value"]
	err := Player.SetCommand([]string{"seek", value})
	if err != nil {
		logger.Log("err", err)
		return
	}
}

func Next(w http.ResponseWriter, req *http.Request) {
	err := Player.SetCommand([]string{"playlist-next"})
	if err != nil {
		logger.Log("err", err)
		return
	}
}

func Previous(w http.ResponseWriter, req *http.Request) {
	err := Player.SetCommand([]string{"playlist-prev"})
	if err != nil {
		logger.Log("err", err)
		return
	}
}
