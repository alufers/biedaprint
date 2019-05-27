package core

import (
	"log"

	"github.com/grandcat/zeroconf"
)

//DiscoveryService handles registering the app using zeroconf so that it can be detected by the various clients.
type DiscoveryService struct {
	app *App
}

//NewDiscoveryService constructs the discovery manager
func NewDiscoveryService(app *App) *DiscoveryService {
	return &DiscoveryService{
		app: app,
	}
}

func (dm *DiscoveryService) Init() {
	log.Printf("Starting Zeroconf DiscoveryService!")
	_, err := zeroconf.Register("Biedaprint instance", "_biedaprint._tcp", "local.", 4444, []string{}, nil)
	if err != nil {
		panic(err)
	}
}
