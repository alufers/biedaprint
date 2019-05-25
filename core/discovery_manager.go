package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type DiscoveryManager struct {
	app *App
}

func NewDiscoveryManager(app *App) *DiscoveryManager {
	return &DiscoveryManager{
		app: app,
	}
}

func (dm *DiscoveryManager) Init() {
	log.Printf("Starting UDP DiscoveryManager!")
	go func() {
		serverConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 4444, Zone: ""})
		if err != nil {
			log.Panic("Failed to ListenUDP for DiscoveryManager: %v", err)
			return
		}
		defer serverConn.Close()
		buf := make([]byte, 1024)
		for {
			n, addr, _ := serverConn.ReadFromUDP(buf)
			strData := string(buf[0:n])
			if strings.HasPrefix(strData, "BIEDAPRINT_DISCOVERY:") {
				func() {

					log.Printf("Recieved discovery packet %v from %v", strData, addr)
					var replyPort int
					fmt.Sscanf(strData, "BIEDAPRINT_DISCOVERY:%d", &replyPort)
					conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
						IP:   addr.IP,
						Port: replyPort,
						Zone: addr.Zone,
					})
					if err != nil {
						log.Printf("Failed to send discovery reply packet %v", err)
						return
					}
					defer conn.Close()
					jsonData, _ := json.Marshal(map[string]interface{}{
						"biedaprint": true,
					})
					conn.Write(append(jsonData, byte('\n')))
				}()
			}

		}
	}()
}
