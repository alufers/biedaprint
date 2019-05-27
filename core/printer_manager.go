package core

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/armon/circbuf"
	"github.com/pkg/errors"
	"go.bug.st/serial.v1"
)

type PrinterService struct {
	app    *App
	serial serial.Port

	serialConnectMutex    *sync.Mutex
	serialReady           bool
	serialReadyCond       *sync.Cond
	job                   *PrintJobInternal
	okSem                 chan bool
	resendSem             chan int
	abortPrintSem         chan bool
	consoleWriteSem       chan string
	printJobSem           chan *PrintJobInternal
	scrollbackBuffer      *circbuf.Buffer
	scrollbackBufferMutex *sync.RWMutex
	consoleBroadcaster    *EventBroadcaster
}

func NewPrinterService(app *App) *PrinterService {
	return &PrinterService{
		app:                   app,
		serialConnectMutex:    &sync.Mutex{},
		serialReadyCond:       sync.NewCond(&sync.Mutex{}),
		okSem:                 make(chan bool),
		resendSem:             make(chan int, 10),
		abortPrintSem:         make(chan bool, 1),
		consoleWriteSem:       make(chan string, 20),
		printJobSem:           make(chan *PrintJobInternal),
		scrollbackBufferMutex: &sync.RWMutex{},
		consoleBroadcaster:    NewEventBroadcaster(),
	}
}

func (pm *PrinterService) Init() {
	go pm.serialWriterGoroutine()
	go pm.serialReaderGoroutine()
}

func (pm *PrinterService) ConnectToSerial() error {
	pm.serialConnectMutex.Lock()
	defer pm.serialConnectMutex.Unlock()
	parity := serial.EvenParity
	if pm.app.GetSettings().Parity == SerialParityNone {
		parity = serial.NoParity
	}

	ser, err := serial.Open(pm.app.GetSettings().SerialPort, &serial.Mode{
		BaudRate: pm.app.GetSettings().BaudRate,
		Parity:   parity,
		DataBits: pm.app.GetSettings().DataBits,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		return errors.Wrap(err, "failed to connect to serial port")
	}
	pm.serial = ser
	pm.scrollbackBufferMutex.Lock()
	defer pm.scrollbackBufferMutex.Unlock()
	pm.scrollbackBuffer, err = circbuf.NewBuffer(int64(pm.app.GetSettings().ScrollbackBufferSize))
	if err != nil {
		panic(err)
	}
	pm.serialReadyCond.L.Lock()
	defer pm.serialReadyCond.L.Unlock()
	pm.serialReady = true
	pm.serialReadyCond.Broadcast()
	pm.app.TrackedValuesManager.TrackedValues["serialStatus"].UpdateValue("connected")
	return nil
}

func (pm *PrinterService) DisconnectFromSerial() error {
	pm.serialConnectMutex.Lock()
	defer pm.serialConnectMutex.Unlock()
	if !pm.serialReady {
		return errors.New("serial not connected")
	}
	pm.serialReady = false
	pm.app.TrackedValuesManager.TrackedValues["serialStatus"].UpdateValue("disconnected")
	return pm.serial.Close()
}

func (pm *PrinterService) waitForSerialReady() {
	pm.serialReadyCond.L.Lock()
	defer pm.serialReadyCond.L.Unlock()
	for !pm.serialReady {
		pm.serialReadyCond.Wait()
	}
}

func (pm *PrinterService) serialWriterGoroutine() {
	for {
		pm.waitForSerialReady()
		log.Print("Serial writer: serial ready")
		for {
			if pm.serial == nil {
				break
			}
			select {
			case c := <-pm.consoleWriteSem:
				log.Print("Serial writer: serialConsoleWrite", c)
				_, err := pm.serial.Write([]byte(c))
				if err != nil {
					log.Printf("error while writing from serial console to serial: %v", err)
				}
			case job := <-pm.printJobSem:
				log.Printf("New job %v", job)
				lineChan, err := job.jobLines()
				if err != nil {
					log.Printf("Failed to read job lines: %v", err)
					break
				}
				var sendAndMaybeResend func(string)
				sendAndMaybeResend = func(l string) {
					pm.serial.Write([]byte(l))
					select {
					case <-pm.okSem:
					case num := <-pm.resendSem:
						log.Printf("Resending line %v", num)
						<-pm.okSem
						sendAndMaybeResend(job.getLineForResend(num))
					}
				}
			LineLoop:
				for line := range lineChan {
					select {
					case c := <-pm.consoleWriteSem:
						sendAndMaybeResend(c)
						continue
					case <-pm.abortPrintSem:
						job.abort()
						break LineLoop
					default:
					}
					sendAndMaybeResend(line)
				}
			}
		}
	}
}

func (pm *PrinterService) parseLine(line string) {
	line = strings.TrimSpace(line)
	//log.Println("GOT LINE: ", line)
	if strings.HasPrefix(line, "T:") {
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
		pm.app.TrackedValuesManager.TrackedValues["hotendTemperature"].UpdateValue(temp)
		pm.app.TrackedValuesManager.TrackedValues["targetHotendTemperature"].UpdateValue(target)
		pm.app.TrackedValuesManager.TrackedValues["hotbedTemperature"].UpdateValue(bedTemp)
		pm.app.TrackedValuesManager.TrackedValues["targetHotbedTemperature"].UpdateValue(betTargetTemp)
	} else if strings.HasPrefix(line, "ok") {
		select {
		case pm.okSem <- true:
		default:
		}
	} else if strings.HasPrefix(line, "Resend:") {
		var lineNumber int
		fmt.Sscanf(line, "Resend: %d", lineNumber)
		select {
		case pm.resendSem <- lineNumber:
		default:
		}
	}
}

//serialReader runs on a separate goroutine and handles broadcasting the serial messages to websockets and saving the data in a backbuffer
func (pm *PrinterService) serialReaderGoroutine() {
	for {
		pm.waitForSerialReady()
		lineBuf := []byte{}
		for {
			var data = make([]byte, 512)
			n, err := pm.serial.Read(data)
			if err != nil {
				log.Printf("Serial error %v", err)
				//trackedValues["serialStatus"].updateValue("error")
				break
			}
			func() {
				// scrollback buffer logic
				pm.scrollbackBufferMutex.Lock()
				defer pm.scrollbackBufferMutex.Unlock()
				pm.scrollbackBuffer.Write(data[:n])

				// lineparsing logic
				for i := 0; i < n; i++ {
					lineBuf = append(lineBuf, data[i])
					if data[i] == '\n' {
						pm.parseLine(string(lineBuf))
						lineBuf = lineBuf[0:0]
					}
				}
				//serialConsole logic
				strData := string(data[:n])
				pm.consoleBroadcaster.Broadcast(strData)
			}()

		}

	}
}
