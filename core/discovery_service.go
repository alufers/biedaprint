package core

import (
	"fmt"
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

//Init starts the service
func (dm *DiscoveryService) Init() error {
	log.Printf("Starting Zeroconf DiscoveryService!")
	_, err := zeroconf.Register("Biedaprint instance", "_biedaprint._tcp", "local.", 4444, []string{}, nil)
	if err != nil {
		return fmt.Errorf("failed to register Zeroconf: %w", err)
	}
	return nil
}
