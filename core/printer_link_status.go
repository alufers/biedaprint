package core

type PrinterLinkStatus int

const (
	StatusDisconnected PrinterLinkStatus = iota
	StatusConnected
	StatusError
)

func (pls PrinterLinkStatus) TrackedValueString() string {
	switch pls {
	case StatusConnected:
		return "connected"
	case StatusDisconnected:
		return "disconnected"
	case StatusError:
		return "error"
	default:
		panic("invalid PrinterLinkStatus")
	}
}
