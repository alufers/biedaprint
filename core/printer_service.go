package core

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/armon/circbuf"
	"github.com/pkg/errors"
)

type PrinterService struct {
	app         *App
	printerLink PrinterLink

	serialConnectMutex    *sync.Mutex
	lineBuf               []byte
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
		printerLink:           NewSerialPrinterLink(),
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
	go pm.eventLoop()
	go pm.serialWriterGoroutine()
}

func (pm *PrinterService) ConnectToSerial() error {
	pm.serialConnectMutex.Lock()
	defer pm.serialConnectMutex.Unlock()
	pm.lineBuf = make([]byte, 0)

	settings := pm.app.GetSettings()
	pm.printerLink.(*SerialPrinterLink).SetConfig(SerialPrinterLinkConfigFromSettings(&settings))

	err := pm.printerLink.Connect()
	if err != nil {
		return errors.Wrap(err, "failed to connect to printer")
	}

	pm.scrollbackBufferMutex.Lock()
	defer pm.scrollbackBufferMutex.Unlock()
	pm.scrollbackBuffer, err = circbuf.NewBuffer(int64(settings.ScrollbackBufferSize))
	if err != nil {
		return errors.Wrap(err, "failed to create circbuf")
	}
	return nil
}

func (pm *PrinterService) DisconnectFromSerial() error {
	pm.serialConnectMutex.Lock()
	defer pm.serialConnectMutex.Unlock()
	err := pm.printerLink.Disconnect()
	if err != nil {
		return errors.Wrap(err, "failed to disconnect from printer")
	}
	return nil
}

func (pm *PrinterService) WaitForSerialReady() {
	if pm.printerLink.CurrentStatus() == StatusConnected {
		return
	}
	statusEvent, cancel := pm.printerLink.Status().Subscribe()
	defer cancel()
	for status := range statusEvent {
		if status == StatusConnected {
			return
		}
	}
}

func (pm *PrinterService) eventLoop() {
	statusEvent, _ := pm.printerLink.Status().Subscribe()
	dataEvent, _ := pm.printerLink.Data().Subscribe()
	for {
		select {
		case newStatus := <-statusEvent:
			pm.app.TrackedValuesService.TrackedValues["serialStatus"].UpdateValue(newStatus.(PrinterLinkStatus).TrackedValueString())
		case data := <-dataEvent:
			pm.handleIncomingData(data.([]byte))
		}
	}
}

func (pm *PrinterService) handleIncomingData(data []byte) {
	pm.scrollbackBufferMutex.Lock()
	defer pm.scrollbackBufferMutex.Unlock()
	pm.scrollbackBuffer.Write(data)

	// lineparsing logic
	for i := 0; i < len(data); i++ {
		pm.lineBuf = append(pm.lineBuf, data[i])
		if data[i] == '\n' {
			pm.parseLine(string(pm.lineBuf))
			pm.lineBuf = pm.lineBuf[0:0]
		}
	}
	//serialConsole logic
	strData := string(data)
	pm.consoleBroadcaster.Broadcast(strData)
}

func (pm *PrinterService) serialWriterGoroutine() {
	for {
		pm.WaitForSerialReady()
		log.Print("Serial writer: serial ready")
		for {
			select {
			case c := <-pm.consoleWriteSem:
				log.Print("Serial writer: serialConsoleWrite", c)
				err := pm.printerLink.Write([]byte(c))
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
	var sendAndMaybeResend func(string)
	sendAndMaybeResend = func(l string) {
		pm.printerLink.Write([]byte(l))
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
	} else if strings.HasPrefix(line, "X:") {
		pm.app.ManualMovementService.ProcessPositionReportLine(line)
	}
}
