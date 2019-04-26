package main

import (
	"github.com/chewxy/math32"
	"strconv"
	"strings"
)

var gcodeInstructionConstantPenality float32 = 0.003 // 3 ms

type gcodePrinterStatus struct {
	X float32
	Y float32
	Z float32
	F float32
	E float32
}

func (gps gcodePrinterStatus) nextAfterMovement(segments []string) (next gcodePrinterStatus, time float32, filamentUsed float32, err error) {
	next = gps

	for _, seg := range segments[1:] {
		if strings.HasPrefix(seg, "F") {
			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "F"), 32)
			if err != nil {
				return
			}
			next.F = float32(f64)
		} else if strings.HasPrefix(seg, "X") {

			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "X"), 32)
			if err != nil {
				return
			}
			next.X = float32(f64)
		} else if strings.HasPrefix(seg, "Y") {

			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "Y"), 32)
			if err != nil {
				return
			}
			next.Y = float32(f64)
		} else if strings.HasPrefix(seg, "Z") {

			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "Z"), 32)
			if err != nil {
				return
			}
			next.Z = float32(f64)
		} else if strings.HasPrefix(seg, "E") {

			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "E"), 32)
			if err != nil {
				return
			}
			next.E = float32(f64)
		}
	}
	filamentUsed = next.E - gps.E
	dist := math32.Sqrt((next.X-gps.X)*(next.X-gps.X) + (next.Y-gps.Y)*(next.Y-gps.Y) + (next.Z-gps.Z)*(next.Z-gps.Z))
	if next.F == 0 {
		time = 0
	} else {

		time = dist / next.F * 60
	}
	return
}

type gcodeSimulator struct {
	currentStatus gcodePrinterStatus
	filamentUsed  float32 // milimeters
	time          float32 // seconds
}

func (gs *gcodeSimulator) parseLine(line string) error {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, ";") || line == "" {
		return nil
	}
	segments := strings.Split(line, " ")
	switch segments[0] {
	case "G0":
		fallthrough
	case "G1":
		next, time, filamentUsed, err := gs.currentStatus.nextAfterMovement(segments)
		if err != nil {
			return err
		}
		gs.currentStatus = next
		gs.time += time + gcodeInstructionConstantPenality
		gs.filamentUsed += filamentUsed
	case "G92":
		// http://marlinfw.org/docs/gcode/G092.html Set Position
		next, _, _, err := gs.currentStatus.nextAfterMovement(segments)
		if err != nil {
			return err
		}
		gs.currentStatus = next
	default:
	}
	return nil
}
