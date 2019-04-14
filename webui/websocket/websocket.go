package websocket

import (
	"net/http"
	"time"
	"encoding/json"

	"github.com/Pegasus8/piworker/processment/stats"
	"github.com/Pegasus8/piworker/utilities/log"
	"github.com/gorilla/websocket"
)

// Read and write buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

// Upgrade func takes incoming connections and upgrade the request into a WebSocket connection
func Upgrade(w http.ResponseWriter, request *http.Request) (*websocket.Conn, error) {

	// Allow other origins connection
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// WebSocket connection
	ws, err := upgrader.Upgrade(w, request, nil)
	if err != nil {
		log.Criticalln(err)
		return ws, err
	}
	log.Infoln("Connected with:", request.RemoteAddr)
	// Return WebSocket connection
	return ws, nil
}

// Writer func sends data into WebSocket to the client
func Writer(conn *websocket.Conn) {

	
	// Other way to do that
	// // for {
	// // 	ticker := time.NewTicker(5 * time.Second)
	// // 	for t := range ticker.C { ... }
	// // }
	
	log.Infoln("Sending data to ", conn.RemoteAddr())
	// Send data to client every 1 sec
	for range time.Tick(1 * time.Second) {
		// Get data
		data, err := stats.GetStatistics()
		if err != nil {
			log.Criticalln(err)
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Criticalln(err)
			return
		}
	
		// Send data
		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Criticalln(err)
			return
		}
	}
}
	
