package core

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/armon/circbuf"
	"github.com/pkg/errors"
)

/*
PrinterService handl;es communication with the printer using a PrinterLink.
It manages sending various commands to the printer, recieving messages from the printer and dispatching them to various parts of the application.
*/
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

/*
NewPrinterService constructs a PrinterService.
*/
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

/*
Init initializes the printer service by starting all the necessary goroutines.
*/
func (pm *PrinterService) Init() error {
	go pm.eventLoop()
	go pm.serialWriterGoroutine()
	return nil
}

/*
ConnectToSerial connects to the printer using the PrinterLink.
*/
func (pm *PrinterService) ConnectToSerial() error {
	pm.serialConnectMutex.Lock()
	defer pm.serialConnectMutex.Unlock()
	pm.lineBuf = make([]byte, 0)

	pm.printerLink.(*SerialPrinterLink).SetConfig(SerialPrinterLinkConfigFromSettings(pm.app))

	err := pm.printerLink.Connect()
	if err != nil {
		return errors.Wrap(err, "failed to connect to printer")
	}

	pm.scrollbackBufferMutex.Lock()
	defer pm.scrollbackBufferMutex.Unlock()
	scrollbackBufferSize, err := pm.app.SettingsService.GetInt64("serial.scrollbackBufferSize")
	if err != nil {
		return errors.Wrap(err, "failed to create circbuf")
	}
	pm.scrollbackBuffer, err = circbuf.NewBuffer(scrollbackBufferSize)
	if err != nil {
		return errors.Wrap(err, "failed to create circbuf")
	}
	return nil
}

/*
DisconnectFromSerial disconnects from the printer.
*/
func (pm *PrinterService) DisconnectFromSerial() error {
	pm.serialConnectMutex.Lock()
	defer pm.serialConnectMutex.Unlock()
	err := pm.printerLink.Disconnect()
	if err != nil {
		return errors.Wrap(err, "failed to disconnect from printer")
	}
	return nil
}

/*
WaitForSerialReady blocks until the app connects to the printer.
*/
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

/*
eventLoop subscribes to broadcasts from the printer service and handles them
*/
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
				err := pm.sendLineAndSniff(c)
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
		pm.sendLineAndSniff(l)
		select {
		case <-pm.okSem:
		case num := <-pm.resendSem:
			log.Printf("Resending line %v", num)
			<-pm.okSem
			sendAndMaybeResend(job.getLineForResend(num))
		case <-pm.abortPrintSem:
			job.abort()
			return
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

func (pm *PrinterService) sendLineAndSniff(line string) error {
	fanSpeedRegex := regexp.MustCompile("(L[0-9\\ ]+)?M106 S([0-9]*)")
	matches := fanSpeedRegex.FindSubmatch([]byte(line))
	log.Printf("%v, len(matches) = %v", line, len(matches))
	if len(matches) == 3 {
		log.Printf("FOUDN FANCTL")
		numValue, err := strconv.Atoi(string(matches[2]))
		if err == nil {
			pm.app.TrackedValuesService.TrackedValues["fanSpeed"].UpdateValue(float64(numValue))
		}
	}
	return pm.printerLink.Write([]byte(line))
}

/*
parseLine processes lines incoming from the printer's firmware and dispatches them.
*/
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
