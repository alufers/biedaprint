package main

import (
	"time"

	"github.com/gorilla/websocket"
)

var trackedValues = map[string]*trackedValue{
	"hotendTemperature": &trackedValue{
		PlotColor:        "red",
		Name:             "hotendTemperature",
		Unit:             "°C",
		DisplayType:      "plot",
		Value:            0,
		MaxHistoryLength: 300,
		History:          []interface{}{},
		subscribers:      []*websocket.Conn{},
	},
	"targetHotendTemperature": &trackedValue{
		PlotColor:        "blue",
		Name:             "targetHotendTemperature",
		Unit:             "°C",
		DisplayType:      "plot",
		Value:            0,
		MaxHistoryLength: 300,
		History:          []interface{}{},
		subscribers:      []*websocket.Conn{},
	},
	"printProgress": &trackedValue{
		Name:             "printProgress",
		Unit:             "%",
		DisplayType:      "number",
		Value:            0.0,
		MaxHistoryLength: 0,
		subscribers:      []*websocket.Conn{},
	},
}

type trackedValue struct {
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	DisplayType string `json:"displayType"`
	PlotColor   string `json:"plotColor"`

	Value interface{} `json:"value"`

	LastUpdate time.Time `json:"lastUpdate"`

	History          []interface{} `json:"history"`
	MaxHistoryLength int           `json:"maxHistoryLength"`

	subscribers []*websocket.Conn
}

func (tv *trackedValue) updateValue(val interface{}) {
	if tv.MaxHistoryLength != 0 {
		if len(tv.History) >= tv.MaxHistoryLength {
			tv.History = append(tv.History[1:], val)
		} else {
			tv.History = append(tv.History, val)
		}
	}
	tv.LastUpdate = time.Now()
	go func() {
		handlerMutex.Lock()
		defer handlerMutex.Unlock()
		for _, s := range tv.subscribers {
			s.WriteJSON(jd{
				"type": "trackedValueUpdated",
				"data": jd{
					"name":  tv.Name,
					"value": val,
				},
			})
		}
	}()

}
