package core

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type App struct {
	settings              *Settings
	settingsMutex         *sync.RWMutex
	PrinterService        *PrinterService
	RecentCommandsService *RecentCommandsService
	TrackedValuesService  *TrackedValuesService
	DiscoveryService      *DiscoveryService
	router                *gin.Engine
}

func NewApp() *App {
	app := &App{
		settingsMutex: &sync.RWMutex{},
		settings: &Settings{
			SerialPort:           "<invalid>",
			BaudRate:             250000,
			ScrollbackBufferSize: 1024 * 10, // 10 KiB
			Parity:               SerialParityEven,
			DataBits:             7,
			DataPath:             "./biedaprint_data",
			StartupCommand:       "",
			TemperaturePresets: []*TemperaturePreset{
				&TemperaturePreset{
					Name:              "PLA",
					HotendTemperature: 200,
					HotbedTemperature: 60,
				},
				&TemperaturePreset{
					Name:              "ABS",
					HotendTemperature: 230,
					HotbedTemperature: 95,
				},
			},
		},
	}
	app.PrinterService = NewPrinterService(app)
	app.RecentCommandsService = NewRecentCommandsService(app)
	app.TrackedValuesService = NewTrackedValuesService(app)
	app.DiscoveryService = NewDiscoveryService(app)
	return app
}

func (app *App) Init() {
	app.loadSettings()

	app.runStartupCommand()

	app.RecentCommandsService.LoadRecentCommands()

	app.DiscoveryService.Init()
	app.PrinterService.Init()
}

//GetSettings returns a copy of settigns, safe for concurrent use
func (app *App) GetSettings() Settings {
	app.settingsMutex.RLock()
	defer app.settingsMutex.RUnlock()
	return *app.settings
}
