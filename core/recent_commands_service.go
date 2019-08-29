package core

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"sync"
)

type RecentCommandsService struct {
	app                 *App
	recentCommands      []string
	recentCommandsMutex *sync.RWMutex
}

func NewRecentCommandsService(app *App) *RecentCommandsService {
	return &RecentCommandsService{
		app:                 app,
		recentCommands:      []string{},
		recentCommandsMutex: &sync.RWMutex{},
	}
}

func (rcm *RecentCommandsService) AddRecentCommand(cmd string) error {
	rcm.recentCommandsMutex.Lock()
	defer rcm.recentCommandsMutex.Unlock()
	rcm.recentCommands = append(rcm.recentCommands, cmd)
	f, err := os.OpenFile(filepath.Join(rcm.app.GetDataPath(), "recent_commands.meta"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(rcm.recentCommands)
	if err != nil {
		return err
	}
	return nil
}

func (rcm *RecentCommandsService) LoadRecentCommands() error {
	rcm.recentCommandsMutex.Lock()
	defer rcm.recentCommandsMutex.Unlock()
	f, err := os.Open(filepath.Join(rcm.app.GetDataPath(), "recent_commands.meta"))
	if err != nil {
		return err
	}
	defer f.Close()

	deco := gob.NewDecoder(f)
	err = deco.Decode(&rcm.recentCommands)
	if err != nil {
		return err
	}
	return nil
}

func (rcm *RecentCommandsService) GetRecentCommands() []string {
	rcm.recentCommandsMutex.RLock()
	defer rcm.recentCommandsMutex.RUnlock()
	return rcm.recentCommands
}
