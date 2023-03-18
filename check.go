package main

import (
	log "golang.org/x/exp/slog"
	"time"
)

func StartSystemDCheckTicker() (ticker *time.Ticker) {
	freq := time.Second * time.Duration(conf.UpdateFrequency)
	ticker = time.NewTicker(freq)

	go func() {
		for range ticker.C {
			log.Debug("Checking")
			for _, s := range conf.Services {
				if s.PeriodicCheck {
					go s.CheckStatus()
				}
			}
		}
		log.Warn("Ticker Stopped")
	}()
	return
}
