package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var (
	topicReaders = make(map[string]*kafka.Reader)
	wg           sync.WaitGroup
)

func main() {
	topics := getTopics()

	// Create Kafka readers for subway topics and weather
	for _, topic := range topics {
		uuid := uuid.New().String()
		groupID := "websocket-" + topic + "-group-" + uuid
		topicReaders[topic] = CreateKafkaReader(topic, groupID)
	}

	// Start goroutines to cache Kafka messages for each reader
	for topic, reader := range topicReaders {
		wg.Add(1)
		go func(topic string, reader *kafka.Reader) {
			defer wg.Done()
			CacheKafkaMessages(reader, topic)
		}(topic, reader)
	}

	// Start WebSocket server
	http.HandleFunc("/ws", HandleConnection)
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
	closeKafkaReaders()
}

func getTopics() []string {
	return []string{"subway-a", "subway-b", "subway-c", "weather-data"}
}

func closeKafkaReaders() {
	for _, reader := range topicReaders {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v\n", err)
		}
	}
}
