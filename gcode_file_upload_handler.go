package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func handleGcodeFileUpload(w http.ResponseWriter, r *http.Request) {
	respondErr := func(err error) {
		w.WriteHeader(500)
		respondJSON(w, jd{
			"error": err.Error(),
		})
		log.Printf("Gcode upload error: %v", err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseMultipartForm(10 * 1024 * 1024)
	file, hdr, err := r.FormFile("file")
	if err != nil {
		respondErr(err)
		return
	}
	defer file.Close()

	gcodeFilename := RandStringRunes(8) + ".gcode"
	// copy example
	f, err := os.OpenFile(filepath.Join(globalSettings.DataPath, "gcode_files/"+gcodeFilename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		respondErr(err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		respondErr(err)
		return
	}
	meta := &gcodeFileMeta{
		OriginalName:  hdr.Filename,
		UploadDate:    time.Now(),
		GcodeFileName: gcodeFilename,
	}
	err = meta.AnalyzeGcodeFile()
	if err != nil {
		respondErr(err)
		return
	}
	err = meta.Save()
	if err != nil {
		respondErr(err)
		return
	}
	respondJSON(w, meta)
}
