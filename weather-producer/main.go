package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9093"
	}
	topic := "weather-data"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaURL},
		Topic:   topic,
	})
	defer writer.Close()

	// Run the function immediately
	fetchAndPublishWeatherData(writer)

	// Set the interval for fetching data to every 10 minutes
	interval := 10 * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Set up a goroutine to run the function periodically
	go func() {
		for range ticker.C {
			fetchAndPublishWeatherData(writer)
			log.Println("Weather data fetched and published")
		}
	}()

	log.Println("Weather producer started")

	// Block main thread
	select {}
}

func fetchAndPublishWeatherData(writer *kafka.Writer) {
	apiKey := os.Getenv("API_KEY")
	lat := os.Getenv("LAT")
	long := os.Getenv("LONG")
	weatherEndpoint := "https://api.openweathermap.org/data/3.0/onecall?units=imperial&lat=" + lat + "&lon=" + long + "&appid=" + apiKey

	res, err := http.Get(weatherEndpoint)
	if err != nil {
		log.Println("Error fetching weather data:", err)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("weather"),
			Value: body,
		},
	)
	if err != nil {
		log.Println("Error writing message to Kafka:", err)
	}
}
