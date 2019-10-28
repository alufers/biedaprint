package core

import "log"

/*
App is the root object holding all the diffrent services.
*/
type App struct {
	SettingsService       *SettingsService
	PrinterService        *PrinterService
	RecentCommandsService *RecentCommandsService
	TrackedValuesService  *TrackedValuesService
	DiscoveryService      *DiscoveryService
	HeatingService        *HeatingService
	ManualMovementService *ManualMovementService
	HTTPService           *HTTPService
	DBService             *DBService
	StartupCommandService *StartupCommandService
	DataDirectoryService  *DataDirectoryService
}

/*
NewApp constructs an App
*/
func NewApp() *App {
	app := &App{}
	app.SettingsService = NewSettingsService(app)
	app.PrinterService = NewPrinterService(app)
	app.RecentCommandsService = NewRecentCommandsService(app)
	app.TrackedValuesService = NewTrackedValuesService(app)
	app.DiscoveryService = NewDiscoveryService(app)
	app.HeatingService = NewHeatingService(app)
	app.ManualMovementService = NewManualMovementService(app)
	app.HTTPService = NewHTTPService(app)
	app.DBService = NewDBService(app)
	app.StartupCommandService = NewStartupCommandService(app)
	app.DataDirectoryService = NewDataDirectoryService(app)
	return app
}

/*
Init initializes the app
*/
func (app *App) Init() {
	app.SettingsService.Init()
	err := InitMany(
		app.SettingsService,
		app.StartupCommandService,
		app.DataDirectoryService,
		app.DBService,
		app.RecentCommandsService,
		app.DiscoveryService,
		app.PrinterService,
		app.HeatingService,
	)
	if err != nil {
		log.Fatal(err)
	}
}
