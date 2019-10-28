package core

import (
	"log"
	"os"
	"os/exec"
)

//StartupCommandService handles running the startup comamdn when the application is started
type StartupCommandService struct {
	app *App
}

//NewStartupCommandService constructs a StartupCommandService
func NewStartupCommandService(app *App) *StartupCommandService {
	return &StartupCommandService{
		app: app,
	}
}

//Init initializes the service, by running the command. On error it only prints an information, it always succeeds.
func (scs *StartupCommandService) Init() error {
	sc, err := scs.app.SettingsService.GetString("general.startupCommand")
	if err != nil {
		log.Printf("Failed to get general.startupCommand setting: %v", err)
		return nil
	}
	if sc != "" {
		cmd := exec.Command("sh", "-c", sc)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to execute startup command: %v", err)
		}

	}
	return nil
}
