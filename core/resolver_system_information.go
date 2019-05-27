package core

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/mem"
)

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
	resp["AppVersion"] = AppVersion
	resp["AppReleaseExecutableName"] = AppReleaseExecutableName
	return &resp, nil
}
