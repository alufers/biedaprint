package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

/*
SettingsPath is the path to the settings.
*/
const SettingsPath = "settings.json"

/*
SettingsService handles loading and saving the apps settings.
*/
type SettingsService struct {
	app      *App
	mutex    *sync.RWMutex
	settings interface{}
}

/*
NewSettingsService constructs a SettingsService
*/
func NewSettingsService(app *App) *SettingsService {
	return &SettingsService{
		app:   app,
		mutex: &sync.RWMutex{},
		settings: map[string]interface{}{ // default settings
			"__v": 2,
			"general": map[string]interface{}{
				"dataPath":       "./biedaprint_data",
				"startupCommand": "",
			},
			"serial": map[string]interface{}{
				"baudRate":             250000,
				"dataBits":             7,
				"parity":               "EVEN",
				"scrollbackBufferSize": 10240,
				"serialPort":           "<invalid>",
			},
			"temperatures": map[string]interface{}{
				"temperaturePresets": []interface{}{
					map[string]interface{}{
						"hotbedTemperature": 60,
						"hotendTemperature": 200,
						"name":              "PLA",
					},
					map[string]interface{}{
						"hotbedTemperature": 95,
						"hotendTemperature": 230,
						"name":              "ABS",
					},
				},
			},
		},
	}
}

/*
Init prepares the service.
*/
func (serv *SettingsService) Init() error {
	return serv.loadSettings()
}

/*
loadSettings loads the settings from a json file
*/
func (serv *SettingsService) loadSettings() error {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()
	file, err := ioutil.ReadFile(SettingsPath)
	if err != nil {
		log.Printf("Failed to load settings: %v, trying to save...", err)
		serv.mutex.Unlock() // unlock the mutex so that the settings can be saved
		err = serv.saveSettings()
		serv.mutex.Lock()
		if err != nil {
			return fmt.Errorf("failed to save default settings: %w", err)
		}
		file, err = ioutil.ReadFile(SettingsPath)
		if err != nil {
			return fmt.Errorf("failed to load default settings after saving: %w", err)
		}
	}

	err = json.Unmarshal([]byte(file), &serv.settings)
	if err != nil {
		return fmt.Errorf("failed to parse %v: %w", SettingsPath, err)
	}
	log.Printf("Loaded %v", SettingsPath)
	return nil
}

func (serv *SettingsService) saveSettings() error {
	serv.mutex.RLock()
	defer serv.mutex.RUnlock()
	settingsJSON, err := json.MarshalIndent(serv.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to stringify settings: %w", err)
	}
	err = ioutil.WriteFile(SettingsPath, settingsJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}
	return err
}

/*
GetAllSettings deep-copies the root settings object and returns it. It is safe to be called from multiple goroutines.
*/
func (serv *SettingsService) GetAllSettings() interface{} {
	serv.mutex.RLock()
	defer serv.mutex.RUnlock()
	return serv.copyValue(serv.settings)
}

/*
UpdateAllSettings replaces the whole settings object with the given value. It is safe to be called from multiple goroutines.
*/
func (serv *SettingsService) UpdateAllSettings(newSettings interface{}) {
	serv.mutex.Lock()
	serv.settings = newSettings
	serv.mutex.Unlock()
	serv.saveSettings()
}

/*
GetString returns a string at a path. Returns an error if the value does not exist or is not a string. It is safe to be called from multiple goroutines.
*/
func (serv *SettingsService) GetString(path string) (str string, err error) {

	val, err := serv.GetValue(path)
	if err != nil {
		return
	}
	str, ok := val.(string)
	if !ok {
		err = fmt.Errorf("The settings value at %v is not a string. It is a %T", path, val)
		return
	}
	return
}

/*
GetUint returns a uint at a path. Returns an error if the value does not exist or is not an uint. It is safe to be called from multiple goroutines.
*/
func (serv *SettingsService) GetUint(path string) (u uint, err error) {

	val, err := serv.GetValue(path)
	if err != nil {
		return
	}
	switch f := val.(type) {
	case float64:
		u = uint(f)
	case json.Number:
		i64, conversionError := f.Int64()
		if conversionError != nil {
			err = errors.Wrap(conversionError, "failed to convert json.Number to Int64")
			return
		}
		u = uint(i64)
	default:
		err = fmt.Errorf("The settings value at %v is not an float64 or json.Number. It is a %T", path, val)
		return
	}

	return
}

/*
GetInt64 returns a int64 at a path. Returns an error if the value does not exist or is not an int64. It is safe to be called from multiple goroutines.
*/
func (serv *SettingsService) GetInt64(path string) (i int64, err error) {

	val, err := serv.GetValue(path)
	if err != nil {
		return
	}
	switch f := val.(type) {
	case float64:
		i = int64(f)
	case json.Number:
		var conversionError error
		i, conversionError = f.Int64()
		if conversionError != nil {
			err = errors.Wrap(conversionError, "failed to convert json.Number to Int64")
			return
		}
	default:
		err = fmt.Errorf("The settings value at %v is not an float64 or json.Number. It is a %T", path, val)
		return
	}

	return
}

/*
GetValue returns a value in the settings at a path.
*/
func (serv *SettingsService) GetValue(path string) (ret interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			ret = nil
		}
	}()
	ret = serv.getValue(path)
	return
}

func (serv *SettingsService) getValue(path string) interface{} {
	if path == "" {
		return serv.GetAllSettings()
	}
	serv.mutex.RLock()
	defer serv.mutex.RUnlock()
	segments := strings.Split(path, ".")
	accumulator := serv.settings

	for index, segment := range segments {
		switch v := accumulator.(type) {
		case map[string]interface{}:
			prop, ok := v[segment]
			if !ok {
				panic(fmt.Errorf("no such property at segment %v (%v)", index, segment))
			}
			accumulator = prop
		case []interface{}:
			indexValue, err := strconv.Atoi(segment)
			if err != nil {
				panic(errors.Wrapf(err, "segment %v (%v) failed to be parsed as a number when indexing an array", index, segment))
			}
			if indexValue >= len(v) {
				panic(fmt.Errorf("segment %v (%v) out of range when indexing an array", index, segment))
			}
			accumulator = v[indexValue]
		}
	}
	return serv.copyValue(accumulator)
}

func (serv *SettingsService) copyValue(i interface{}) interface{} {
	if i == nil {
		return nil
	}
	switch v := i.(type) {
	case float64:
		return v
	case string:
		return v
	case bool:
		return v
	case json.Number:
		return v
	case map[string]interface{}:
		newMap := map[string]interface{}{}
		for k, val := range v {
			newMap[k] = serv.copyValue(val)
		}
		return newMap
	case []interface{}:
		newArr := make([]interface{}, len(v))
		for index, val := range v {
			newArr[index] = serv.copyValue(val)
		}
		return newArr
	}
	panic(errors.Errorf("unsupported type to copy: %T", i))
}
