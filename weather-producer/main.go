package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

// Mutex for synchronizing writes to WebSocket connection.
var wsMutex sync.Mutex

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	kafkaURL := os.Getenv("KAFKA_URL")
	apiKey := os.Getenv("WEATHER_API_KEY")
	if kafkaURL == "" {
		kafkaURL = "localhost:9093"
	}
	topic := "weather-data"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaURL},
		Topic:   topic,
	})
	defer writer.Close()

	log.Println("Weather producer started")

	// Run the function immediately
	go fetchAndPublishWeatherData(writer, apiKey)

	// Set the interval for fetching data to every 10 minutes
	interval := 10 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Set up a goroutine to run the function periodically
	go func() {
		for range ticker.C {
			fetchAndPublishWeatherData(writer, apiKey)
		}
	}()

	// Block main thread
	select {}
}

func fetchAndPublishWeatherData(writer *kafka.Writer, apiKey string) {
	wsMutex.Lock()
	defer wsMutex.Unlock()

	lat := "40.781433"
	long := "-73.972143"
	weatherEndpoint := "https://api.openweathermap.org/data/3.0/onecall?units=imperial&lat=" + lat + "&lon=" + long + "&appid=" + apiKey

	log.Println("Fetching weather data from:", weatherEndpoint)

	res, err := http.Get(weatherEndpoint)
	if err != nil {
		log.Println("Error fetching weather data:", err)
		return
	}
	defer res.Body.Close()

	log.Println("Weather data fetched successfully")

	body, _ := io.ReadAll(res.Body)
	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("weather"),
			Value: body,
		},
	)

	log.Println("Weather data published to Kafka")

	if err != nil {
		log.Println("Error writing message to Kafka:", err)
	}
}
