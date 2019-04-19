package main

import (
	"log"

	"go.bug.st/serial.v1"
)

var globalSerial serial.Port
var globalSerialStatus = "disconnected"
var serialReady = make(chan bool)

//serialReader runs on a separate goroutine a
func serialReader() {
	for {
		<-serialReady // wait for somebody to open the serial
		for {
			if globalSerial == nil {
				break
			}
			var data = make([]byte, 512)
			n, err := globalSerial.Read(data)
			if err != nil {
				log.Printf("Serial error %v", err)
				globalSerialStatus = "error"
				break
			}
			func() {
				handlerMutex.Lock()
				defer handlerMutex.Unlock()
				strData := string(data[:n])
				for _, a := range activeConnections {
					a.WriteJSON(jd{
						"type": "serialConsole",
						"data": jd{
							"data": strData,
						},
					})
				}
			}()

		}

	}
}
