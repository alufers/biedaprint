package main

import (
	"log"
	"net/http"
	"sync"

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
		serialConsoleSubscribers = removeConnFromArray(c, serialConsoleSubscribers)
	}()
	defer c.Close()
	exit := false
	for {
		func() {

			var msgData map[string]interface{}
			if err := c.ReadJSON(&msgData); err != nil {
				log.Printf("Bad json from conn: %v", err)
				c.Close()
				exit = true
				return
			}

			handlerMutex.Lock()
			defer handlerMutex.Unlock()

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
