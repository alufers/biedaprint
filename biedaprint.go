package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

var activeConnections = []*websocket.Conn{}

var handlerMutex = &sync.Mutex{}

func handleWs(w http.ResponseWriter, r *http.Request) {
	log.Printf("New connection from %v (%v) \n", r.RemoteAddr, r.Header.Get("User-Agent"))

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	activeConnections = append(activeConnections, c)
	log.Printf("Active connections: %d", len(activeConnections))
	defer func() { // cleanup function
		handlerMutex.Lock()
		defer handlerMutex.Unlock()
		log.Println("Connection closed")
		// remove the connection when it drops
		activeConnections = removeConnFromArray(c, activeConnections)
		for _, tv := range trackedValues {
			tv.subscribers = removeConnFromArray(c, tv.subscribers)
		}
	}()
	defer c.Close()
	exit := false
	for {
		func() {
			handlerMutex.Lock()
			defer handlerMutex.Unlock()
			var msgData map[string]interface{}
			if err := c.ReadJSON(&msgData); err != nil {
				c.WriteJSON(map[string]interface{}{
					"msg":   "error",
					"error": "Bad json",
				})
				c.Close()
				exit = true
				return
			}

			handler, ok := messageHandlers[msgData["type"].(string)]
			if !ok {
				c.WriteJSON(jd{
					"type": "alert",
					"data": jd{
						"type":    "danger",
						"content": "Unknown action " + msgData["type"].(string),
					},
				})
				return
			}
			handler(c, msgData["data"])
		}()
		if exit {
			break
		}
	}
}

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
