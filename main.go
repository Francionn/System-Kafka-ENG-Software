package main

import (
	"context"
	"log"

	"kafka_system/consumer"
	"kafka_system/infra"
)

func main() {
	ctx := context.Background()

	fsClient, err := infra.NewClient(ctx, "v1sensor")
	if err != nil {
		log.Fatalf("Erro Firestore: %v", err)
	}
	defer fsClient.Close()

	collection := fsClient.Collection("sensores")

	processMessage := infra.NewFirestoreHandler(ctx, collection)

	brokers := []string{"localhost:9094"}
	numWorkers := 4
	bufferSize := 1000
	topic := "my-topic-test-new"
	groupID := "consumer-group-test"

	consumer.RunConsumerWithHandler(ctx, groupID, topic, brokers, numWorkers, bufferSize, processMessage)
}
