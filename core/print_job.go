package core

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

//PrintJobInternal wraps PrintJob with some fields used internally by the various methods
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
	pj.gcodeFile, err = os.Open(filepath.Join(pj.app.GetDataPath(), "gcode_files/", pj.GcodeMeta.GcodeFileName))
	if err != nil {
		return nil, err
	}

	pj.scanner = bufio.NewScanner(pj.gcodeFile)
	c := make(chan string)
	go func() {
		defer close(c)
		defer pj.gcodeFile.Close()
		log.Printf("Starting jobLines goroutine...")
		pj.app.TrackedValuesService.TrackedValues["printOriginalName"].UpdateValue(pj.GcodeMeta.OriginalName)
		pj.app.TrackedValuesService.TrackedValues["isPrinting"].UpdateValue(true)
		pj.app.TrackedValuesService.TrackedValues["printStartTime"].UpdateValue(pj.StartedTime.Format(time.RFC3339))
		pj.app.TrackedValuesService.TrackedValues["printCurrentLayer"].UpdateValue(0)
		pj.app.TrackedValuesService.TrackedValues["printTotalLayers"].UpdateValue(len(pj.GcodeMeta.LayerIndexes))
		defer pj.app.TrackedValuesService.TrackedValues["isPrinting"].UpdateValue(false)
		c <- "M110 N0\r\n"
		for pj.scanner.Scan() {
			rawLine := strings.Split(pj.scanner.Text(), ";")[0]

			lineWithNumber := fmt.Sprintf("N%d %v", pj.currentNonBlankLine+1, rawLine)
			lineWithChecksum := fmt.Sprintf("%v*%v\r\n", lineWithNumber, pj.computeLineChecksum(lineWithNumber)) // add the checksum and newlines after a star character
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
				delete(pj.lineResendBuffer, pj.currentNonBlankLine-10) // store last 10 lines for resending
				pj.lineResendBufferMutex.Unlock()
				pj.app.TrackedValuesService.TrackedValues["printProgress"].UpdateValue((float64(pj.currentLine) / float64(pj.GcodeMeta.TotalLines)) * 100)
			}
			pj.currentLine++

			if pj.currentLayerIndex < len(pj.GcodeMeta.LayerIndexes)-1 {
				if pj.currentLine >= pj.GcodeMeta.LayerIndexes[pj.currentLayerIndex].LineNumber {
					pj.currentLayerIndex++
					pj.app.TrackedValuesService.TrackedValues["printCurrentLayer"].UpdateValue(pj.GcodeMeta.LayerIndexes[pj.currentLayerIndex].LayerNumber)
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
