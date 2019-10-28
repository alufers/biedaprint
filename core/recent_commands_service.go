package core

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"sync"
)

/*
RecentCommandsService handles saving the commands entered by the user in the serial console window for later use by pressing the up arrow just like in a normal shell.
*/
type RecentCommandsService struct {
	app                 *App
	recentCommands      []string
	recentCommandsMutex *sync.RWMutex
}

/*
NewRecentCommandsService constructs a RecentCommandsService.
*/
func NewRecentCommandsService(app *App) *RecentCommandsService {
	return &RecentCommandsService{
		app:                 app,
		recentCommands:      []string{},
		recentCommandsMutex: &sync.RWMutex{},
	}
}

/*
Init initializes the recent commands service.
*/
func (rcm *RecentCommandsService) Init() error {
	err := rcm.LoadRecentCommands()
	if err != nil {
		log.Print(err)
	}
	return nil
}

/*
AddRecentCommand saves a new command in the recent commands file.
*/
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

/*
LoadRecentCommands should be run on start, to load the recently used commands.
*/
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

/*
GetRecentCommands retunrs the recent commands.
*/
func (rcm *RecentCommandsService) GetRecentCommands() []string {
	rcm.recentCommandsMutex.RLock()
	defer rcm.recentCommandsMutex.RUnlock()
	return rcm.recentCommands
}
