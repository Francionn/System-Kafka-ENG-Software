package consumer

import (
	"context"
	"encoding/json"
	"log"
	"sync/atomic"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/IBM/sarama"
)

func RunConsumerWithHandler(
	ctx context.Context,
	groupID string,
	topic string,
	brokers []string,
	numWorkers int,
	bufferSize int,
	processMessage func(msg []byte),
	firestoreClient *firestore.Client,
	firestoreCollection string,
) {
	// Kafka config
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Infra
	messages := make(chan *sarama.ConsumerMessage, bufferSize)
	var counter uint64
	collection := firestoreClient.Collection(firestoreCollection)

	// Worker pool
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			for msg := range messages {
				atomic.AddUint64(&counter, 1)
				current := atomic.LoadUint64(&counter)

				processMessage(msg.Value)
	
				func(data []byte) {
					var sensorData SensorData

					if err := json.Unmarshal(data, &sensorData); err != nil {
						log.Printf("[Worker-%d] JSON unmarshal error: %v", workerID, err)
						return
					}

					sensorData.Status = analyzeSensorStatus(&sensorData)
					sensorData.QIA = CalculateQIA(&sensorData)

					ts := time.Time(sensorData.Timestamp)

					doc := map[string]interface{}{
						"sensorID":      sensorData.SensorID,
						"name":          sensorData.Name,
						"latitude":      sensorData.Latitude,
						"longitude":     sensorData.Longitude,
						"temperature":   sensorData.Temperature,
						"humidity":      sensorData.Humidity,
						"co2":           sensorData.CO2,
						"pm25":          sensorData.PM25,
						"pm10":          sensorData.PM10,
						"windDirection": sensorData.WindDirection,
						"status":        sensorData.Status,
						"qia":           sensorData.QIA,

						"timestamp": ts,

						"year":  ts.Year(),
						"month": int(ts.Month()),
						"day":   ts.Day(),

						"topic":     msg.Topic,
						"partition": msg.Partition,
						"offset":    msg.Offset,
					}

					if _, _, err := collection.Add(ctx, doc); err != nil {
						log.Printf("[Worker-%d] Firestore error: %v", workerID, err)
					}
				}(msg.Value)

				log.Printf(
					"[Worker-%d] Topic: %s | Partition: %d | Offset: %d | Total processed: %d",
					workerID,
					msg.Topic,
					msg.Partition,
					msg.Offset,
					current,
				)
			}
		}(w)
	}

	// Consumer handler
	handler := &ConsumerGroupHandler{Messages: messages}

	for {
		if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("Error during consumption: %v", err)
		}
		if ctx.Err() != nil {
			log.Println("Context cancelled, stopping consumer...")
			break
		}
	}
}
