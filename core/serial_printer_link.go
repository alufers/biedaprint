package core

import (
	"log"
	"sync"

	"github.com/pkg/errors"
	"go.bug.st/serial.v1"
)

type SerialPrinterLinkConfig struct {
	SerialPort string
	SerialMode *serial.Mode
}

func SerialPrinterLinkConfigFromSettings(settings *Settings) *SerialPrinterLinkConfig {
	parity := serial.EvenParity
	if settings.Parity == SerialParityNone {
		parity = serial.NoParity
	}
	return &SerialPrinterLinkConfig{
		SerialPort: settings.SerialPort,
		SerialMode: &serial.Mode{
			BaudRate: settings.BaudRate,
			Parity:   parity,
			DataBits: settings.DataBits,
			StopBits: serial.OneStopBit,
		},
	}
}

type SerialPrinterLink struct {
	config            *SerialPrinterLinkConfig
	connection        serial.Port
	statusMutex       *sync.RWMutex
	status            PrinterLinkStatus
	statusBroadcaster *EventBroadcaster
	dataBroadcaster   *EventBroadcaster
}

func NewSerialPrinterLink() *SerialPrinterLink {
	return &SerialPrinterLink{
		statusBroadcaster: NewEventBroadcaster(),
		dataBroadcaster:   NewEventBroadcaster(),
		statusMutex:       &sync.RWMutex{},
	}
}

func (spl *SerialPrinterLink) SetConfig(config *SerialPrinterLinkConfig) {
	spl.config = config
}

func (spl *SerialPrinterLink) Connect() error {
	var err error
	log.Printf("connecting with serial %#v", spl.config.SerialMode)
	spl.connection, err = serial.Open(spl.config.SerialPort, spl.config.SerialMode)
	if err != nil {
		return errors.Wrapf(err, "failed to open SerialPrinterLink at %v", spl.config.SerialPort)
	}
	spl.updateStatus(StatusConnected)
	go spl.readerRoutine()
	return nil
}

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
		log.Printf("broadcasting %v", data[:n])
		if err != nil {
			if portErr, ok := err.(serial.PortError); ok {
				if portErr.Code() == serial.PortClosed {
					return
				}
			}
			log.Printf("SerialPrinterLink.readerRoutine error: %v", err)
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

func (spl *SerialPrinterLink) Data() *EventBroadcaster {
	return spl.dataBroadcaster
}
func (spl *SerialPrinterLink) CurrentStatus() PrinterLinkStatus {
	spl.statusMutex.RLock()
	defer spl.statusMutex.RUnlock()
	return spl.status
}
func (spl *SerialPrinterLink) Status() *EventBroadcaster {
	return spl.statusBroadcaster
}
