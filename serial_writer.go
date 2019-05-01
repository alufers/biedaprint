package main

import (
	"log"
)

var serialConsoleWrite = make(chan string, 20)
var serialOkSem = make(chan bool)
var serialResendSem = make(chan int, 10)
var serialPrintJobSem = make(chan *printJob)

//serialWriter runs on a separate goroutine
func serialWriter() {
	for {
		waitForSerialReady()
		log.Print("Serial writer: serial ready")
		for {
			if globalSerial == nil {
				break
			}
			select {
			case c := <-serialConsoleWrite:
				log.Print("Serial writer: serialConsoleWrite", c)
				_, err := globalSerial.Write([]byte(c))
				if err != nil {
					log.Printf("error while writing from serial console to serial: %v", err)
				}
			case job := <-serialPrintJobSem:
				log.Printf("New job %v", job)
				lineChan, err := job.jobLines()
				if err != nil {
					log.Printf("Failed to read job lines: %v", err)
					break
				}
				var sendAndMaybeResend func(string)
				sendAndMaybeResend = func(l string) {
					log.Printf("sending: %#v", l)
					globalSerial.Write([]byte(l))
					select {
					case <-serialOkSem:
					case num := <-serialResendSem:
						sendAndMaybeResend(job.getLineForResend(num))
					}
				}
				for line := range lineChan {
					sendAndMaybeResend(line)
				}
			}
		}
	}
}
