package main

import (
	"fmt"
	"log"

	"github.com/brutella/hap/accessory"
)

type SystemDService struct {
	Name        string            `yaml:"name"`
	ServiceName string            `yaml:"service_name"`
	OnCommand   string            `yaml:"oncmd"`
	OffCommand  string            `yaml:"offcmd"`
	Accessory   *accessory.Switch `yaml:"-"`
}

func (s *SystemDService) Off() {
	cmd := fmt.Sprintf("services %s stop", s.ServiceName)
	if s.OffCommand != "" {
		cmd = s.OffCommand
	}

	log.Printf("Running %s", cmd)
}

func (s *SystemDService) On() {
	cmd := fmt.Sprintf("services %s start", s.ServiceName)
	if s.OnCommand != "" {
		cmd = s.OnCommand
	}

	log.Printf("Running %s", cmd)
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
