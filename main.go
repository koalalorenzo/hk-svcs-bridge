package main

import (
	"context"
	log "golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

var app_version string

func init() {
	if len(app_version) == 0 {
		app_version = "dev"
	}
}

func main() {
	log.Info("Setting up the bridge")
	bridge := SetupBridge()

	log.Info("Loading the accessories")
	svcsA := []*accessory.A{}
	services = []SystemDService{}
	for _, svc := range conf.Services {
		services = append(services, svc.Init())
		svcsA = append(svcsA, svc.Accessory.A)
	}

	log.Debug("Setting up the server...")

	// Create the hap server.
	fs := hap.NewFsStore(conf.DatabasePath)
	server, err := hap.NewServer(fs, bridge.A, svcsA...)
	if err != nil {
		// stop if an error happens
		log.Error("Error setting HomeKit Server", "error", err)
	}

	log.Info("Using HomeKit Pairing Pin Code", "code", conf.PairingCode)
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

	// Run the update Ticker
	t := StartSystemDCheckTicker()
	defer t.Stop()

	// Run the server.
	log.Info("Ready")
	server.ListenAndServe(ctx)
}
