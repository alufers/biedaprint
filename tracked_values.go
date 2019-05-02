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
	// current print values
	"isPrinting": &trackedValue{
		Name:        "isPrinting",
		DisplayType: "boolean",
		Value:       false,
		subscribers: []*websocket.Conn{},
	},
	"printProgress": &trackedValue{
		Name:              "printProgress",
		Unit:              "%",
		DisplayType:       "number",
		Value:             0.0,
		MinUpdateInterval: time.Second * 1,
		subscribers:       []*websocket.Conn{},
	},
	"printOriginalName": &trackedValue{
		Name:        "printOriginalName",
		DisplayType: "string",
		Value:       "",
		subscribers: []*websocket.Conn{},
	},
	"printStartTime": &trackedValue{
		Name:        "printStartTime",
		DisplayType: "time",
		Value:       nil,
		subscribers: []*websocket.Conn{},
	},
	"printCurrentLayer": &trackedValue{
		Name:        "printCurrentLayer",
		DisplayType: "number",
		Value:       0,
		subscribers: []*websocket.Conn{},
	},
}

type trackedValue struct {
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	DisplayType string `json:"displayType"`
	PlotColor   string `json:"plotColor"`

	Value interface{} `json:"value"`

	LastUpdate        time.Time `json:"lastUpdate"`
	lastSent          time.Time
	MinUpdateInterval time.Duration `json:"minUpdateInterval"`

	History          []interface{} `json:"history"`
	MaxHistoryLength int           `json:"maxHistoryLength"`

	subscribers []*websocket.Conn
}

func (tv *trackedValue) updateValue(val interface{}) {

	go func() {
		handlerMutex.Lock()
		defer handlerMutex.Unlock()
		if tv.MaxHistoryLength != 0 {
			if len(tv.History) >= tv.MaxHistoryLength {
				tv.History = append(tv.History[1:], val)
			} else {
				tv.History = append(tv.History, val)
			}
		}
		tv.Value = val
		tv.LastUpdate = time.Now()
		if time.Now().Sub(tv.lastSent) > tv.MinUpdateInterval {
			tv.lastSent = time.Now()
			tv.broadcastValue()
		}

	}()

}

func (tv *trackedValue) broadcastValue() {
	for _, s := range tv.subscribers {
		s.WriteJSON(jd{
			"type": "trackedValueUpdated",
			"data": jd{
				"name":  tv.Name,
				"value": tv.Value,
			},
		})
	}
}
