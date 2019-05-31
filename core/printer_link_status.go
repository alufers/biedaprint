package core

type PrinterLinkStatus int

const (
	StatusDisconnected PrinterLinkStatus = iota
	StatusConnected
	StatusError
)
