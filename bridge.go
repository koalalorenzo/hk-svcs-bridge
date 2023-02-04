package main

import (
	"github.com/brutella/hap/accessory"
)

func SetupBridge() *accessory.Bridge {
	return accessory.NewBridge(accessory.Info{
		Name: conf.Name,
	})
}