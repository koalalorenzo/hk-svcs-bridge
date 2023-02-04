package main

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Name          string           `yaml:"name"`
	DatabnasePath string           `yaml:"db_path"`
	PairingCode   string           `yaml:"pairing_code"`
	UpdateDelay   int              `yaml:"update_delay_seconds"`
	Services      []SystemDService `yaml:"services"`
}

var conf Config

func init() {
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

	if conf.UpdateDelay <= 3 {
		conf.UpdateDelay = 3
	}

	if conf.DatabnasePath == "" {
		conf.DatabnasePath = "db"
	}
}
