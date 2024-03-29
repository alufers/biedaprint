package core

import (
	"io"
	"log"
	"os"
	"sync"

	errorsNative "errors"

	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
)

/*
SerialPrinterLinkConfig holds the data required to connect to the printer.
*/
type SerialPrinterLinkConfig struct {
	SerialPort string
	SerialMode *serial.OpenOptions
}

/*
SerialPrinterLinkConfigFromSettings reads the settings from the app and creates a SerialPrinterLinkConfig from it.
*/
func SerialPrinterLinkConfigFromSettings(app *App) *SerialPrinterLinkConfig {
	var parity serial.ParityMode
	paritySetting, err := app.SettingsService.GetString("serial.parity")
	if err != nil {
		panic(err)
	}
	if paritySetting == "NONE" {
		parity = serial.PARITY_NONE
	} else if paritySetting == "EVEN" {
		parity = serial.PARITY_EVEN
	} else if paritySetting == "EVEN" {
		parity = serial.PARITY_EVEN
	}

	serialPort, err := app.SettingsService.GetString("serial.serialPort")
	if err != nil {
		panic(err)
	}
	baudRate, err := app.SettingsService.GetUint("serial.baudRate")
	if err != nil {
		panic(err)
	}
	dataBits, err := app.SettingsService.GetUint("serial.dataBits")
	if err != nil {
		panic(err)
	}
	return &SerialPrinterLinkConfig{
		SerialPort: serialPort,
		SerialMode: &serial.OpenOptions{
			PortName:        serialPort,
			BaudRate:        baudRate,
			ParityMode:      parity,
			DataBits:        dataBits,
			StopBits:        1,
			MinimumReadSize: 1,
		},
	}
}

/*
SerialPrinterLink manages communication with the printer using a serial port. It broadcasts the incoming data using an EventBroadcaster.
*/
type SerialPrinterLink struct {
	config            *SerialPrinterLinkConfig
	connection        io.ReadWriteCloser
	statusMutex       *sync.RWMutex
	status            PrinterLinkStatus
	statusBroadcaster *EventBroadcaster
	dataBroadcaster   *EventBroadcaster
}

/*
NewSerialPrinterLink constructs a SerialPrinterLink.
*/
func NewSerialPrinterLink() *SerialPrinterLink {
	return &SerialPrinterLink{
		statusBroadcaster: NewEventBroadcaster(),
		dataBroadcaster:   NewEventBroadcaster(),
		statusMutex:       &sync.RWMutex{},
	}
}

/*
SetConfig sets the Link's configuration.
*/
func (spl *SerialPrinterLink) SetConfig(config *SerialPrinterLinkConfig) {
	spl.config = config
}

/*
Connect connects to the priter using the configuration stored in the struct.
*/
func (spl *SerialPrinterLink) Connect() error {
	var err error
	log.Printf("connecting with serial %#v", spl.config.SerialMode)
	spl.connection, err = serial.Open(*spl.config.SerialMode)
	if err != nil {
		return errors.Wrapf(err, "failed to open SerialPrinterLink at %v", spl.config.SerialPort)
	}
	spl.updateStatus(StatusConnected)
	go spl.readerRoutine()
	return nil
}

/*
Disconnect disconnects from the printer.
*/
func (spl *SerialPrinterLink) Disconnect() error {
	err := spl.connection.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close SerialPrinterLink")
	}
	spl.updateStatus(StatusDisconnected)
	return nil
}

func (spl *SerialPrinterLink) updateStatus(sta PrinterLinkStatus) {
	spl.statusMutex.Lock()
	defer spl.statusMutex.Unlock()
	spl.status = sta
	spl.statusBroadcaster.Broadcast(sta)
}

func (spl *SerialPrinterLink) readerRoutine() {

	for {
		data := make([]byte, 512)
		n, err := spl.connection.Read(data)
		if err != nil {

			log.Printf("SerialPrinterLink.readerRoutine error: %v", err)
			if spl.status == StatusDisconnected && errorsNative.Is(err, os.ErrClosed) {
				return // we don't want to report an error if the connection was closed by the user.
			}
			spl.updateStatus(StatusError)
			return
		}

		spl.dataBroadcaster.Broadcast(data[:n])
	}
}

func (spl *SerialPrinterLink) Write(data []byte) error {
	_, err := spl.connection.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write to SerialPrinterLink")
	}
	return nil
}

/*
Data returns the EventBroadcaster that is fired everytime some data is recieved from the printer.
*/
func (spl *SerialPrinterLink) Data() *EventBroadcaster {
	return spl.dataBroadcaster
}

/*
CurrentStatus returns the status of the link. It is safe to be used by many goroutines.
*/
func (spl *SerialPrinterLink) CurrentStatus() PrinterLinkStatus {
	spl.statusMutex.RLock()
	defer spl.statusMutex.RUnlock()
	return spl.status
}

/*
Status returns an EventBroadcaster that is fired everytime the status of the printer changes.
*/
func (spl *SerialPrinterLink) Status() *EventBroadcaster {
	return spl.statusBroadcaster
}
