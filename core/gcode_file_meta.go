package core

import (
	"bufio"
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

func init() {
	gob.Register(GcodeFileMeta{})
	gob.Register(GcodeLayerIndex{})
}

//GcodeFileMeta stores information about a gcode file
type GcodeFileMeta struct {
	Model
	OriginalName      string             `json:"originalName"`
	GcodeFileName     string             `json:"gcodeFileName"`
	UploadDate        time.Time          `json:"uploadDate"`
	TotalLines        int                `json:"totalLines"`
	PrintTime         float64            `json:"printTime"`
	FilamentUsedMm    float64            `json:"filamentUsedMm"`
	LayerIndexes      []*GcodeLayerIndex `json:"layerIndexes"`
	HotendTemperature float64            `json:"hotendTemperature"`
	HotbedTemperature float64            `json:"hotbedTemperature"`
}

func loadGcodeFileMeta(path string) (*GcodeFileMeta, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	deco := gob.NewDecoder(f)
	var out *GcodeFileMeta
	err = deco.Decode(&out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (gfm *GcodeFileMeta) Save(dataPath string) error {
	f, err := os.OpenFile(filepath.Join(dataPath, "gcode_files/", gfm.GcodeFileName+".meta"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to open gcode meta file for writing")
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(gfm)
	if err != nil {
		return errors.Wrap(err, "failed to encode gcode meta file")
	}
	return nil
}

func (gfm *GcodeFileMeta) AnalyzeGcodeFile(dataPath string) error {
	f, err := os.Open(filepath.Join(dataPath, "gcode_files/", gfm.GcodeFileName))
	if err != nil {
		return err
	}
	startTime := time.Now()
	s := bufio.NewScanner(f)
	gfm.TotalLines = 0
	sim := &gcodeSimulator{
		layerIndexes: []GcodeLayerIndex{},
	}
	lineNumber := 0
	for s.Scan() {
		line := s.Text()
		gfm.TotalLines++
		lineNumber++
		sim.parseLine(line, lineNumber)
	}
	gfm.PrintTime = sim.time
	gfm.FilamentUsedMm = sim.filamentUsed
	gfm.LayerIndexes = []*GcodeLayerIndex{}
	gfm.HotendTemperature = sim.hotendTemperature
	gfm.HotbedTemperature = sim.hotbedTemperature
	for _, li := range sim.layerIndexes {
		gfm.LayerIndexes = append(gfm.LayerIndexes, &li)
	}
	log.Printf("Finished gcode analysis in %v seconds", time.Now().Sub(startTime).Seconds())
	runtime.GC()
	return nil
}
