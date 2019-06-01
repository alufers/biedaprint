package core

import (
	"io"
	"log"
	"sync"

	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
)

type SerialPrinterLinkConfig struct {
	SerialPort string
	SerialMode *serial.OpenOptions
}

func SerialPrinterLinkConfigFromSettings(settings *Settings) *SerialPrinterLinkConfig {
	parity := serial.PARITY_EVEN
	if settings.Parity == SerialParityNone {
		parity = serial.PARITY_NONE
	}
	return &SerialPrinterLinkConfig{
		SerialPort: settings.SerialPort,
		SerialMode: &serial.OpenOptions{
			PortName:        settings.SerialPort,
			BaudRate:        uint(settings.BaudRate),
			ParityMode:      parity,
			DataBits:        uint(settings.DataBits),
			StopBits:        1,
			MinimumReadSize: 1,
		},
	}
}

type SerialPrinterLink struct {
	config            *SerialPrinterLinkConfig
	connection        io.ReadWriteCloser
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
	spl.connection, err = serial.Open(*spl.config.SerialMode)
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
		if err != nil {

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
