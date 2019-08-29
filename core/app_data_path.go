package core

/*
GetDataPath returns the data path where all the data should be kept.
*/
func (app *App) GetDataPath() string {
	path, err := app.SettingsService.GetString("general.dataPath")
	if err != nil {
		panic(err)
	}
	return path
}
