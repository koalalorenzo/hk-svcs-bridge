package main

import (
	"io/ioutil"
	"os"

	log "golang.org/x/exp/slog"

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
			log.Warn("Unable to get hostname", "error", err)
			conf.Name = "GoHomeKitBridge"
		} else {
			conf.Name = hn
		}
	}
}

var conf Config

func init() {
	// Set Log Level
	var logLevel = new(log.LevelVar)
	ll, _ := os.LookupEnv("LOG_LEVEL")
	switch ll {
	case "debug":
		logLevel.Set(log.LevelDebug)
	case "warn":
		logLevel.Set(log.LevelWarn)
	case "error":
		logLevel.Set(log.LevelError)
	}

	h := log.HandlerOptions{Level: logLevel}.NewTextHandler(os.Stderr)
	log.SetDefault(log.New(h))

	if app_version == "" {
		app_version = "local-dev"
	}
	log := log.With("app_version", app_version, "app_build", app_build)

	// Loads config
	path, a := os.LookupEnv("CONFIG")
	if !a {
		path = "config.yaml"
	}
	log = log.With("configPath", path)

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Warn("Error reading yaml", "err", err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Warn("Error Unmarshal config file", "err", err)
	}

	if err := defaults.Set(&conf); err != nil {
		log.Warn("Error setting default values", "err", err)
	}

	log.Debug("Configuration loaded")
}
