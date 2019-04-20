package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Settings struct {
	SerialPort           string `json:"serialPort"`
	BaudRate             int    `json:"baudRate"`
	ScrollbackBufferSize int    `json:"scrollbackBufferSize"`
}

var globalSettings = &Settings{
	SerialPort:           "<invalid>",
	BaudRate:             250000,
	ScrollbackBufferSize: 1024 * 10, // 10 KiB
}

func loadSettings() {
	file, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Printf("Failed to load settings: %v, trying to save...", err)
		err = saveSettings()
		if err != nil {
			log.Fatalf("Failed to save  default settings: %v", err)
			return
		}
		file, err = ioutil.ReadFile("settings.json")
		if err != nil {
			log.Fatalf("Failed to load default settings after saving: %v", err)
		}
	}

	err = json.Unmarshal([]byte(file), &globalSettings)
	if err != nil {
		log.Fatalf("Failed to parse settings.json. Check your syntax. %v", err)
	}
}

func saveSettings() error {
	settingsJson, err := json.Marshal(globalSettings)
	if err != nil {
		log.Printf("Failed to stringify settings %v", err)
		return err
	}
	err = ioutil.WriteFile("settings.json", settingsJson, 0644)
	if err != nil {
		log.Printf("Failed to save settings %v", err)
	}
	return err
}
