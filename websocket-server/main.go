package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
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

	// Send latest messages for each topic
	for topic := range topicReaders {
		wg.Add(1)
		go func(topic string, conn *websocket.Conn) {
			defer wg.Done()
			if err := sendLatestMessage(topic, conn); err != nil {
				log.Printf("Error sending latest message for topic %s: %v\n", topic, err)
			}
		}(topic, conn)
	}

	// Start goroutines to read Kafka messages for each reader
	for topic, reader := range topicReaders {
		wg.Add(1)
		go func(topic string, reader *kafka.Reader) {
			defer wg.Done()
			if err := readKafkaMessages(reader, conn, topic); err != nil {
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

// sendLatestMessage fetches the latest message from Kafka for a given topic and sends it to the WebSocket connection.
func sendLatestMessage(topic string, websocketConnection *websocket.Conn) error {
	latestReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{getKafkaURL()},
		Topic:       topic,
		MaxWait:     1 * time.Second,
		StartOffset: kafka.SeekCurrent,
	})

	defer latestReader.Close()

	msg, err := latestReader.ReadMessage(context.Background())
	if err != nil {
		log.Printf("Error reading latest message for topic %s: %v", topic, err)
		return err
	}

	log.Printf("Sending latest message for topic %s", topic)
	return writeToWebsocket(msg, websocketConnection)
}

// createKafkaReader creates a Kafka reader for the specified topic.
func createKafkaReader(topic string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{getKafkaURL()},
		Topic:   topic,
		GroupID: groupID,
	})
}

// readKafkaMessages reads messages from Kafka and sends them to the WebSocket connection.
func readKafkaMessages(reader *kafka.Reader, websocketConnection *websocket.Conn, topic string) error {
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

		log.Printf("Received message from Kafka topic %s", topic)
		if err := writeToWebsocket(msg, websocketConnection); err != nil {
			log.Printf("Error writing WebSocket message for topic %s: %v\n", topic, err)
			return err
		}
	}
}

// writeToWebsocket writes a Kafka message to the WebSocket connection.
func writeToWebsocket(msg kafka.Message, websocketConnection *websocket.Conn) error {
	wsMutex.Lock()
	defer wsMutex.Unlock()

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
		err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait))
		if err != nil {
			log.Printf("Ping error: %v\n", err)
			return
		}
	}
}

// getKafkaURL retrieves the Kafka URL from environment variables.
func getKafkaURL() string {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		return "localhost:9093"
	}
	return kafkaURL
}

// main starts the HTTP server and handles WebSocket connections.
func main() {
	http.HandleFunc("/ws", handleConnection)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("WebSocket server started on port %s\n", port)
	serverAddr := ":" + port
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
