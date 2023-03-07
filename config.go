package main

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	// Name is the Name of the Bridge, visible in the HomeKit app
	Name string `yaml:"name"`
	// DatabasePath is the path where the files will be stored
	DatabasePath string `yaml:"db_path"`
	// PairingCode consist into a customizable pairing code for your accessory
	PairingCode string `yaml:"pairing_code"`
	// UpdateFrequency is the frequency in seconds to check for SystemD service
	UpdateFrequency int `yaml:"update_frequency"`
	// Services is the list of SystemD Services to add as accessories
	Services []SystemDService `yaml:"services"`
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

	log.Printf("%v", conf)

	// Set Defaults
	if len(conf.PairingCode) < 8 {
		conf.PairingCode = "10042001"
	}

	if len(conf.Name) <= 3 {
		hn, err := os.Hostname()
		if err != nil {
			log.Printf("Unable to get hostname: #%v", err)
			conf.Name = "SystemD"
		} else {
			conf.Name = hn
		}
	}

	if conf.UpdateFrequency <= 3 {
		conf.UpdateFrequency = 3
	}

	if conf.DatabasePath == "" {
		conf.DatabasePath = "db"
	}
}
