package core

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
	return app
}

/*
Init initializes the app
*/
func (app *App) Init() {
	app.SettingsService.Init()
	app.runStartupCommand()

	app.RecentCommandsService.LoadRecentCommands()

	app.DiscoveryService.Init()
	app.PrinterService.Init()
	app.HeatingService.Init()
}
