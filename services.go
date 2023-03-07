package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/brutella/hap/accessory"
)

type SystemDService struct {
	// Name visible in the Homekit app
	Name string `yaml:"name"`
	// ServiceName is the SystemD Service name
	ServiceName string `yaml:"service_name"`
	// OnCommand is useful if we want to customize what command to run when
	// the trigger is ON
	OnCommand string `yaml:"oncmd"`
	// OffCommand is like OnCommand but runs when the trigger is set to Off
	OffCommand string `yaml:"offcmd"`
	// PeriodicCheck if true will periodically check the status of the
	// ServiceName.by running systemctl or the custom command
	PeriodicCheck bool `yaml:"periodic_check"`
	// PeriodicCheckCmd is the command that if returns 0 will set to
	PeriodicCheckCmd string `yaml:"periodic_check_cmd"`

	// Accessory is t he HAP accessory
	Accessory *accessory.Switch `yaml:"-"`
	// Updating is used to prevent the CheckStatus to interfere with the cmds
	IsUpdating bool `yaml:"-"`
}

func (s *SystemDService) runCmd(cmd string, succSetVal, failSetVal bool) {
	if s.IsUpdating {
		return
	}
	s.IsUpdating = true
	defer func() { s.IsUpdating = false }()

	log.Printf("Running %s", cmd)

	args := strings.Split(cmd, " ")
	run := exec.Command(args[0], args[1:]...)

	out, err := run.CombinedOutput()
	if err != nil {
		s.Accessory.Switch.On.SetValue(failSetVal)

		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			log.Printf("Cmd returned %v:\n %v", exitCode, out)
		} else {
			log.Printf("Error running command: %v", err)
		}

		return
	}

	s.Accessory.Switch.On.SetValue(succSetVal)

	// Prevent Updating during the to avoid overlapping
	time.Sleep(time.Duration(250) * time.Millisecond)
}

func (s *SystemDService) CheckStatus() {
	if s.IsUpdating {
		return
	}
	s.IsUpdating = true
	defer func() { s.IsUpdating = false }()

	s.runCmd(s.PeriodicCheckCmd, true, false)
}

func (s *SystemDService) Init() {
	sw := accessory.NewSwitch(accessory.Info{
		Name: s.Name,
	})

	// Set Defaults
	if s.ServiceName == "" {
		s.ServiceName = s.Name
	}

	if s.OffCommand == "" {
		s.OffCommand = fmt.Sprintf("systemctl stop %s", s.ServiceName)
	}

	if s.OnCommand == "" {
		s.OnCommand = fmt.Sprintf("systemctl start %s", s.ServiceName)
	}

	if s.PeriodicCheckCmd == "" {
		s.PeriodicCheckCmd = fmt.Sprintf("systemctl is-active %s", s.ServiceName)
	}

	// We assume that the service is already running
	sw.Switch.On.SetValue(true)

	// Adds event for on-off switching
	sw.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			s.runCmd(s.OffCommand, false, false)
		} else {
			s.runCmd(s.OnCommand, true, false)
		}
	})

	s.Accessory = sw
}
