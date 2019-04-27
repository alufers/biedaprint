package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type gcodeLineIndex struct {
	LineNumber int   `json:"lineNumber"`
	Offset     int64 `json:"offset"`
}

type gcodeFileMeta struct {
	OriginalName  string    `json:"originalName"`
	GcodeFileName string    `json:"gcodeFileName"`
	UploadDate    time.Time `json:"uploadDate"`

	TotalLines     int               `json:"totalLines"`
	PrintTime      float32           `json:"printTime"`
	FilamentUsedMM float32           `json:"filamentUsedMm"`
	LayerIndexes   []gcodeLayerIndex `json:"layerIndexes"`
}

type gcodeLayerIndex struct {
	LineNumber  int `json:"lineNumber"`
	LayerNumber int `json:"layerNumber"`
}

func init() {
	gob.Register(gcodeLineIndex{})
	gob.Register(gcodeFileMeta{})
	gob.Register(gcodeLayerIndex{})
}

func (gfm *gcodeFileMeta) Save() error {
	f, err := os.OpenFile(filepath.Join(globalSettings.DataPath, "gcode_files/", gfm.GcodeFileName+".meta"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(gfm)
	if err != nil {
		return err
	}
	return nil
}

func loadGcodeFileMeta(path string) (*gcodeFileMeta, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	deco := gob.NewDecoder(f)
	var out *gcodeFileMeta
	err = deco.Decode(&out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (gfm *gcodeFileMeta) AnalyzeGcodeFile() error {
	f, err := os.Open(filepath.Join(globalSettings.DataPath, "gcode_files/", gfm.GcodeFileName))
	if err != nil {
		return err
	}
	startTime := time.Now()
	s := bufio.NewScanner(f)
	gfm.TotalLines = 0
	sim := &gcodeSimulator{
		layerIndexes: []gcodeLayerIndex{},
	}
	lineNumber := 0
	for s.Scan() {
		line := s.Text()
		gfm.TotalLines++
		lineNumber++
		sim.parseLine(line, lineNumber)
	}
	gfm.PrintTime = sim.time
	gfm.FilamentUsedMM = sim.filamentUsed
	gfm.LayerIndexes = sim.layerIndexes
	log.Printf("Finished gcode analysis in %v seconds", time.Now().Sub(startTime).Seconds())
	runtime.GC()
	return nil
}
