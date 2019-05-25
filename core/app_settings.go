package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const SettingsPath = "settings.json"

func (app *App) loadSettings() {
	app.settingsMutex.Lock()
	defer app.settingsMutex.Unlock()
	file, err := ioutil.ReadFile(SettingsPath)
	if err != nil {
		log.Printf("Failed to load settings: %v, trying to save...", err)
		app.settingsMutex.Unlock() // unlock the mutex so that the settings can be saved
		err = app.saveSettings()
		app.settingsMutex.Lock()
		if err != nil {
			log.Fatalf("Failed to save  default settings: %v", err)
			return
		}
		file, err = ioutil.ReadFile(SettingsPath)
		if err != nil {
			log.Fatalf("Failed to load default settings after saving: %v", err)
		}
	}

	err = json.Unmarshal([]byte(file), &app.settings)
	if err != nil {
		log.Fatalf("Failed to parse %v. Check your syntax. %v", SettingsPath, err)
	}
	log.Printf("Loaded %v", SettingsPath)
}

func (app *App) saveSettings() error {
	app.settingsMutex.RLock()
	defer app.settingsMutex.RUnlock()
	settingsJSON, err := json.Marshal(app.settings)
	if err != nil {
		log.Printf("Failed to stringify settings %v", err)
		return err
	}
	err = ioutil.WriteFile(SettingsPath, settingsJSON, 0644)
	if err != nil {
		log.Printf("Failed to save settings %v", err)
	}
	return err
}
