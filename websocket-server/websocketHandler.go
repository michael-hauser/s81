package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

const (
	writeWait        = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait         = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod       = (pongWait * 9) / 10 // Send pings to peer with this period.
	closeGracePeriod = 10 * time.Second    // Time to wait before force close on connection.
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins
		},
	}
	wsMutex sync.Mutex
)

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while connecting to WebSocket: %v\n", err)
		return
	}
	defer func() {
		log.Println("WebSocket connection closed")
		conn.Close()
	}()

	log.Println("New WebSocket connection established")

	sendCachedMessages(conn)
	StreamKafkaMessages(conn)

	go ping(conn)

	wg.Wait()
}

func sendCachedMessages(conn *websocket.Conn) {
	topics := getTopics()

	for _, topic := range topics {
		sendLatestMessage(topic, conn)
	}
}

func sendLatestMessage(s string, conn *websocket.Conn) {
	cache.Lock()
	msg := cache.data[s]
	cache.Unlock()

	if err := writeToWebsocket(kafka.Message{
		Topic: s,
		Value: []byte(msg.Value),
	}, conn); err != nil {
		log.Printf("Error writing WebSocket message for topic %s: %v\n", s, err)
	}
}

func writeToWebsocket(msg kafka.Message, websocketConnection *websocket.Conn) error {
	wsMutex.Lock()
	defer wsMutex.Unlock()

	if msg.Value == nil || len(msg.Value) == 0 {
		return nil
	}

	webSocketValue := WebSocketValue{
		Key:   msg.Topic,
		Value: string(msg.Value),
	}

	jsonValue, err := json.Marshal(webSocketValue)
	if err != nil {
		log.Printf("Error marshaling WebSocketValue: %v\n", err)
		return err
	}

	return websocketConnection.WriteMessage(websocket.TextMessage, jsonValue)
}

func ping(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for range ticker.C {
		wsMutex.Lock()
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf("Error writing ping message: %v\n", err)
			wsMutex.Unlock()
			return
		}
		wsMutex.Unlock()
	}
}
