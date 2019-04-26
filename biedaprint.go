package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr"
)

func main() {
	log.Printf("Starting biedaprint...\n")
	loadSettings()
	err := os.MkdirAll(filepath.Join(globalSettings.DataPath, "gcode_files"), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create gcode_files data directory: %v", err)
		return
	}
	go serialReader()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWs)
	mux.HandleFunc("/gcode-file-upload", handleGcodeFileUpload)
	box := packr.NewBox("./static")
	mux.Handle("/", interceptHandler(http.FileServer(box), func(w http.ResponseWriter, status int) {
		data, _ := box.FindString("index.html")
		w.Header().Set("Content-type", "text/html")
		fmt.Fprint(w, data)
	}))

	log.Printf("%v\n", http.ListenAndServe("0.0.0.0:4444", mux))
}
