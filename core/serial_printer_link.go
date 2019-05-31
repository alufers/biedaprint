package core

import (
	"go.bug.st/serial.v1"
)

type SerialPrinterLink struct {
	connection        serial.Port
	status            PrinterLinkStatus
	statusBroadcaster *EventBroadcaster
	dataBroadcaster   *EventBroadcaster
}

func NewSerialPrinterLink() *SerialPrinterLink {
	return &SerialPrinterLink{
		statusBroadcaster: NewEventBroadcaster(),
		dataBroadcaster:   NewEventBroadcaster(),
	}
}

func (spl *SerialPrinterLink) Connect() error {
	return nil
}
func (spl *SerialPrinterLink) Disconnect() error {
	return nil
}
func (spl *SerialPrinterLink) Write([]byte) error {
	return nil
}
func (spl *SerialPrinterLink) Data() *EventBroadcaster {
	return spl.dataBroadcaster
}
func (spl *SerialPrinterLink) CurrentStatus() PrinterLinkStatus {
	return spl.status
}
func (spl *SerialPrinterLink) Status() *EventBroadcaster {
	return spl.statusBroadcaster
}
