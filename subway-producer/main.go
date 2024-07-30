package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	gtfs_realtime "github.com/michael-hauser/s81/subway-producer/gtfs-realtime"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

// SubwayConfig holds the configuration for each train line
type SubwayConfig struct {
	Name        string
	Endpoint    string
	TripRouteID string
	Stops       []string
	Topic       string
}

// Global configuration for all train lines
var trainConfigs = map[string]SubwayConfig{
	"A": {
		Name:        "A",
		Endpoint:    "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace",
		TripRouteID: "A",
		Stops:       []string{"A21N", "A21S"},
		Topic:       "subway-a",
	},
	"B": {
		Name:        "B",
		Endpoint:    "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-bdfm",
		TripRouteID: "D",
		Stops:       []string{"B21N", "B21S"},
		Topic:       "subway-b",
	},
	"C": {
		Name:        "C",
		Endpoint:    "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace",
		TripRouteID: "C",
		Stops:       []string{"A21N", "A21S"},
		Topic:       "subway-c",
	},
}

// Main function
func main() {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9093"
	}

	// Create Kafka writers for each train line
	writers := make(map[string]*kafka.Writer)
	for _, config := range trainConfigs {
		writers[config.Name] = kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{kafkaURL},
			Topic:   config.Topic,
		})
		defer writers[config.Name].Close()
	}

	// Set the interval for fetching data
	interval := 30 * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			fetchAndPublishSubwayData(writers)
			log.Println("Fetched and published subway data")
		}
	}()

	log.Println("Subway data producer started")

	select {}
}

// fetchAndPublishSubwayData fetches the subway data for each line and publishes it to Kafka
func fetchAndPublishSubwayData(writers map[string]*kafka.Writer) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, config := range trainConfigs {
		res, err := client.Get(config.Endpoint)
		if err != nil {
			log.Printf("Error fetching data for %s: %v", config.Name, err)
			continue
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading response body for %s: %v", config.Name, err)
			res.Body.Close()
			continue
		}
		res.Body.Close()

		feedMessage := &gtfs_realtime.FeedMessage{}
		err = proto.Unmarshal(body, feedMessage)
		if err != nil {
			log.Printf("Error unmarshalling feed for %s: %v", config.Name, err)
			continue
		}

		filteredFeed := filterFeedForLine(feedMessage, config)
		if err := publishToKafka(writers[config.Name], config.Name, filteredFeed); err != nil {
			log.Printf("Error writing %s message to Kafka: %v", config.Name, err)
		}
	}
}

// filterFeedForLine filters the feed message for a specific train line
func filterFeedForLine(feedMessage *gtfs_realtime.FeedMessage, config SubwayConfig) *gtfs_realtime.FeedMessage {
	var filteredEntities []*gtfs_realtime.FeedEntity

	for _, entity := range feedMessage.Entity {
		if entity.Vehicle != nil && entity.Vehicle.StopId != nil {
			stopId := *entity.Vehicle.StopId
			if contains(config.Stops, stopId) {
				filteredEntities = append(filteredEntities, entity)
			}
		}
		if entity.TripUpdate != nil {
			// Filter stop_time_update for relevant stops
			filteredStopTimeUpdates := filterStopTimeUpdates(entity.TripUpdate.StopTimeUpdate, config.Stops)
			if len(filteredStopTimeUpdates) > 0 {
				if entity.TripUpdate.Trip != nil {
					routeId := *entity.TripUpdate.Trip.RouteId
					if routeId == config.TripRouteID {
						entity.TripUpdate.StopTimeUpdate = filteredStopTimeUpdates
						filteredEntities = append(filteredEntities, entity)
					}
				}
			}
		}
	}

	// Log the number of entities filtered for the line
	log.Printf("Filtered %d entities for %s", len(filteredEntities), config.Name)

	return &gtfs_realtime.FeedMessage{Entity: filteredEntities}
}

// filterStopTimeUpdates filters the stop time updates for relevant stops
func filterStopTimeUpdates(updates []*gtfs_realtime.TripUpdate_StopTimeUpdate, relevantStops []string) []*gtfs_realtime.TripUpdate_StopTimeUpdate {
	var filtered []*gtfs_realtime.TripUpdate_StopTimeUpdate
	for _, update := range updates {
		if update.StopId != nil {
			stopId := *update.StopId
			if contains(relevantStops, stopId) {
				filtered = append(filtered, update)
			}
		}
	}
	return filtered
}

// contains checks if a slice contains a specific item
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// publishToKafka publishes the feed message to Kafka
func publishToKafka(writer *kafka.Writer, key string, feedMessage *gtfs_realtime.FeedMessage) error {
	feedMessageJSON, err := json.Marshal(feedMessage)
	if err != nil {
		return err
	}

	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: feedMessageJSON,
		},
	)
}
