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

// Main function to start the WebSocket server.
func main() {
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

// handleConnection handles incoming WebSocket connections and Kafka messages.
func handleConnection(w http.ResponseWriter, r *http.Request) {
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

	topics := []string{"subway-a", "subway-b", "subway-c", "weather-data"}
	var wg sync.WaitGroup

	// Start a Kafka consumer for each topic
	for _, topic := range topics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			consumeLatestAndStream(conn, topic)
		}(topic)
	}

	go ping(conn)

	// Wait for all goroutines to finish
	wg.Wait()
}

// consumeLatestAndStream reads the latest message from Kafka and then continues to stream new messages.
func consumeLatestAndStream(conn *websocket.Conn, topic string) {
	reader := createKafkaReader(topic)
	defer reader.Close()

	// Consume the latest message first
	consumeLatestMessage(reader, conn, topic)

	// Continue to stream new messages
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			if err.Error() == "context canceled" {
				// Context was canceled, likely due to connection close
				return
			}
			log.Printf("Error reading Kafka message for topic %s: %v\n", topic, err)
			return
		}

		if err := writeToWebsocket(msg, conn); err != nil {
			log.Printf("Error writing WebSocket message for topic %s: %v\n", topic, err)
			return
		}
	}
}

// consumeLatestMessage consumes the latest message from Kafka.
func consumeLatestMessage(reader *kafka.Reader, conn *websocket.Conn, topic string) {
	// Seek to the last offset to get the latest message
	reader.SetOffset(kafka.LastOffset)

	// Read the latest message
	msg, err := reader.ReadMessage(context.Background())
	if err != nil {
		log.Printf("Error reading latest Kafka message for topic %s: %v\n", topic, err)
		return
	}

	// Send the latest message to the WebSocket connection
	if err := writeToWebsocket(msg, conn); err != nil {
		log.Printf("Error writing latest WebSocket message for topic %s: %v\n", topic, err)
	}
}

// createKafkaReader creates a Kafka reader for the specified topic.
func createKafkaReader(topic string) *kafka.Reader {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9093"
	}
	groupID := "websocket-" + topic + "-group-" + uuid.New().String()
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaURL},
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
	})
}

// writeToWebsocket writes a Kafka message to the WebSocket connection.
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
