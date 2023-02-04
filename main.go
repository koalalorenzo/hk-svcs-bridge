package main

import (
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"

	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Print("Setting up the bridge")
	bridge := SetupBridge()

	log.Print("Loading the accessories")
	svcsA := []*accessory.A{}
	for _, svc := range conf.Services {
		svc.Init()
		svcsA = append(svcsA, svc.Accessory.A)
	}

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore("./db")

	log.Print("Setting up the server")

	// Create the hap server.
	server, err := hap.NewServer(fs, bridge.A, svcsA...)
	if err != nil {
		// stop if an error happens
		log.Panic(err)
	}

	log.Printf("Using HomeKit Pairing Pin Code: %s", conf.PairingCode)
	server.Pin = conf.PairingCode

	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals.
		signal.Stop(c)
		// Cancel the context to stop the server.
		cancel()
	}()

	// Run the server.
	log.Print("READY")
	server.ListenAndServe(ctx)
}
