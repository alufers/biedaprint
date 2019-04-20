package main

import (
	"bytes"
	"log"

	"go.bug.st/serial.v1"
)

var globalSerial serial.Port
var globalSerialStatus = "disconnected"
var serialReady = make(chan bool)
var lineBuf = bytes.NewBuffer(make([]byte, 512))

func lineParser() {
	// for {
	// 	var line string
	// 	fmt.Fscanln(lineBuf, line)
	// 	log.Println("GOT LINE", line)
	// }
}

//serialReader runs on a separate goroutine and handles broadcasting the serial messages to websockets and saving the data in a backbuffer
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
				scrollbackBuffer.Write(data[:n])
				lineBuf.Write(data[:n])
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
