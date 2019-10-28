package core

import (
	"fmt"
	"os"
)

//DataDirectoryService handles running the startup comamdn when the application is started
type DataDirectoryService struct {
	app *App
}

//NewDataDirectoryService constructs a DataDirectoryService
func NewDataDirectoryService(app *App) *DataDirectoryService {
	return &DataDirectoryService{
		app: app,
	}
}

//Init initializes the service, by running the command. On error it only prints an information, it always succeeds.
func (dds *DataDirectoryService) Init() error {
	err := os.MkdirAll(dds.app.GetDataPath(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create the data directory: %w", err)
	}
	return nil
}
