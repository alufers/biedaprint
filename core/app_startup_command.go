package core

import (
	"log"
	"os"
	"os/exec"
)

func (app *App) runStartupCommand() {
	sc, err := app.SettingsService.GetString("general.startupCommand")
	if err != nil {
		log.Printf("Failed to get general.startupCommand setting: %v", err)
		return
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
}
