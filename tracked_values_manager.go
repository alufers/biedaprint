package biedaprint

import "time"

type TrackedValuesManager struct {
	app           *App
	TrackedValues map[string]*TrackedValueInternal
}

var zero = 0

func NewTrackedValuesManager(app *App) *TrackedValuesManager {
	return &TrackedValuesManager{
		app: app,
		TrackedValues: map[string]*TrackedValueInternal{
			"hotendTemperature": NewTrackedValueInternal(&TrackedValue{
				PlotColor:        "red",
				Name:             "hotendTemperature",
				Unit:             "°C",
				DisplayType:      TrackedValueDisplayTypePlot,
				Value:            0,
				MaxHistoryLength: 300,
				History:          []interface{}{},
			}),
			"targetHotendTemperature": NewTrackedValueInternal(&TrackedValue{
				PlotColor:        "blue",
				Name:             "targetHotendTemperature",
				Unit:             "°C",
				DisplayType:      TrackedValueDisplayTypePlot,
				Value:            0,
				MaxHistoryLength: 300,
				History:          []interface{}{},
			}),
			"hotbedTemperature": NewTrackedValueInternal(&TrackedValue{
				PlotColor:        "orange",
				Name:             "hotbedTemperature",
				Unit:             "°C",
				DisplayType:      TrackedValueDisplayTypePlot,
				Value:            0,
				MaxHistoryLength: 300,
				History:          []interface{}{},
			}),
			"targetHotbedTemperature": NewTrackedValueInternal(&TrackedValue{
				PlotColor:        "green",
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
