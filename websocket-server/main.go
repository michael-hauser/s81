package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

const (
	writeWait        = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait         = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod       = (pongWait * 9) / 10 // Send pings to peer with this period.
	closeGracePeriod = 10 * time.Second    // Time to wait before force close on connection.
)

// Topics list
var topics = []string{"subway-a", "subway-b", "subway-c", "weather-data"}

// WebSocket upgrader for upgrading HTTP connections to WebSocket.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// WebSocketValue represents the message format sent over WebSocket.
type WebSocketValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Mutex for synchronizing writes to WebSocket connection.
var wsMutex sync.Mutex

// ConnectionManager manages active WebSocket connections and broadcasts messages.
type ConnectionManager struct {
	mu          sync.Mutex
	connections map[*websocket.Conn]bool
}

var manager = &ConnectionManager{
	connections: make(map[*websocket.Conn]bool),
}

// Main function to start the WebSocket server.
func main() {
	// Generate a unique identifier for this instance
	instanceID := uuid.New().String()

	// Start a Kafka consumer for each topic to broadcast messages
	for _, topic := range topics {
		go consumeAndBroadcast(topic, instanceID)
	}

	// Start WebSocket server
	http.HandleFunc("/ws", handleConnection)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("WebSocket server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

// handleConnection handles incoming WebSocket connections.
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while connecting to WebSocket: %v\n", err)
		return
	}
	defer func() {
		manager.removeConnection(conn)
		conn.Close()
	}()

	manager.addConnection(conn)

	log.Println("New WebSocket connection established")

	go ping(conn)
}

// consumeAndBroadcast reads messages from Kafka and broadcasts them to all active WebSocket connections.
func consumeAndBroadcast(topic string, instanceID string) {
	reader := createKafkaReader(topic, instanceID)
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			if err.Error() == "context canceled" {
				return
			}
			log.Printf("Error reading Kafka message for topic %s: %v\n", topic, err)
			continue
		}

		// Directly broadcast the message to all active connections
		manager.broadcastMessage(msg)
	}
}

// createKafkaReader creates a Kafka reader for the specified topic with a unique consumer group ID.
func createKafkaReader(topic string, instanceID string) *kafka.Reader {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9093"
	}
	groupID := "websocket-broadcast-" + topic + "-" + instanceID
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaURL},
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
	})
}

// broadcastMessage sends a Kafka message to all active WebSocket connections.
func (m *ConnectionManager) broadcastMessage(msg kafka.Message) {
	webSocketValue := WebSocketValue{
		Key:   msg.Topic,
		Value: string(msg.Value),
	}

	jsonValue, err := json.Marshal(webSocketValue)
	if err != nil {
		log.Printf("Error marshaling WebSocketValue: %v\n", err)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	for conn := range m.connections {
		wsMutex.Lock()
		err := conn.WriteMessage(websocket.TextMessage, jsonValue)
		wsMutex.Unlock()
		if err != nil {
			log.Printf("Error writing message to WebSocket: %v\n", err)
			m.removeConnection(conn)
			conn.Close()
		}
	}
}

// addConnection adds a new WebSocket connection to the manager and sends the latest message.
func (m *ConnectionManager) addConnection(conn *websocket.Conn) {
	m.mu.Lock()
	m.connections[conn] = true
	m.mu.Unlock()

	m.broadcastLatestMessages(conn)
}

// broadcastLatestMessages sends the latest message for each topic to the new connection.
func (m *ConnectionManager) broadcastLatestMessages(conn *websocket.Conn) {
	for _, topic := range topics {
		// Create a temporary Kafka reader to get the latest message for the topic
		reader := createKafkaReader(topic, uuid.New().String())
		defer reader.Close()

		// Seek to the last offset to get the latest message
		reader.SetOffset(kafka.LastOffset)

		// Read the latest message
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading latest Kafka message for topic %s: %v\n", topic, err)
			continue
		}

		// Send the latest message directly to the newly established WebSocket connection
		webSocketValue := WebSocketValue{
			Key:   msg.Topic,
			Value: string(msg.Value),
		}

		jsonValue, err := json.Marshal(webSocketValue)
		if err != nil {
			log.Printf("Error marshaling WebSocketValue: %v\n", err)
			return
		}

		wsMutex.Lock()
		errWs := conn.WriteMessage(websocket.TextMessage, jsonValue)
		wsMutex.Unlock()
		if errWs != nil {
			log.Printf("Error sending latest message to new WebSocket connection: %v\n", errWs)
			return
		}
	}
}

// removeConnection removes a WebSocket connection from the manager.
func (m *ConnectionManager) removeConnection(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.connections, conn)
}

// ping sends ping messages to keep the WebSocket connection alive.
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
