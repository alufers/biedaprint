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
	pm.app.TrackedValuesService.TrackedValues["serialStatus"].UpdateValue("connected")
	return nil
}

func (pm *PrinterService) DisconnectFromSerial() error {
	pm.serialConnectMutex.Lock()
	defer pm.serialConnectMutex.Unlock()
	if !pm.serialReady {
		return errors.New("serial not connected")
	}
	pm.serialReady = false
	pm.app.TrackedValuesService.TrackedValues["serialStatus"].UpdateValue("disconnected")
	return pm.serial.Close()
}

func (pm *PrinterService) WaitForSerialReady() {
	pm.serialReadyCond.L.Lock()
	defer pm.serialReadyCond.L.Unlock()
	for !pm.serialReady {
		pm.serialReadyCond.Wait()
	}
}

func (pm *PrinterService) serialWriterGoroutine() {
	for {
		pm.WaitForSerialReady()
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
				pm.handleJob(job)
			}
		}
	}
}

func (pm *PrinterService) handleJob(job *PrintJobInternal) {
	log.Printf("New job %v", job)
	lineChan, err := job.jobLines()
	if err != nil {
		log.Printf("Failed to read job lines: %v", err)
		return
	}

	log.Printf("Starting smart heating. Hotend target: %v, hotbed target: %v", job.GcodeMeta.HotendTemperature, job.GcodeMeta.HotbedTemperature)
	heatingWaitChan := make(chan bool)
	abortHeatingChan := make(chan bool, 1)
	go func() {
		pm.app.HeatingService.SmartHeatUp(job.GcodeMeta.HotendTemperature, job.GcodeMeta.HotbedTemperature, abortHeatingChan)
		heatingWaitChan <- true
	}()

	select {
	case <-pm.abortPrintSem:
		abortHeatingChan <- true
		job.abort()
		return
	case <-heatingWaitChan:
	}
	log.Printf("Smart heating finished")
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
	for line := range lineChan {
		select {
		case c := <-pm.consoleWriteSem:
			sendAndMaybeResend(c)
			continue
		case <-pm.abortPrintSem:
			job.abort()
			return
		default:
		}
		sendAndMaybeResend(line)
	}
}

func (pm *PrinterService) parseLine(line string) {
	line = strings.TrimSpace(line)
	//log.Println("GOT LINE: ", line)
	if strings.HasPrefix(line, "T:") {
		pm.app.HeatingService.processTemperatureReport(line)
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
		pm.WaitForSerialReady()
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
