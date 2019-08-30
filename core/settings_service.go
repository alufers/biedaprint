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
		app:      app,
		mutex:    &sync.RWMutex{},
		settings: nil,
	}
}

/*
Init prepares the service.
*/
func (serv *SettingsService) Init() {
	serv.loadSettings()
}

/*
loadSettings loads the settings from a json file
*/
func (serv *SettingsService) loadSettings() {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()
	file, err := ioutil.ReadFile(SettingsPath)
	if err != nil {
		log.Printf("Failed to load settings: %v, trying to save...", err)
		serv.mutex.Unlock() // unlock the mutex so that the settings can be saved
		err = serv.saveSettings()
		serv.mutex.Lock()
		if err != nil {
			log.Fatalf("Failed to save default settings: %v", err)
			return
		}
		file, err = ioutil.ReadFile(SettingsPath)
		if err != nil {
			log.Fatalf("Failed to load default settings after saving: %v", err)
		}
	}

	err = json.Unmarshal([]byte(file), &serv.settings)
	if err != nil {
		log.Fatalf("Failed to parse %v. Check your syntax. %v", SettingsPath, err)
	}
	log.Printf("Loaded %v", SettingsPath)
}

func (serv *SettingsService) saveSettings() error {
	serv.mutex.RLock()
	defer serv.mutex.RUnlock()
	settingsJSON, err := json.Marshal(serv.settings)
	if err != nil {
		log.Printf("Failed to stringify settings %v", err)
		return err
	}
	err = ioutil.WriteFile(SettingsPath, settingsJSON, 0644)
	if err != nil {
		log.Printf("Failed to save settings %v", err)
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
	f, ok := val.(float64)
	if !ok {
		err = fmt.Errorf("The settings value at %v is not an float64. It is a %T", path, val)
		return
	}
	u = uint(f)
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
	f, ok := val.(float64)
	if !ok {
		err = fmt.Errorf("The settings value at %v is not an uint. It is a %T", path, val)
		return
	}
	i = int64(f)
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
	panic(errors.New("unsupported type to copy"))
}
