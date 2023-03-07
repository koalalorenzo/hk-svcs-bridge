package main

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/creasty/defaults"
)

type Config struct {
	// Name is the Name of the Bridge, visible in the HomeKit app
	Name string `yaml:"name" default:"-"`
	// DatabasePath is the path where the files will be stored
	DatabasePath string `yaml:"db_path" default:"./db"`
	// PairingCode consist into a customizable pairing code for your accessory
	PairingCode string `yaml:"pairing_code" default:"10042001"`
	// UpdateFrequency is the frequency in seconds to check for SystemD service
	UpdateFrequency int `yaml:"update_frequency" default:"3"`
	// Services is the list of SystemD Services to add as accessories
	Services []SystemDService `yaml:"services"`
}

func (c *Config) SetDefaults() {
	if len(c.PairingCode) < 8 {
		c.PairingCode = "10042001"
	}

	if defaults.CanUpdate(c.Name) {
		hn, err := os.Hostname()
		if err != nil {
			log.Printf("Unable to get hostname: #%v", err)
			conf.Name = "GoHomeKitBridge"
		} else {
			conf.Name = hn
		}
	}
}

var conf Config

func init() {
	if app_version == "" {
		app_version = "local-dev"
	}

	path, a := os.LookupEnv("CONFIG")
	if !a {
		path = "config.yaml"
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading %s: #%v ", path, err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	if err := defaults.Set(&conf); err != nil {
		log.Fatalf("Error seting Defaults: %v", err)
	}
}
