package main

import (
	"fmt"
	"log"
	"strings"

	"go.bug.st/serial.v1"
)

var globalSerial serial.Port
var globalSerialStatus = "disconnected"
var serialReady = make(chan bool)

func parseLine(line string) {
	line = strings.TrimSpace(line)
	//log.Println("GOT LINE: ", line)
	if strings.HasPrefix(line, "T:") {
		//trackedValues["hotendTemperature"]
		var temp float64
		var target float64
		var power int
		fmt.Sscanf(line, "T:%f /%f @:%d", &temp, &target, &power)
		trackedValues["hotendTemperature"].updateValue(temp)
		trackedValues["targetHotendTemperature"].updateValue(target)
	}
}

//serialReader runs on a separate goroutine and handles broadcasting the serial messages to websockets and saving the data in a backbuffer
func serialReader() {
	for {
		<-serialReady // wait for somebody to open the serial
		lineBuf := []byte{}
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
				scrollbackBuffer.Write(data[:n])

				for i := 0; i < n; i++ {
					lineBuf = append(lineBuf, data[i])
					if data[i] == '\n' {
						parseLine(string(lineBuf))
						lineBuf = lineBuf[0:0]
					}
				}
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
