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
	writeWait        = 15 * time.Second    // Increased time allowed to write a message to the peer.
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

// ConnectionManager manages active WebSocket connections and broadcasts messages.
type ConnectionManager struct {
	mu             sync.RWMutex
	connections    map[*websocket.Conn]struct{}
	latestMessages map[string]kafka.Message
}

var manager = &ConnectionManager{
	connections:    make(map[*websocket.Conn]struct{}),
	latestMessages: make(map[string]kafka.Message),
}

// Main function to start the WebSocket server.
func main() {
	// Generate a unique identifier for this instance
	instanceID := uuid.New().String()

	// Start Kafka consumers
	for _, topic := range topics {
		go consumeAndSendDirectly(topic, instanceID)
	}

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

	manager.addConnection(conn)
	log.Println("New WebSocket connection established")

	go ping(conn)
}

// consumeAndSendDirectly reads messages from Kafka and immediately sends them to all active WebSocket connections.
func consumeAndSendDirectly(topic string, instanceID string) {
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
		// Update the latest message
		manager.updateLatestMessage(topic, msg)
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

	m.mu.RLock()
	defer m.mu.RUnlock()

	for conn := range m.connections {
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		err := conn.WriteMessage(websocket.TextMessage, jsonValue)
		if err != nil {
			log.Printf("Error writing message to WebSocket: %v\n", err)
			m.removeAndCloseConnection(conn)
		}
	}
}

// addConnection adds a new WebSocket connection to the manager and sends the latest messages.
func (m *ConnectionManager) addConnection(conn *websocket.Conn) {
	m.mu.Lock()
	m.connections[conn] = struct{}{}
	m.mu.Unlock()

	m.sendLatestMessages(conn)
}

// sendLatestMessages sends the latest message for each topic to the new connection.
func (m *ConnectionManager) sendLatestMessages(conn *websocket.Conn) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, msg := range m.latestMessages {
		webSocketValue := WebSocketValue{
			Key:   msg.Topic,
			Value: string(msg.Value),
		}

		jsonValue, err := json.Marshal(webSocketValue)
		if err != nil {
			log.Printf("Error marshaling WebSocketValue: %v\n", err)
			continue
		}

		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.TextMessage, jsonValue); err != nil {
			log.Printf("Error sending latest message to new WebSocket connection: %v\n", err)
			m.removeAndCloseConnection(conn)
			return
		}
	}
}

// updateLatestMessage updates the latest message for a given topic.
func (m *ConnectionManager) updateLatestMessage(topic string, msg kafka.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.latestMessages[topic] = msg
}

// removeAndCloseConnection removes a WebSocket connection from the manager and ensures it's properly closed.
func (m *ConnectionManager) removeAndCloseConnection(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the connection is still in the map before attempting to remove and close it
	if _, ok := m.connections[conn]; ok {
		log.Println("Removing and closing connection due to error")
		conn.Close()                // Close the WebSocket connection
		delete(m.connections, conn) // Remove the connection from the map
	}
}

// ping sends ping messages to keep the WebSocket connection alive.
func ping(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for range ticker.C {
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf("Error writing ping message: %v\n", err)
			manager.removeAndCloseConnection(conn) // Properly remove the connection if ping fails
			return
		}
	}
}
