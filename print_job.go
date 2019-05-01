package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type printJob struct {
	gcodeMeta             *gcodeFileMeta
	lineResendBuffer      map[int]string
	lineResendBufferMutex *sync.RWMutex
	currentLine           int
	currentNonBlankLine   int
	gcodeFile             *os.File
	scanner               *bufio.Scanner
}

func (pj *printJob) computeLineChecksum(line string) int {
	var cs int
	for i := 0; i < len(line) && line[i] != '*'; i++ {
		cs = cs ^ int(line[i])
	}
	cs &= 0xff
	return cs
}

//jobLines returns a channel which sends lines together with newline chars and checksums
func (pj *printJob) jobLines() (chan string, error) {
	pj.lineResendBufferMutex = &sync.RWMutex{}
	pj.lineResendBuffer = make(map[int]string)
	var err error
	pj.gcodeFile, err = os.Open(filepath.Join(globalSettings.DataPath, "gcode_files/", pj.gcodeMeta.GcodeFileName))
	if err != nil {
		return nil, err
	}

	pj.scanner = bufio.NewScanner(pj.gcodeFile)
	c := make(chan string)
	go func() {
		defer close(c)
		defer pj.gcodeFile.Close()
		log.Printf("Starting jobLines goroutine...")
		c <- "M110 N0\r\n"
		for pj.scanner.Scan() {
			rawLine := strings.Split(pj.scanner.Text(), ";")[0]

			lineWithNumber := fmt.Sprintf("N%d %v", pj.currentNonBlankLine+1, rawLine)
			lineWithChecksum := fmt.Sprintf("%v*%v\r\n", lineWithNumber, pj.computeLineChecksum(lineWithNumber))
			if strings.TrimSpace(rawLine) != "" {
				log.Printf("Sending gcode line %v of %v", pj.currentLine+1, pj.gcodeMeta.TotalLines)
				c <- lineWithChecksum
				pj.lineResendBufferMutex.Lock()
				pj.lineResendBuffer[pj.currentNonBlankLine] = lineWithChecksum
				pj.currentNonBlankLine++
				delete(pj.lineResendBuffer, pj.currentNonBlankLine-10)
				pj.lineResendBufferMutex.Unlock()
				trackedValues["printProgress"].updateValue((float64(pj.currentLine) / float64(pj.gcodeMeta.TotalLines)) * 100)
			}
			pj.currentLine++

		}
	}()

	return c, nil
}

func (pj *printJob) getLineForResend(number int) string {
	pj.lineResendBufferMutex.RLock()
	defer pj.lineResendBufferMutex.RUnlock()
	return pj.lineResendBuffer[number]
}
