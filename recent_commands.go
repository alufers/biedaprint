package main

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

var recentCommands = make([]string, 0)

func addRecentCommand(cmd string) error {
	recentCommands = append(recentCommands, cmd)
	f, err := os.OpenFile(filepath.Join(globalSettings.DataPath, "recent_commands.meta"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(recentCommands)
	if err != nil {
		return err
	}
	return nil
}

func loadRecentCommands() error {
	f, err := os.Open(filepath.Join(globalSettings.DataPath, "recent_commands.meta"))
	if err != nil {
		return err
	}
	defer f.Close()

	deco := gob.NewDecoder(f)
	err = deco.Decode(&recentCommands)
	if err != nil {
		return err
	}
	return nil
}
