package biedaprint

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type PrintJobInternal struct {
	*PrintJob
	app                   *App
	lineResendBuffer      map[int]string
	lineResendBufferMutex *sync.RWMutex
	currentLine           int
	currentNonBlankLine   int
	currentLayerIndex     int
	gcodeFile             *os.File
	scanner               *bufio.Scanner
	abortSem              chan bool
}

func (pj *PrintJobInternal) computeLineChecksum(line string) int {
	var cs int
	for i := 0; i < len(line) && line[i] != '*'; i++ {
		cs = cs ^ int(line[i])
	}
	cs &= 0xff
	return cs
}

//jobLines returns a channel which sends lines together with newline chars and checksums
func (pj *PrintJobInternal) jobLines() (chan string, error) {
	pj.lineResendBufferMutex = &sync.RWMutex{}
	pj.abortSem = make(chan bool)
	pj.lineResendBuffer = make(map[int]string)
	var err error
	pj.gcodeFile, err = os.Open(filepath.Join(pj.app.GetSettings().DataPath, "gcode_files/", pj.GcodeMeta.GcodeFileName))
	if err != nil {
		return nil, err
	}

	pj.scanner = bufio.NewScanner(pj.gcodeFile)
	c := make(chan string)
	go func() {
		defer close(c)
		defer pj.gcodeFile.Close()
		log.Printf("Starting jobLines goroutine...")
		pj.app.TrackedValuesManager.TrackedValues["printOriginalName"].UpdateValue(pj.GcodeMeta.OriginalName)
		pj.app.TrackedValuesManager.TrackedValues["isPrinting"].UpdateValue(true)
		pj.app.TrackedValuesManager.TrackedValues["printStartTime"].UpdateValue(pj.StartedTime.Format(time.RFC3339))
		pj.app.TrackedValuesManager.TrackedValues["printCurrentLayer"].UpdateValue(0)
		pj.app.TrackedValuesManager.TrackedValues["printTotalLayers"].UpdateValue(len(pj.GcodeMeta.LayerIndexes))
		defer pj.app.TrackedValuesManager.TrackedValues["isPrinting"].UpdateValue(false)
		c <- "M110 N0\r\n"
		for pj.scanner.Scan() {
			rawLine := strings.Split(pj.scanner.Text(), ";")[0]

			lineWithNumber := fmt.Sprintf("N%d %v", pj.currentNonBlankLine+1, rawLine)
			lineWithChecksum := fmt.Sprintf("%v*%v\r\n", lineWithNumber, pj.computeLineChecksum(lineWithNumber))
			if strings.TrimSpace(rawLine) != "" {
				log.Printf("Sending gcode line %v of %v", pj.currentLine+1, pj.GcodeMeta.TotalLines)
				select {
				case c <- lineWithChecksum:
				case <-pj.abortSem:
					return
				}

				pj.lineResendBufferMutex.Lock()
				pj.lineResendBuffer[pj.currentNonBlankLine] = lineWithChecksum
				pj.currentNonBlankLine++
				delete(pj.lineResendBuffer, pj.currentNonBlankLine-10)
				pj.lineResendBufferMutex.Unlock()
				pj.app.TrackedValuesManager.TrackedValues["printProgress"].UpdateValue((float64(pj.currentLine) / float64(pj.GcodeMeta.TotalLines)) * 100)
			}
			pj.currentLine++

			if pj.currentLayerIndex < len(pj.GcodeMeta.LayerIndexes) {
				if pj.currentLine >= pj.GcodeMeta.LayerIndexes[pj.currentLayerIndex].LineNumber {
					pj.currentLayerIndex++
					pj.app.TrackedValuesManager.TrackedValues["printCurrentLayer"].UpdateValue(pj.GcodeMeta.LayerIndexes[pj.currentLayerIndex].LayerNumber)
				}
			}

		}
	}()

	return c, nil
}

func (pj *PrintJobInternal) getLineForResend(number int) string {
	pj.lineResendBufferMutex.RLock()
	defer pj.lineResendBufferMutex.RUnlock()
	return pj.lineResendBuffer[number]
}

func (pj *PrintJobInternal) abort() {
	pj.abortSem <- true
}
