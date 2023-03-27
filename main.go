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
var app_build string
var services []SystemDService

func init() {
	if len(app_version) == 0 {
		app_version = "dev"
	}
}

func main() {
	if len(os.Args) == 0 {
		log.Warn("Use hk-svcs-bridge serve to start the server", "version", app_version, "build", app_build)
		os.Exit(0)
	}

	log.Info("Starting Go HomeKit Services Bridge", "version", app_version, "build", app_build)
	bridge := SetupBridge()

	log.Debug("Loading the accessories...")
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

	log.Info("Updating HomeKit Pairing PIN", "code", conf.PairingCode)
	server.Pin = conf.PairingCode

	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	c := make(chan os.Signal, 1)
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
