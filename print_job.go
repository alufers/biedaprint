package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type printJob struct {
	gcodeMeta             *gcodeFileMeta
	lineResendBuffer      map[int]string
	lineResendBufferMutex *sync.RWMutex
	currentLine           int
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
		for pj.scanner.Scan() {
			rawLine := pj.scanner.Text()
			lineWithNumber := fmt.Sprintf("N%d %v ", pj.currentLine+1, rawLine)
			lineWithChecksum := fmt.Sprintf("%v*%v\r\n", lineWithNumber, pj.computeLineChecksum(lineWithNumber))
			log.Printf("Sending gcode line %v of %v", pj.currentLine+1, pj.gcodeMeta.TotalLines)
			c <- lineWithChecksum
			pj.lineResendBufferMutex.Lock()
			pj.lineResendBuffer[pj.currentLine] = lineWithChecksum
			pj.currentLine++
			delete(pj.lineResendBuffer, pj.currentLine-10)
			pj.lineResendBufferMutex.Unlock()
		}
	}()

	return c, nil
}

func (pj *printJob) getLineForResend(number int) string {
	pj.lineResendBufferMutex.RLock()
	defer pj.lineResendBufferMutex.RUnlock()
	return pj.lineResendBuffer[number]
}
