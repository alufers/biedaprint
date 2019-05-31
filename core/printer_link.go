package core

type PrinterLink interface {
	Connect() error
	Disconnect() error
	Write([]byte) error
	Data() *EventBroadcaster
	CurrentStatus() PrinterLinkStatus
	Status() *EventBroadcaster
}
