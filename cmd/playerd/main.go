package main

import (
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/player/http"
	"github.com/pmdcosta/player/mpv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// create logger.
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// gracefully shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// create client.
	flags := map[string]bool{"no-resume-playback": true, "keep-open": true, "ytdl": true, "video": false}
	options := map[string]string{}
	client := mpv.NewClient(logger)
	err := client.Open(flags, options)
	if err != nil {
		panic(err.Error())
	}

	// create player
	ps := client.PlayerService()
	go ps.GetEvents()
	logger.Log("msg", "starting player")

	// start http server.
	s := http.NewServer(logger, &http.Handler{PlayerHandler: http.NewPlayerHandler(ps, logger)}, ":3000")
	if err := s.Open(); err != nil {
		panic(err)
	}
	defer s.Close()

	// wait for user signal.
	<-sigs
	logger.Log("msg", "shutting down")
}
