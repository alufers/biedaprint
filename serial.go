package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

var globalSerial io.ReadWriteCloser
var globalSerialStatus = "disconnected"
var serialReady bool
var serialReadySems = []chan bool{}
var serialConsoleSubscribers = []*websocket.Conn{}

func waitForSerialReady() {
	if serialReady {
		return
	}
	c := make(chan bool)
	serialReadySems = append(serialReadySems, c)
	<-c
	for i, cc := range serialReadySems {
		if cc == c {
			serialReadySems = append(serialReadySems[:i], serialReadySems[i+1:]...)
		}
	}
}

func parseLine(line string) {
	line = strings.TrimSpace(line)
	//log.Println("GOT LINE: ", line)
	if strings.HasPrefix(line, "T:") {
		//trackedValues["hotendTemperature"]
		var temp float64
		var target float64
		var power int
		var bedTemp float64
		var betTargetTemp float64
		var bedPower float64
		if strings.Contains(line, "B:") { // has heated bed
			fmt.Sscanf(line, "T:%f /%f B:%f /%f @:%d B@:%d", &temp, &target, &bedTemp, &betTargetTemp, &power, &bedPower)
		} else {
			fmt.Sscanf(line, "T:%f /%f @:%d", &temp, &target, &power)
		}
		trackedValues["hotendTemperature"].updateValue(temp)
		trackedValues["targetHotendTemperature"].updateValue(target)
		trackedValues["hotbedTemperature"].updateValue(bedTemp)
		trackedValues["targetHotbedTemperature"].updateValue(betTargetTemp)
	} else if strings.HasPrefix(line, "ok") {
		select {
		case serialOkSem <- true:
		default:
		}
	} else if strings.HasPrefix(line, "Resend:") {
		var lineNumber int
		fmt.Sscanf(line, "Resend: %d", lineNumber)
		select {
		case serialResendSem <- lineNumber:
		default:
		}
	}
}

//serialReader runs on a separate goroutine and handles broadcasting the serial messages to websockets and saving the data in a backbuffer
func serialReader() {
	for {
		waitForSerialReady()
		lineBuf := []byte{}
		for {
			if globalSerial == nil {
				break
			}
			var data = make([]byte, 512)
			n, err := globalSerial.Read(data)
			if err != nil {
				log.Printf("Serial error %v", err)
				trackedValues["serialStatus"].updateValue("error")
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
				for _, a := range serialConsoleSubscribers {
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
