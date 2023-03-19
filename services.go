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
	PeriodicCheck bool `yaml:"periodic_check" default:"true"`
	// PeriodicCheckCmd is the command that if returns 0 will set to
	PeriodicCheckCmd string `yaml:"periodic_check_cmd"`

	// Accessory is t he HAP accessory
	Accessory *accessory.Switch `yaml:"-" default:"-"`
	// Updating is used to prevent the CheckStatus to interfere with the cmds
	IsUpdating bool `yaml:"-" default:"-"`
}

func (s *SystemDService) runCmd(cmd string, succSetVal, failSetVal bool) {
	if s.IsUpdating {
		return
	}
	s.IsUpdating = true
	defer func() { s.IsUpdating = false }()
	log := log.With("cmd", cmd)

	log.Debug("Running...")
	args := strings.Split(cmd, " ")
	run := exec.Command(args[0], args[1:]...)

	out, err := run.CombinedOutput()
	log = log.With("output", out)
	if err != nil {
		s.Accessory.Switch.On.SetValue(failSetVal)
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			log.Warn("Error running output", "exitCode", exitCode)
		} else {
			log.Warn("Error", "error", err)
		}

		return
	}

	s.Accessory.Switch.On.SetValue(succSetVal)

	// Prevent Updating during the to avoid overlapping
	time.Sleep(time.Duration(250) * time.Millisecond)
}

func (s *SystemDService) CheckStatus() {
	log := log.With("svc_name", s.Name)
	if !s.PeriodicCheck {
		log.Debug("Skipping periodic check")
		return
	}

	log.Debug("Checking status for service")
	s.runCmd(s.PeriodicCheckCmd, true, false)
}

func (s *SystemDService) SetDefaults() {
	log.Debug("Loading default config for svc", "svc", s)

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

	log.Debug("Done loading default config for svc", "svc", s)
}

func (s *SystemDService) Init() SystemDService {
	s.IsUpdating = false

	sw := accessory.NewSwitch(accessory.Info{
		Name: s.Name,
	})

	// We assume that the service is already running
	sw.Switch.On.SetValue(true)

	// Adds event for on-off switching
	sw.Switch.On.OnValueRemoteUpdate(func(on bool) {
		// If it fails running any cmd, don't change the value of the switch
		oldValue := s.Accessory.Switch.On.Value()
		if on == true {
			s.runCmd(s.OffCommand, false, oldValue)
		} else {
			s.runCmd(s.OnCommand, true, oldValue)
		}
	})

	s.Accessory = sw
	return *s
}
