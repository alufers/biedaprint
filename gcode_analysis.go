package main

import (
	"math"
	"strconv"
	"strings"
)

var gcodeInstructionConstantPenality = 0.003 // 3 ms

type gcodePrinterStatus struct {
	X float64
	Y float64
	Z float64
	F float64
	E float64
}

func (gps gcodePrinterStatus) nextAfterMovement(segments []string) (next gcodePrinterStatus, time float64, filamentUsed float64, err error) {
	next = gps

	for _, seg := range segments[1:] {
		if strings.HasPrefix(seg, "F") {
			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "F"), 32)
			if err != nil {
				return
			}
			next.F = float64(f64)
		} else if strings.HasPrefix(seg, "X") {

			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "X"), 32)
			if err != nil {
				return
			}
			next.X = float64(f64)
		} else if strings.HasPrefix(seg, "Y") {

			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "Y"), 32)
			if err != nil {
				return
			}
			next.Y = float64(f64)
		} else if strings.HasPrefix(seg, "Z") {

			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "Z"), 32)
			if err != nil {
				return
			}
			next.Z = float64(f64)
		} else if strings.HasPrefix(seg, "E") {

			var f64 float64
			f64, err = strconv.ParseFloat(strings.TrimPrefix(seg, "E"), 32)
			if err != nil {
				return
			}
			next.E = float64(f64)
		}
	}
	filamentUsed = next.E - gps.E
	dist := math.Sqrt((next.X-gps.X)*(next.X-gps.X) + (next.Y-gps.Y)*(next.Y-gps.Y) + (next.Z-gps.Z)*(next.Z-gps.Z))
	if next.F == 0 {
		time = 0
	} else {

		time = dist / next.F * 60
	}
	return
}

type gcodeSimulator struct {
	currentStatus gcodePrinterStatus
	filamentUsed  float64 // milimeters
	time          float64 // seconds
	layerIndexes  []gcodeLayerIndex
	layer         int
}

func (gs *gcodeSimulator) parseLine(line string, number int) error {
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
		if gs.currentStatus.Z < next.Z { // next layer
			gs.layerIndexes = append(gs.layerIndexes, gcodeLayerIndex{
				LineNumber:  number,
				LayerNumber: gs.layer,
			})
			gs.layer++
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
