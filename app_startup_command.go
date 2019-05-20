package biedaprint

import (
	"log"
	"os"
	"os/exec"
)

func (app *App) runStartupCommand() {
	set := app.GetSettings()
	if set.StartupCommand != "" {
		cmd := exec.Command("sh", "-c", set.StartupCommand)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to execute starup command %v", err)
		}

	}
}
