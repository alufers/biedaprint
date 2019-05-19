package main

import (
	"log"
	"os"
	"os/exec"
)

func runStartupCommand() {
	if globalSettings.StartupCommand != "" {
		cmd := exec.Command("sh", "-c", globalSettings.StartupCommand)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to execute starup command %v", err)
		}

	}
}
