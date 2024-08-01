package main

import (
	"context"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type WebSocketValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var cache = struct {
	sync.RWMutex
	data map[string]WebSocketValue
}{data: make(map[string]WebSocketValue)}

func CacheKafkaMessages(reader *kafka.Reader, topic string) {
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

		cache.Lock()
		cache.data[topic] = webSocketValue
		cache.Unlock()
	}
}
