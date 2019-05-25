package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/phayes/freeport"
)

func main() {
	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}

	serverConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: port, Zone: ""})
	if err != nil {
		log.Panic("Failed to ListenUDP for discovery: %v", err)
		return
	}
	defer serverConn.Close()
	go func() {
		for {
			func() {
				addr, err := net.ResolveUDPAddr("udp", "192.168.1.255:4444")
				if err != nil {
					panic(err)
				}
				conn, err := net.DialUDP("udp", nil, addr)
				if err != nil {
					log.Printf("Failed to send discovery packet %v", err)
					return
				}
				defer conn.Close()
				fmt.Fprintf(conn, "BIEDAPRINT_DISCOVERY:%d\n", port)
			}()
			time.Sleep(5 * time.Second)
		}
	}()

	buf := make([]byte, 1024)
	for {

		n, addr, _ := serverConn.ReadFromUDP(buf)
		var jsonData interface{}
		err = json.Unmarshal(buf[0:n], &jsonData)
		if err != nil {
			log.Panic("Bad json response %v", err)
			return
		}
		if m, ok := jsonData.(map[string]interface{}); ok {
			if val, _ := m["biedaprint"].(bool); val {
				fmt.Printf("Found biedaprint instance at http://%v:4444\n", addr.IP.String())
			}
		}

	}

}
