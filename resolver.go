package biedaprint

import (
	"context"
	"fmt"
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
	defer r.App.settingsMutex.Unlock()
	r.App.settings.SerialPort = settings.SerialPort
	r.App.settings.BaudRate = settings.BaudRate
	r.App.settings.ScrollbackBufferSize = settings.ScrollbackBufferSize
	r.App.settings.DataPath = settings.DataPath
	r.App.settings.StartupCommand = settings.StartupCommand
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
func (r *mutationResolver) UploadGcode(ctx context.Context, file graphql.Upload) (*bool, error) {
	panic("not implemented")
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
	panic("not implemented")
}
func (r *mutationResolver) AbortPrintJob(ctx context.Context, void *bool) (*bool, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Settings(ctx context.Context) (*Settings, error) {
	set := r.App.GetSettings()
	return &set, nil
}
func (r *queryResolver) TrackedValues(ctx context.Context) ([]*TrackedValue, error) {
	panic("not implemented")
}
func (r *queryResolver) TrackedValue(ctx context.Context, name string) (*TrackedValue, error) {
	panic("not implemented")
}
func (r *queryResolver) ScrollbackBuffer(ctx context.Context) (string, error) {
	r.App.PrinterManager.scrollbackBufferMutex.Lock()
	defer r.App.PrinterManager.scrollbackBufferMutex.Unlock()
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

func (r *subscriptionResolver) TrackedValueUpdated(ctx context.Context, name string) (<-chan *TrackedValue, error) {
	panic("not implemented")
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
