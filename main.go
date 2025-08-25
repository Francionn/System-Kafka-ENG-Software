package main

import (
	"context"
	"log"

	"kafka_system/consumer"
	"kafka_system/producer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// -------- PRODUCER --------
	prod, err := producer.NewProducer(ctx, "v1sensor", "my-topic-test", []string{"localhost:9094"})
	if err != nil {
		log.Fatalf("Error: create producer: %v", err)
	}

	// roda producer em paralelo
	go prod.Run("sensores")

	// -------- CONSUMERS --------
	brokers := []string{"localhost:9094", "localhost:9095", "localhost:9096"}
	consumer.RunConsumers(ctx, "consumer-group-1", "my-topic-test", brokers, 4)
}
