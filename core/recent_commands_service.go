package core

import (
	"fmt"
)

/*
RecentCommandsService handles saving the commands entered by the user in the serial console window for later use by pressing the up arrow just like in a normal shell.
*/
type RecentCommandsService struct {
	app *App
}

/*
NewRecentCommandsService constructs a RecentCommandsService.
*/
func NewRecentCommandsService(app *App) *RecentCommandsService {
	return &RecentCommandsService{
		app: app,
	}
}

/*
Init initializes the recent commands service.
*/
func (rcm *RecentCommandsService) Init() error {
	err := rcm.app.DBService.DB.AutoMigrate(&RecentCommand{}).Error
	if err != nil {
		return fmt.Errorf("failed to migrate RecentCommand: %w", err)
	}
	return nil
}

/*
AddRecentCommand saves a new command in the recent commands file.
*/
func (rcm *RecentCommandsService) AddRecentCommand(cmd string) error {
	record := &RecentCommand{
		Command: cmd,
	}
	if err := rcm.app.DBService.DB.Save(record).Error; err != nil {
		return fmt.Errorf("failed to insert recent command: %v", err)
	}
	return nil
}

/*
GetRecentCommands retunrs the recent commands.
*/
func (rcm *RecentCommandsService) GetRecentCommands() (ret []string, err error) {
	var rows []*RecentCommand
	if dbErr := rcm.app.DBService.DB.Find(&rows).Error; dbErr != nil {
		err = fmt.Errorf("failed to query recent commands: %v", dbErr)
		return
	}
	ret = make([]string, len(rows))
	for i, r := range rows {
		ret[i] = r.Command
	}
	return
}

//RecentCommand represents a record stored in the databse anout a recent command
type RecentCommand struct {
	Model
	Command string
}
