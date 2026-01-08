package main

import (
	"context"
	"log"

	"kafka_system/consumer"
	"kafka_system/infra"
)

func main() {
	ctx := context.Background()

	//  Kafka configs
	brokers := []string{"localhost:9094"}
	numWorkers := 4
	bufferSize := 500
	topic := "my-topic-test-new"
	groupID := "consumer-group-test"

	// Firestore configs
	projectID := "v1sensor"
	firestoreCollection := "sensoresv3"

	firestoreClient, err := infra.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
	defer firestoreClient.Close()

	processMessage := consumer.DefaultProcessMessage

	log.Println("Starting consumer Kafka: API and Firestore")

	consumer.RunConsumerWithHandler(
		ctx,
		groupID,
		topic,
		brokers,
		numWorkers,
		bufferSize,
		processMessage,
		firestoreClient,
		firestoreCollection,
	)
}
