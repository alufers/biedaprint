package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/mem"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	App *App
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

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
func (r *mutationResolver) ConnectToSerial(ctx context.Context, void *bool) (*bool, error) {
	return nil, r.App.PrinterManager.ConnectToSerial()
}
func (r *mutationResolver) DisconnectFromSerial(ctx context.Context, void *bool) (*bool, error) {
	return nil, r.App.PrinterManager.DisconnectFromSerial()
}
func (r *mutationResolver) SendGcode(ctx context.Context, cmd string) (*bool, error) {
	r.App.PrinterManager.consoleWriteSem <- cmd + "\r\n"
	return nil, nil
}
func (r *mutationResolver) SendConsoleCommand(ctx context.Context, cmd string) (*bool, error) {
	r.App.PrinterManager.consoleWriteSem <- cmd + "\r\n"
	return nil, r.App.RecentCommandsManager.AddRecentCommand(cmd)
}
func (r *mutationResolver) UploadGcode(ctx context.Context, file graphql.Upload) (*GcodeFileMeta, error) {
	dataPath := r.App.GetSettings().DataPath
	gcodeFilename := RandStringRunes(8) + ".gcode"
	// copy example
	f, err := os.OpenFile(filepath.Join(dataPath, "gcode_files/"+gcodeFilename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open gcode file for writing")
	}
	defer f.Close()
	_, err = io.Copy(f, file.File)
	if err != nil {
		return nil, errors.Wrap(err, "failed to copy file from request")
	}
	meta := &GcodeFileMeta{
		OriginalName:  file.Filename,
		UploadDate:    time.Now(),
		GcodeFileName: gcodeFilename,
	}
	err = meta.AnalyzeGcodeFile(dataPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to analyze gcode file")
	}
	err = meta.Save(dataPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save gcode file meta")
	}
	return meta, nil
}
func (r *mutationResolver) DeleteGcodeFile(ctx context.Context, gcodeFilename string) (*bool, error) {
	dataPath := r.App.GetSettings().DataPath
	gcodeName := filepath.Join(dataPath, "gcode_files/", gcodeFilename)
	gcodeMetaName := filepath.Join(dataPath, "gcode_files/", gcodeFilename+".meta")
	err := os.Remove(gcodeName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove gcode file")
	}
	err = os.Remove(gcodeMetaName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove gcode file meta")
	}
	return nil, nil
}
func (r *mutationResolver) StartPrintJob(ctx context.Context, gcodeFilename string) (*bool, error) {
	meta, err := loadGcodeFileMeta(filepath.Join(r.App.GetSettings().DataPath, "gcode_files/", gcodeFilename+".meta"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load gcode meta")
	}
	job := &PrintJobInternal{
		app: r.App,
		PrintJob: &PrintJob{
			GcodeMeta:   meta,
			StartedTime: time.Now(),
		},
	}

	select {
	case r.App.PrinterManager.printJobSem <- job:
	default:
		return nil, errors.New("serial writer busy with antoher job")
	}
	return nil, nil
}
func (r *mutationResolver) AbortPrintJob(ctx context.Context, void *bool) (*bool, error) {
	select {
	case r.App.PrinterManager.abortPrintSem <- true:
	default:
	}
	return nil, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) SerialPorts(ctx context.Context) ([]string, error) {
	return []string{"/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/ttyUSB2", "/dev/ttyUSB3", "/dev/ttyACM0", "/dev/ttyACM1", "/dev/ttyACM2", "/dev/cu.wchusbserial14d10"}, nil
}

func (r *queryResolver) Settings(ctx context.Context) (*Settings, error) {
	set := r.App.GetSettings()
	return &set, nil
}
func (r *queryResolver) TrackedValues(ctx context.Context) (resp []*TrackedValue, err error) {
	resp = []*TrackedValue{}
	for _, tv := range r.App.TrackedValuesManager.TrackedValues {
		tv.ValueMutex.RLock()
		val := *tv.TrackedValue
		resp = append(resp, &val)
		tv.ValueMutex.RUnlock()
	}
	return resp, nil
}
func (r *queryResolver) TrackedValue(ctx context.Context, name string) (*TrackedValue, error) {
	tv, ok := r.App.TrackedValuesManager.TrackedValues[name]
	tv.ValueMutex.RLock()
	defer tv.ValueMutex.RUnlock()
	if !ok {
		return nil, errors.New("tracked value with this name not found")
	}
	val := *tv.TrackedValue
	return &val, nil
}
func (r *queryResolver) ScrollbackBuffer(ctx context.Context) (string, error) {
	r.App.PrinterManager.scrollbackBufferMutex.Lock()
	defer r.App.PrinterManager.scrollbackBufferMutex.Unlock()
	if r.App.PrinterManager.scrollbackBuffer == nil {
		return "", errors.New("No scrollback buffer")
	}
	return r.App.PrinterManager.scrollbackBuffer.String(), nil
}
func (r *queryResolver) RecentCommands(ctx context.Context) ([]string, error) {
	return r.App.RecentCommandsManager.GetRecentCommands(), nil
}
func (r *queryResolver) GcodeFileMetas(ctx context.Context) ([]*GcodeFileMeta, error) {
	metas := []*GcodeFileMeta{}
	metafilePaths, _ := filepath.Glob(filepath.Join(r.App.GetSettings().DataPath, "gcode_files/", "*.gcode.meta"))
	for _, fp := range metafilePaths {
		meta, err := loadGcodeFileMeta(fp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read gcode file meta")
		}
		metas = append(metas, meta)

	}
	sort.Slice(metas, func(i int, j int) bool {
		return !metas[i].UploadDate.Before(metas[j].UploadDate)
	})
	return metas, nil
}
func (r *queryResolver) CurrentPrintJob(ctx context.Context) (*PrintJob, error) {
	panic("not implemented")
}
func (r *queryResolver) SystemInformation(ctx context.Context) (*map[string]interface{}, error) {
	resp := map[string]interface{}{}
	v, _ := mem.VirtualMemory()
	resp["AppName"] = "Biedaprint"
	resp["SystemUsedMemoryPercent"] = fmt.Sprintf("%4.2f%%", v.UsedPercent)
	resp["SystemTotalMemory"] = byteCountBinary(int64(v.Total))
	resp["SystemUsedMemory"] = byteCountBinary(int64(v.Used))
	resp["SystemFreeMemory"] = byteCountBinary(int64(v.Free))
	resp["SystemTime"] = time.Now().Format(time.RFC1123Z)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	resp["AppSysMemory"] = byteCountBinary(int64(m.Sys))
	resp["AppAlloc"] = byteCountBinary(int64(m.Alloc))
	resp["AppNumGC"] = m.NumGC
	resp["GCCPUFractionPercent"] = fmt.Sprintf("%4.2f%%", m.GCCPUFraction*100)
	return &resp, nil
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) TrackedValueUpdated(ctx context.Context, name string) (<-chan interface{}, error) {
	tv, ok := r.App.TrackedValuesManager.TrackedValues[name]
	tv.ValueMutex.RLock()
	defer tv.ValueMutex.RUnlock()
	if !ok {
		return nil, errors.New("tracked value with this name not found")
	}
	outChan := make(chan interface{})
	sub, can := tv.ValueUpdatedBroadcaster.Subscribe()
	go func() {
	Loop:
		for {
			select {
			case data := <-sub:
				outChan <- data
			case <-ctx.Done():
				can()
				break Loop
			}
		}
	}()

	return outChan, nil
}
func (r *subscriptionResolver) CurrentPrintJobUpdated(ctx context.Context) (<-chan *PrintJob, error) {
	panic("not implemented")
}
func (r *subscriptionResolver) SerialConsoleData(ctx context.Context) (<-chan string, error) {
	outChan := make(chan string)
	sub, can := r.App.PrinterManager.consoleBroadcaster.Subscribe()
	go func() {
	Loop:
		for {
			select {
			case data := <-sub:
				if data != nil {
					outChan <- data.(string)
				}
			case <-ctx.Done():
				can()
				break Loop
			}
		}
	}()

	return outChan, nil
}
