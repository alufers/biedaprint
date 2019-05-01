package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/mem"
	"go.bug.st/serial.v1"
)

type jd map[string]interface{} //json data

var messageHandlers = map[string]func(*websocket.Conn, interface{}){
	"serialList":              handleSerialListMessage,
	"getSettings":             handleGetSettingsMessage,
	"saveSettings":            handleSaveSettingsMessage,
	"getSerialStatus":         handleGetSerialStatusMessage,
	"connectToSerial":         handleConnectToSerialMessage,
	"disconnectFromSerial":    handleDisconnectFromSerialMessage,
	"serialWrite":             handleSerialWriteMessage,
	"getSysteminfo":           handleGetSystemInfoMessage,
	"sendGCODE":               handleSendGCODEMessage,
	"getScrollbackBuffer":     handleGetScrollbackBufferMessage,
	"getTrackedValues":        handleGetTrackedValuesMessage,
	"getTrackedValue":         handleGetTrackedValueMessage,
	"subscribeToTrackedValue": handleSubscribeToTrackedValueMessage,
	"getGcodeFileMetas":       handleGetGcodeFileMetas,
	"deleteGcodeFile":         handleDeleteGcodeFile,
	"startPrintJob":           handleStartPrintJob,
}

func sendError(c *websocket.Conn, err error) {
	log.Printf("Handler error %v", err)
	c.WriteJSON(jd{
		"type": "alert",
		"data": jd{
			"type":    "danger",
			"content": err.Error(),
		},
	})
}

func handleSerialListMessage(c *websocket.Conn, data interface{}) {
	// ports, err := serial.GetPortsList()
	// if err != nil {
	// 	log.Print(err)
	// 	sendError(c, err)
	// 	return
	// }
	c.WriteJSON(jd{
		"type": "serialList",
		"data": jd{
			"ports": []string{"/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/ttyUSB2", "/dev/ttyUSB3", "/dev/ttyACM0", "/dev/ttyACM1"},
		},
	})
}

func handleGetSettingsMessage(c *websocket.Conn, data interface{}) {

	c.WriteJSON(jd{
		"type": "getSettings",
		"data": globalSettings,
	})
}

func handleSaveSettingsMessage(c *websocket.Conn, data interface{}) {
	err := mapstructure.Decode(data, &globalSettings)
	if err != nil {
		sendError(c, errors.Wrap(err, "failed to decode settings"))
		return
	}
	err = saveSettings()
	if err != nil {
		sendError(c, errors.Wrap(err, "failed to save settings"))
		return
	}
	c.WriteJSON(jd{
		"type": "alert",
		"data": jd{
			"type":    "success",
			"content": "Settings saved!",
		},
	})
}

func handleGetSerialStatusMessage(c *websocket.Conn, data interface{}) {

	c.WriteJSON(jd{
		"type": "getSerialStatus",
		"data": jd{
			"status": globalSerialStatus,
		},
	})
}

func handleConnectToSerialMessage(c *websocket.Conn, data interface{}) {

	var err error
	globalSerial, err = serial.Open(globalSettings.SerialPort, &serial.Mode{
		BaudRate: globalSettings.BaudRate,
		Parity:   serial.EvenParity,
		DataBits: 7,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		sendError(c, errors.Wrap(err, "failed to connect to printer"))
		globalSerialStatus = "error"
		return
	}
	resetScrollback()
	globalSerialStatus = "connected"
	serialReady = true
	for _, cc := range serialReadySems {
		cc <- true
	}

	globalSerial.Write([]byte("M155 S1\r\n"))

	c.WriteJSON(jd{
		"type": "getSerialStatus",
		"data": jd{
			"status": globalSerialStatus,
		},
	})
}

func handleDisconnectFromSerialMessage(c *websocket.Conn, data interface{}) {
	if globalSerial == nil {
		sendError(c, errors.New("Not connected to serial port"))
		return
	}

	err := globalSerial.Close()
	globalSerialStatus = "disconnected"
	serialReady = false
	if err != nil {
		sendError(c, err)

		return
	}
	c.WriteJSON(jd{
		"type": "getSerialStatus",
		"data": jd{
			"status": globalSerialStatus,
		},
	})
}

func handleSerialWriteMessage(c *websocket.Conn, data interface{}) {
	if globalSerial == nil {
		sendError(c, errors.New("Not connected to serial port"))
		return
	}

	// _, err := globalSerial.Write([]byte((data.(map[string]interface{}))["data"].(string)))
	// if err != nil {
	// sendError(c, errors.Wrap(err, "failed to write to serial"))
	// return
	// }
	serialConsoleWrite <- (data.(map[string]interface{}))["data"].(string)
}

func handleSendGCODEMessage(c *websocket.Conn, data interface{}) {
	if globalSerial == nil {
		sendError(c, errors.New("Not connected to serial port"))
		return
	}
	dataMap := data.(map[string]interface{})
	gcodeStr := dataMap["data"].(string)
	//_, err := globalSerial.Write([]byte((gcodeStr + "\r\n")))
	serialConsoleWrite <- gcodeStr + "\r\n"
	// if err != nil {
	// 	sendError(c, errors.Wrap(err, "failed to write to serial"))
	// 	return
	// }
}

func handleGetSystemInfoMessage(c *websocket.Conn, data interface{}) {
	var resp = jd{}
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
	// resp["FreeRamPercent"] = int((float64(sys.FreeRam) / float64(sys.TotalRam)) * 100.0)
	// resp["FreeRamPercentNoBuffer"] = int((float64(sys.FreeRam+sys.BufferRam) / float64(sys.TotalRam)) * 100.0)
	c.WriteJSON(jd{
		"type": "getSysteminfo",
		"data": resp,
	})
}

func handleGetScrollbackBufferMessage(c *websocket.Conn, data interface{}) {
	if scrollbackBuffer == nil {
		sendError(c, errors.New("No scrollback buffer"))
		return
	}
	c.WriteJSON(jd{
		"type": "getScrollbackBuffer",
		"data": jd{
			"data": string(scrollbackBuffer.Bytes()),
		},
	})
}

func handleGetTrackedValuesMessage(c *websocket.Conn, data interface{}) {
	c.WriteJSON(jd{
		"type": "getTrackedValues",
		"data": jd{
			"trackedValues": trackedValues,
		},
	})
}

func handleGetTrackedValueMessage(c *websocket.Conn, data interface{}) {
	dataMap := data.(map[string]interface{})
	t, _ := trackedValues[dataMap["name"].(string)]
	c.WriteJSON(jd{
		"type": "getTrackedValue",
		"data": jd{
			"trackedValue": t,
		},
	})
}

func handleSubscribeToTrackedValueMessage(c *websocket.Conn, data interface{}) {
	dataMap := data.(map[string]interface{})
	t, ok := trackedValues[dataMap["name"].(string)]
	if !ok {
		sendError(c, errors.New("No such trackedValue"))
		return
	}
	log.Printf("Client subscribed to %v", dataMap["name"].(string))
	for _, s := range t.subscribers {
		if s == c {
			return // duplicate
		}
	}
	t.subscribers = append(t.subscribers, c)
}

func handleGetGcodeFileMetas(c *websocket.Conn, data interface{}) {
	handlerMutex.Unlock()
	metas := []*gcodeFileMeta{}
	metafilePaths, _ := filepath.Glob(filepath.Join(globalSettings.DataPath, "gcode_files/", "*.gcode.meta"))
	for _, fp := range metafilePaths {
		meta, err := loadGcodeFileMeta(fp)
		if err != nil {
			// sendError(c, err)
			// don't abort
		} else {
			metas = append(metas, meta)
		}
	}
	sort.Slice(metas, func(i int, j int) bool {
		return !metas[i].UploadDate.Before(metas[j].UploadDate)
	})
	handlerMutex.Lock()

	c.WriteJSON(jd{
		"type": "getGcodeFileMetas",
		"data": metas,
	})

}

func handleDeleteGcodeFile(c *websocket.Conn, data interface{}) {
	dataMap := data.(map[string]interface{})
	gcodeFileName := dataMap["gcodeFileName"].(string)
	gcodeName := filepath.Join(globalSettings.DataPath, "gcode_files/", gcodeFileName)
	gcodeMetaName := filepath.Join(globalSettings.DataPath, "gcode_files/", gcodeFileName+".meta")
	err := os.Remove(gcodeName)
	if err != nil {
		sendError(c, err)
		return
	}
	err = os.Remove(gcodeMetaName)
	if err != nil {
		sendError(c, err)
		return
	}
	c.WriteJSON(jd{
		"type": "alert",
		"data": jd{
			"type":    "success",
			"content": "Gcode file deleted!",
		},
	})
}

func handleStartPrintJob(c *websocket.Conn, data interface{}) {
	dataMap := data.(map[string]interface{})
	gcodeFileName := dataMap["gcodeFileName"].(string)
	meta, err := loadGcodeFileMeta(filepath.Join(globalSettings.DataPath, "gcode_files/", gcodeFileName+".meta"))
	if err != nil {
		sendError(c, err)
		return
	}
	job := &printJob{
		gcodeMeta: meta,
	}

	select {
	case serialPrintJobSem <- job:
	default:
		sendError(c, errors.New("serial writer busy with antoher job"))
	}

}
