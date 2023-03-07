package main

import (
	"fmt"
	"log"
	"time"
)

func StartSystemDCheckTicker() (ticker *time.Ticker) {
	freq := time.Second * time.Duration(conf.UpdateFrequency)
	ticker = time.NewTicker(freq)

	go func() {
		for range ticker.C {
			log.Print("Checking")
			for _, s := range conf.Services {
				if s.PeriodicCheck {
					go s.CheckStatus()
				}
			}
		}
		fmt.Printf("Ticker Stopped")
	}()
	return
}
