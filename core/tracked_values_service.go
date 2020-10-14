package core

import "time"

/*
TrackedValuesService holds the definitions of tracked values. And tracks their histories.
*/
type TrackedValuesService struct {
	app           *App
	TrackedValues map[string]*TrackedValueInternal
}

var zero = 0

/*
NewTrackedValuesService constructs a TrackedValuesService.
*/
func NewTrackedValuesService(app *App) *TrackedValuesService {
	return &TrackedValuesService{
		app: app,
		TrackedValues: map[string]*TrackedValueInternal{
			"hotendTemperature": NewTrackedValueInternal(&TrackedValue{
				PlotColor:        "#ff3860",
				PlotDash:         []float64{},
				Name:             "hotendTemperature",
				Unit:             "°C",
				DisplayType:      TrackedValueDisplayTypePlot,
				Value:            0,
				MaxHistoryLength: 300,
				History:          []interface{}{},
			}),
			"targetHotendTemperature": NewTrackedValueInternal(&TrackedValue{
				PlotColor:        "#ff3860",
				PlotDash:         []float64{5, 5},
				Name:             "targetHotendTemperature",
				Unit:             "°C",
				DisplayType:      TrackedValueDisplayTypePlot,
				Value:            0,
				MaxHistoryLength: 300,
				History:          []interface{}{},
			}),
			"hotbedTemperature": NewTrackedValueInternal(&TrackedValue{
				PlotColor:        "#b86bff",
				PlotDash:         []float64{},
				Name:             "hotbedTemperature",
				Unit:             "°C",
				DisplayType:      TrackedValueDisplayTypePlot,
				Value:            0,
				MaxHistoryLength: 300,
				History:          []interface{}{},
			}),
			"fanSpeed": NewTrackedValueInternal(&TrackedValue{
				PlotColor:        "#a9a9a9",
				PlotDash:         []float64{},
				Name:             "fanSpeed",
				Unit:             "",
				DisplayType:      TrackedValueDisplayTypePlot,
				Value:            0,
				MaxHistoryLength: 300,
				History:          []interface{}{},
			}),
			"targetHotbedTemperature": NewTrackedValueInternal(&TrackedValue{
				PlotColor:        "#b86bff",
				PlotDash:         []float64{5, 5},
				Name:             "targetHotbedTemperature",
				Unit:             "°C",
				DisplayType:      TrackedValueDisplayTypePlot,
				Value:            0,
				MaxHistoryLength: 300,
				History:          []interface{}{},
			}),
			"serialStatus": NewTrackedValueInternal(&TrackedValue{
				Name:        "serialStatus",
				DisplayType: TrackedValueDisplayTypeString,
				Value:       "disconnected",
			}),
			"isPrinting": NewTrackedValueInternal(&TrackedValue{
				Name:        "isPrinting",
				DisplayType: TrackedValueDisplayTypeBoolean,
				Value:       false,
			}),
			"printProgress": NewTrackedValueInternal(&TrackedValue{
				Name:              "printProgress",
				Unit:              "%",
				DisplayType:       TrackedValueDisplayTypeNumber,
				Value:             0.0,
				MinUpdateInterval: 1000,
			}),
			"printOriginalName": NewTrackedValueInternal(&TrackedValue{
				Name:        "printOriginalName",
				DisplayType: TrackedValueDisplayTypeString,
				Value:       "",
			}),
			"printStartTime": NewTrackedValueInternal(&TrackedValue{
				Name:        "printStartTime",
				DisplayType: TrackedValueDisplayTypeTime,
				Value:       time.Now(),
			}),
			"printCurrentLayer": NewTrackedValueInternal(&TrackedValue{
				Name:        "printCurrentLayer",
				DisplayType: TrackedValueDisplayTypeNumber,
				Value:       0,
			}),
			"printTotalLayers": NewTrackedValueInternal(&TrackedValue{
				Name:        "printTotalLayers",
				DisplayType: TrackedValueDisplayTypeNumber,
				Value:       0,
			}),
		},
	}
}
