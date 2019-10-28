package core

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const RoomTemperatureMax = 38.0

type HeatingService struct {
	app                   *App
	autoreportingAttempts int
	temperatureReportChan chan bool
	lastHotendTarget      float64
	hotendHeatingStart    *time.Time
	lastHotbedTarget      float64
	hotbedHeatingStart    *time.Time
	HotendTimings         map[float64]time.Duration // a map which holds the time required to heat up the element to a given key temperature
	HotbedTimings         map[float64]time.Duration
}

func NewHeatingService(app *App) *HeatingService {
	return &HeatingService{
		app:                   app,
		temperatureReportChan: make(chan bool, 14),
		HotendTimings:         map[float64]time.Duration{},
		HotbedTimings:         map[float64]time.Duration{},
	}
}

func (hs *HeatingService) Init() error {
	err := hs.loadTemperatureTimings()
	if err != nil {
		log.Printf("failed to load temperature timings: %v", err)
	}
	go hs.communicateWithPrinter()
	return nil
}

func (hs *HeatingService) temperatureTimingsFilePath() string {
	dataPath := hs.app.GetDataPath()
	return filepath.Join(dataPath, "temperature_timings.meta")
}

func (hs *HeatingService) loadTemperatureTimings() error {
	f, err := os.Open(hs.temperatureTimingsFilePath())
	if err != nil {
		return errors.Wrap(err, "failed to open temperature timings file for reading")
	}
	defer f.Close()

	deco := gob.NewDecoder(f)

	err = deco.Decode(&hs.HotendTimings)
	err = deco.Decode(&hs.HotbedTimings)
	if err != nil {
		return errors.Wrap(err, "failed to decode temperature timings")
	}
	return nil
}

func (hs *HeatingService) saveTemperatureTimings() error {
	f, err := os.OpenFile(hs.temperatureTimingsFilePath(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to temperature timings file for writing")
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(hs.HotendTimings)
	err = enc.Encode(hs.HotbedTimings)
	if err != nil {
		return errors.Wrap(err, "failed to encode gcode meta file")
	}
	return nil
}

func (hs *HeatingService) processTemperatureTimings(currentHotendTemperature, targetHotendTemperature, currentHotbedTemperature, targetHotbedTemperature float64) {
	var shouldSaveTimings bool
	if currentHotendTemperature < RoomTemperatureMax && targetHotendTemperature > hs.lastHotendTarget {
		// heating on hotend started from room temperature
		t := time.Now()
		hs.hotendHeatingStart = &t
	}
	if targetHotendTemperature < hs.lastHotendTarget {
		// heating aborted or temperature lowered. Abort the measurement
		hs.hotendHeatingStart = nil
	}
	hs.lastHotendTarget = targetHotendTemperature
	if currentHotendTemperature >= targetHotendTemperature {
		if hs.hotendHeatingStart != nil {
			// hotend temperature reached. Save the measurement
			hs.HotendTimings[targetHotendTemperature] = time.Now().Sub(*hs.hotendHeatingStart)
			log.Printf("Hotend heated up to %v in %v", targetHotendTemperature, hs.HotendTimings[targetHotendTemperature].String())
			hs.hotendHeatingStart = nil
			shouldSaveTimings = true
		}
	}
	if currentHotbedTemperature < RoomTemperatureMax && targetHotbedTemperature > hs.lastHotbedTarget {
		t := time.Now()
		hs.hotbedHeatingStart = &t
	}
	if targetHotbedTemperature < hs.lastHotbedTarget {
		// heating aborted or temperature lowered. Abort the measurement
		hs.hotbedHeatingStart = nil
	}
	hs.lastHotbedTarget = targetHotbedTemperature
	if currentHotbedTemperature >= targetHotbedTemperature {
		if hs.hotbedHeatingStart != nil {
			// hotbed temperature reached. Save the measurement
			hs.HotbedTimings[targetHotbedTemperature] = time.Now().Sub(*hs.hotbedHeatingStart)
			log.Printf("Hotbed heated up to %v in %v", targetHotbedTemperature, hs.HotbedTimings[targetHotbedTemperature].String())
			hs.hotbedHeatingStart = nil
			shouldSaveTimings = true
		}
	}
	if shouldSaveTimings {
		err := hs.saveTemperatureTimings()
		if err != nil {
			log.Printf("failed to save temperature timings: %v", err)
		}
	}
}

func (hs *HeatingService) processTemperatureReport(temperatureReport string) {
	select {
	case hs.temperatureReportChan <- true:
	default:
	}
	var temp float64
	var target float64
	var power int
	var bedTemp float64
	var bedTargetTemp float64
	var bedPower float64
	if strings.Contains(temperatureReport, "B:") { // has heated bed
		fmt.Sscanf(temperatureReport, "T:%f /%f B:%f /%f @:%d B@:%d", &temp, &target, &bedTemp, &bedTargetTemp, &power, &bedPower)
	} else {
		fmt.Sscanf(temperatureReport, "T:%f /%f @:%d", &temp, &target, &power)
	}
	hs.app.TrackedValuesService.TrackedValues["hotendTemperature"].UpdateValue(temp)
	hs.app.TrackedValuesService.TrackedValues["targetHotendTemperature"].UpdateValue(target)
	hs.app.TrackedValuesService.TrackedValues["hotbedTemperature"].UpdateValue(bedTemp)
	hs.app.TrackedValuesService.TrackedValues["targetHotbedTemperature"].UpdateValue(bedTargetTemp)
	hs.processTemperatureTimings(temp, target, bedTemp, bedTargetTemp)
}

func (hs *HeatingService) communicateWithPrinter() {
	for {
		hs.app.PrinterService.WaitForSerialReady()
		if hs.autoreportingAttempts < 5 {
			time.Sleep(4 * time.Second)
			hs.app.PrinterService.consoleWriteSem <- "M155 S1\r\n"
			hs.autoreportingAttempts++

		} else {
			log.Printf("Failed to recieve temeprature report using auto-reporting %d times. Reverting to polling", hs.autoreportingAttempts)
			for {
				hs.app.PrinterService.consoleWriteSem <- "M105\r\n"
				time.Sleep(time.Second)
			}
		}
	AutoreportingCheckLoop:
		for {
			select {
			case <-time.After(time.Second * 2):
				break AutoreportingCheckLoop
			case <-hs.temperatureReportChan:
				if hs.autoreportingAttempts > 0 {
					hs.autoreportingAttempts--
				}
			}
		}
	}
}

//SmartHeatUp uses temperature timings to heat up all the elements in a most time and energy-effiecient way
func (hs *HeatingService) SmartHeatUp(hotendTarget, hotbedTarget float64, abort chan bool) {
	hotendTiming := hs.getClosestHotendTiming(hotendTarget)
	hotbedTiming := hs.getClosestHotbedTiming(hotbedTarget)
	if hotendTiming == 0 || hotbedTiming == 0 {
		// no timings, abort and set the temperatures the old way
		hs.app.PrinterService.consoleWriteSem <- fmt.Sprintf("M104 %v\r\n", hotendTarget)
		hs.app.PrinterService.consoleWriteSem <- fmt.Sprintf("M140 %v\r\n", hotbedTarget)
		return
	}
	hs.app.PrinterService.consoleWriteSem <- fmt.Sprintf("M140 %v\r\n", hotbedTarget)
	if hotbedTiming-hotendTiming > 0 {
		select {
		case <-time.After(hotbedTiming - hotendTiming):
		case <-abort:
			return
		}
	}
	hs.app.PrinterService.consoleWriteSem <- fmt.Sprintf("M104 %v\r\n", hotendTarget)
}

func (hs *HeatingService) getClosestHotendTiming(target float64) time.Duration {
	var bestTiming time.Duration
	var bestTimingTempDiff = math.Inf(1)
	for timingTemp, timing := range hs.HotendTimings {
		diff := math.Abs(timingTemp - target)
		if diff < bestTimingTempDiff {
			bestTimingTempDiff = diff
			bestTiming = timing
		}
	}
	return bestTiming
}

func (hs *HeatingService) getClosestHotbedTiming(target float64) time.Duration {
	var bestTiming time.Duration
	var bestTimingTempDiff = math.Inf(1)
	for timingTemp, timing := range hs.HotbedTimings {
		diff := math.Abs(timingTemp - target)
		if diff < bestTimingTempDiff {
			bestTimingTempDiff = diff
			bestTiming = timing
		}
	}
	return bestTiming
}
