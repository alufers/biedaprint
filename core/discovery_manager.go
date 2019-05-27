package core

import (
	"log"

	"github.com/grandcat/zeroconf"
)

type DiscoveryManager struct {
	app *App
}

func NewDiscoveryManager(app *App) *DiscoveryManager {
	return &DiscoveryManager{
		app: app,
	}
}

func (dm *DiscoveryManager) Init() {
	log.Printf("Starting Zeroconf DiscoveryManager!")
	_, err := zeroconf.Register("Biedaprint instance", "_biedaprint._tcp", "local.", 4444, []string{}, nil)
	if err != nil {
		panic(err)
	}
}
