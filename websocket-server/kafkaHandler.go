package main

import (
	"context"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

func CreateKafkaReader(topic string, groupID string) *kafka.Reader {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9093"
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaURL},
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
	})
}

func StreamKafkaMessages(conn *websocket.Conn) {
	for topic, reader := range topicReaders {
		wg.Add(1)
		go func(topic string, reader *kafka.Reader) {
			defer wg.Done()
			if err := streamKafkaTopicMessages(reader, conn, topic); err != nil {
				log.Printf("Error reading Kafka messages for topic %s: %v\n", topic, err)
			}
		}(topic, reader)
	}
}

func streamKafkaTopicMessages(reader *kafka.Reader, websocketConnection *websocket.Conn, topic string) error {
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			if err.Error() == "context canceled" {
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
