package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"go-torrent-web-client/torrent"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketManager struct {
	clients   map[*websocket.Conn]bool
	broadcast chan map[string]float64
	tm        *torrent.TorrentManager
}

func NewWebSocketManager(tm *torrent.TorrentManager) *WebSocketManager {
	return &WebSocketManager{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan map[string]float64),
		tm:        tm,
	}
}

func (wsm *WebSocketManager) SendProgresses(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("New WebSocket client connected")

	for {
		progresses := wsm.tm.GetProgresses()
		progressData, err := json.Marshal(progresses)
		if err != nil {
			fmt.Println("Error marshaling progress data:", err)
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, progressData)
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}

		time.Sleep(1 * time.Second)
	}
}
