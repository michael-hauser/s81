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

// Cache to store the latest message for each topic.
var cache = struct {
	sync.RWMutex
	data map[string]WebSocketValue
}{data: make(map[string]WebSocketValue)}

// Main function to start the WebSocket server.
func main() {
	topics := []string{"subway-a", "subway-b", "subway-c", "weather-data"}
	topicReaders := make(map[string]*kafka.Reader)

	// Create Kafka readers for subway topics and weather
	for _, topic := range topics {
		uuid := uuid.New().String()
		groupID := "websocket-" + topic + "-group-" + uuid
		topicReaders[topic] = createKafkaReader(topic, groupID)
	}

	var wg sync.WaitGroup

	// Start goroutines to cache Kafka messages for each reader
	for topic, reader := range topicReaders {
		wg.Add(1)
		go func(topic string, reader *kafka.Reader) {
			defer wg.Done()
			cacheKafkaMessages(reader, topic)
		}(topic, reader)
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

	// Wait for all goroutines to finish
	wg.Wait()

	// Close all Kafka readers
	for _, reader := range topicReaders {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v\n", err)
		}
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
	topicReaders := make(map[string]*kafka.Reader)

	// Create Kafka readers for subway topics and weather
	for _, topic := range topics {
		groupID := "websocket-" + topic + "-group"
		topicReaders[topic] = createKafkaReader(topic, groupID)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Send latest messages for each topic from cache
	for topic := range topics {
		wg.Add(1)
		go func(topic string, conn *websocket.Conn) {
			defer wg.Done()
			sendLatestMessage(topic, conn)
		}(topics[topic], conn)
	}

	// Start goroutines to stream Kafka messages for each reader
	for topic, reader := range topicReaders {
		wg.Add(1)
		go func(topic string, reader *kafka.Reader) {
			defer wg.Done()
			if err := streamKafkaMessages(reader, conn, topic); err != nil {
				log.Printf("Error reading Kafka messages for topic %s: %v\n", topic, err)
			}
		}(topic, reader)
	}

	go ping(conn)

	// Wait for all goroutines to finish
	wg.Wait()

	// Close all Kafka readers
	mu.Lock()
	for _, reader := range topicReaders {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v\n", err)
		}
	}
	mu.Unlock()
}

// sendLatestMessage sends the latest message for the specified topic to the WebSocket connection.
func sendLatestMessage(s string, conn *websocket.Conn) {
	cache.Lock()
	msg := cache.data[s]
	cache.Unlock()

	// Send the message to the WebSocket connection.
	if err := writeToWebsocket(kafka.Message{
		Topic: s,
		Value: []byte(msg.Value),
	}, conn); err != nil {
		log.Printf("Error writing WebSocket message for topic %s: %v\n", s, err)
	}
}

// createKafkaReader creates a Kafka reader for the specified topic.
func createKafkaReader(topic string, groupID string) *kafka.Reader {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9093"
	}
	uuid := uuid.New().String()
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaURL},
		Topic:       topic,
		GroupID:     groupID + uuid,
		StartOffset: kafka.LastOffset,
	})
}

// streamKafkaMessages reads messages from Kafka and sends them to the WebSocket connection.
func streamKafkaMessages(reader *kafka.Reader, websocketConnection *websocket.Conn, topic string) error {
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			if err.Error() == "context canceled" {
				// Context was canceled, likely due to connection close
				return nil
			}
			log.Printf("Error reading Kafka message for topic %s: %v\n", topic, err)
			return err
		}

		if err := writeToWebsocket(msg, websocketConnection); err != nil {
			log.Printf("Error writing WebSocket message for topic %s: %v\n", topic, err)
			return err
		}
	}
}

// readKafkaMessages reads messages from Kafka, updates the cache, and sends them to WebSocket connections.
func cacheKafkaMessages(reader *kafka.Reader, topic string) {
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading Kafka message for topic %s: %v\n", topic, err)
			return
		}

		webSocketValue := WebSocketValue{
			Key:   msg.Topic,
			Value: string(msg.Value),
		}

		// Update the cache with the latest message
		cache.Lock()
		cache.data[topic] = webSocketValue
		cache.Unlock()
	}
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
