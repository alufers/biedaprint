package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func byteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func removeConnFromArray(c *websocket.Conn, array []*websocket.Conn) []*websocket.Conn {
	newConnections := make([]*websocket.Conn, 0, len(array))
	for _, a := range array {
		if a != c {
			newConnections = append(newConnections, a)
		}
	}
	return newConnections
}

func respondJSON(w io.Writer, dat interface{}) {
	enc := json.NewEncoder(w)
	enc.Encode(dat)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
