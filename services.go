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
	// SystemDPeriodicCheck if true will periodically check the status of th
	// SystemD ServiceName.
	SystemDPeriodicCheck bool `yaml:"systemd_check"`

	// Accessory is t he HAP accessory
	Accessory *accessory.Switch `yaml:"-"`
	// Updating is used to prevent the CheckStatus to interfere with the cmds
	IsUpdating bool `yaml:"-"`
}

func (s *SystemDService) runCmd(cmd string) {
	if s.IsUpdating {
		return
	}
	s.IsUpdating = true
	defer func() { s.IsUpdating = false }()

	log.Printf("Running %s", s.OffCommand)
	run := exec.Command(strings.Split(s.OffCommand, " ")...)
	out, err := run.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			log.Printf("Error encountered %v:\n %v", exitCode, out)
		}
	}
	s.Accessory.Switch.On.SetValue(false)

	// Prevent Updating during the to avoid overlapping
	time.Sleep(time.Duration(500) * time.Millisecond)
}

func (s *SystemDService) Init() {
	sw := accessory.NewSwitch(accessory.Info{
		Name: s.ServiceName,
	})

	// We assume that the service is already running
	sw.Switch.On.SetValue(true)

	// Adds event for on-off switching
	sw.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			s.Off()
		} else {
			s.On()
		}
	})

	s.Accessory = sw
}
