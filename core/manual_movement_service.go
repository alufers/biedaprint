package core

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

type PositionVector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	E float64 `json:"e"`
}

func (pv PositionVector) Add(other PositionVector) PositionVector {
	return PositionVector{
		X: pv.X + other.X,
		Y: pv.Y + other.Y,
		Z: pv.Z + other.Z,
		E: pv.E + other.E,
	}
}

type ManualMovementService struct {
	app                *App
	positionReportChan chan string
}

func NewManualMovementService(app *App) *ManualMovementService {
	return &ManualMovementService{
		app:                app,
		positionReportChan: make(chan string),
	}
}

func (mms *ManualMovementService) ProcessPositionReportLine(line string) {
	log.Printf("Recieved position report line in manual movement service")
	select {
	case mms.positionReportChan <- line:
	default:
		log.Print("Ignoring position report line, no manual uperation in progress")
	}
}

// X:10.00 Y:0.00 Z:190.00 E:0.00 Count X:1600 Y:0 Z:30400
func (mms *ManualMovementService) parsePositionReportLine(line string) (pos *PositionVector) {
	pos = &PositionVector{}
	fmt.Sscanf(line, "X:%f Y:%f Z:%f E:%f", &pos.X, &pos.Y, &pos.Z, &pos.E)
	return
}

func (mms *ManualMovementService) MoveRelative(vec *PositionVector) error {
	mms.app.PrinterService.consoleWriteSem <- "M114\r\n"
	select {
	case reportLine := <-mms.positionReportChan:
		currentPosition := mms.parsePositionReportLine(reportLine)
		newPosition := currentPosition.Add(*vec)
		mms.app.PrinterService.consoleWriteSem <- fmt.Sprintf("G0 X%f Y%f Z%f E%f\r\n", newPosition.X, newPosition.Y, newPosition.Z, newPosition.E)
	case <-time.After(time.Second * 5):
		return errors.New("position report request timed out. Is the printer busy?")
	}

	return nil
}
