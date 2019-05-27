package core

import "context"

func (r *mutationResolver) UpdateSettings(ctx context.Context, settings NewSettings) (*Settings, error) {
	r.App.settingsMutex.Lock()
	defer func() {
		r.App.settingsMutex.Unlock()
		r.App.saveSettings()
	}()
	r.App.settings.SerialPort = settings.SerialPort
	r.App.settings.BaudRate = settings.BaudRate
	r.App.settings.ScrollbackBufferSize = settings.ScrollbackBufferSize
	r.App.settings.DataPath = settings.DataPath
	r.App.settings.Parity = settings.Parity
	r.App.settings.DataBits = settings.DataBits
	r.App.settings.StartupCommand = settings.StartupCommand
	r.App.settings.TemperaturePresets = []*TemperaturePreset{}
	for _, tp := range settings.TemperaturePresets {
		r.App.settings.TemperaturePresets = append(r.App.settings.TemperaturePresets, &TemperaturePreset{
			Name:              tp.Name,
			HotendTemperature: tp.HotendTemperature,
			HotbedTemperature: tp.HotbedTemperature,
		})
	}
	return r.App.settings, nil
}

func (r *queryResolver) SerialPorts(ctx context.Context) ([]string, error) {
	return []string{"/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/ttyUSB2", "/dev/ttyUSB3", "/dev/ttyACM0", "/dev/ttyACM1", "/dev/ttyACM2", "/dev/cu.wchusbserial14d10"}, nil
}

func (r *queryResolver) Settings(ctx context.Context) (*Settings, error) {
	set := r.App.GetSettings()
	return &set, nil
}
