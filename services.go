package main

import (
	"fmt"
	log "golang.org/x/exp/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/brutella/hap/accessory"
	"github.com/creasty/defaults"
)

type SystemDService struct {
	// Name visible in the Homekit app
	Name string `yaml:"name"`
	// ServiceName is the SystemD Service name
	ServiceName string `yaml:"service_name"`
	// OnCommand is useful if we want to customize what command to run when
	// the trigger is ON
	OnCommand string `yaml:"on_cmd"`
	// OffCommand is like OnCommand but runs when the trigger is set to Off
	OffCommand string `yaml:"off_cmd"`
	// PeriodicCheck if true will periodically check the status of the
	// ServiceName.by running systemctl or the custom command
	PeriodicCheck bool `yaml:"-" default:"-"`
	// PeriodicCheckCmd is the command that if returns 0 will set to
	PeriodicCheckCmd string `yaml:"check_cmd"`

	// Accessory is t he HAP accessory
	Accessory *accessory.Switch `yaml:"-" default:"-"`
	// Updating is used to prevent the CheckStatus to interfere with the cmds
	IsUpdating bool `yaml:"-" default:"-"`
}

// SetState will change the state only if needed.
func (s *SystemDService) SetState(newState bool) {
	oldState := s.Accessory.Switch.On.Value()
	log := log.With("oldState", oldState, "svcName", s.Name)

	if oldState != newState {
		log.Debug("Changing state", "newState", newState, "oldState", oldState)
		s.Accessory.Switch.On.SetValue(newState)
	}
}

func (s *SystemDService) runCmd(cmd string, succSetVal, failSetVal bool) {
	if s.IsUpdating {
		return
	}
	s.IsUpdating = true
	defer func() { s.IsUpdating = false }()
	log := log.With("cmd", cmd, "svcName", s.Name)

	log.Debug("Running...")
	args := strings.Split(cmd, " ")
	run := exec.Command(args[0], args[1:]...)

	out, err := run.CombinedOutput()
	log = log.With("output", out)
	if err != nil {
		// There was an error, let's change to failure
		s.SetState(failSetVal)
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			log.Warn("Error running output", "exitCode", exitCode)
		} else {
			log.Warn("Error", "error", err)
		}

		return
	}

	// Success! Let's change the state
	s.SetState(succSetVal)

	// Prevent Updating during the to avoid overlapping
	time.Sleep(time.Duration(500) * time.Millisecond)
}

func (s *SystemDService) CheckStatus() {
	log := log.With("svcName", s.Name)
	if !s.PeriodicCheck {
		log.Debug("Skipping periodic check")
		return
	}

	log.Debug("Checking status for service")
	s.runCmd(s.PeriodicCheckCmd, true, false)
}

func (s *SystemDService) SetDefaults() {
	if defaults.CanUpdate(s.ServiceName) {
		s.ServiceName = s.Name
	}

	if defaults.CanUpdate(s.OffCommand) {
		s.OffCommand = fmt.Sprintf("systemctl stop %s", s.ServiceName)
	}

	if defaults.CanUpdate(s.OnCommand) {
		s.OnCommand = fmt.Sprintf("systemctl start %s", s.ServiceName)
	}

	if defaults.CanUpdate(s.PeriodicCheckCmd) {
		s.PeriodicCheckCmd = fmt.Sprintf("systemctl is-active %s", s.ServiceName)
	}

	// Disable the periodic check if the cmd is "Disabled"
	// this is a workaround due to setting default with falsy values
	// TODO: Set the boolean from the YAML instead of force custom cmd
	s.PeriodicCheck = s.PeriodicCheckCmd != "disabled"
}

func (s *SystemDService) Init() SystemDService {
	s.IsUpdating = false

	sw := accessory.NewSwitch(accessory.Info{
		Name: s.Name,
	})

	// We assume that the service is already running
	// sw.Switch.On.SetValue(true)

	// Adds event for on-off switching
	sw.Switch.On.OnValueRemoteUpdate(func(newState bool) {
		if newState == false {
			s.runCmd(s.OffCommand, false, true)
		} else {
			s.runCmd(s.OnCommand, true, false)
		}
	})

	s.Accessory = sw
	return *s
}
