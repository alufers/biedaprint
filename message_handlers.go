package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/mem"
	"go.bug.st/serial.v1"
)

type jd map[string]interface{} //json data

var messageHandlers = map[string]func(*websocket.Conn, interface{}){
	"serialList":           handleSerialListMessage,
	"getSettings":          handleGetSettingsMessage,
	"saveSettings":         handleSaveSettingsMessage,
	"getSerialStatus":      handleGetSerialStatusMessage,
	"connectToSerial":      handleConnectToSerialMessage,
	"disconnectFromSerial": handleDisconnectFromSerialMessage,
	"serialWrite":          handleSerialWriteMessage,
	"getSysteminfo":        handleGetSystemInfoMessage,
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
			"ports": []string{"/dev/ttyUSB0", "/dev/ttyUSB1", "/dev/ttyUSB2", "/dev/ttyUSB3"},
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
	globalSerialStatus = "connected"
	select {
	case serialReady <- true:
	default:
	}

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

	_, err := globalSerial.Write([]byte((data.(map[string]interface{}))["data"].(string)))
	if err != nil {
		sendError(c, errors.Wrap(err, "failed to write to serial"))

		return
	}
}

func handleGetSystemInfoMessage(c *websocket.Conn, data interface{}) {
	var resp = jd{}
	v, _ := mem.VirtualMemory()
	resp["AppName"] = "Biedaprint"
	resp["SystemUsedMemoryPercent"] = fmt.Sprintf("%4.2f%%", v.UsedPercent)
	resp["SystemTotalMemory"] = byteCountBinary(int64(v.Total))
	resp["SystemUsedMemory"] = byteCountBinary(int64(v.Used))
	resp["SystemFreeMemory"] = byteCountBinary(int64(v.Free))
	// resp["FreeRamPercent"] = int((float64(sys.FreeRam) / float64(sys.TotalRam)) * 100.0)
	// resp["FreeRamPercentNoBuffer"] = int((float64(sys.FreeRam+sys.BufferRam) / float64(sys.TotalRam)) * 100.0)
	c.WriteJSON(jd{
		"type": "getSysteminfo",
		"data": resp,
	})
}
