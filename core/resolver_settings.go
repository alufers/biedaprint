package core

import "context"

func (r *mutationResolver) UpdateSettings(ctx context.Context, settings interface{}) (interface{}, error) {
	r.App.SettingsService.UpdateAllSettings(settings)
	return r.App.SettingsService.GetAllSettings(), nil
}

func (r *queryResolver) SerialPorts(ctx context.Context) ([]string, error) {
	return []string{"/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/ttyUSB2", "/dev/ttyUSB3", "/dev/ttyACM0", "/dev/ttyACM1", "/dev/ttyACM2", "/dev/cu.wchusbserial14d10"}, nil
}

func (r *queryResolver) Settings(ctx context.Context, path *string) (interface{}, error) {
	if path == nil {
		empty := ""
		path = &empty
	}
	set, err := r.App.SettingsService.GetValue(*path)
	return &set, err
}
