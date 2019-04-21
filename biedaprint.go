package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/packr"
)

func main() {
	log.Printf("Starting biedaprint...\n")
	loadSettings()
	go serialReader()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWs)
	box := packr.NewBox("./static")
	mux.Handle("/", interceptHandler(http.FileServer(box), func(w http.ResponseWriter, status int) {
		data, _ := box.FindString("index.html")
		w.Header().Set("Content-type", "text/html")
		fmt.Fprint(w, data)
	}))

	log.Printf("%v\n", http.ListenAndServe("0.0.0.0:4444", mux))
}
