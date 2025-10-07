package consumer

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
	Messages chan *sarama.ConsumerMessage
	Counter  *uint64
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("[Consumer] Session started")
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("[Consumer] Session finished")
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.Messages <- msg
		session.MarkMessage(msg, "")
	}
	return nil
}

func RunConsumerWithHandler(
	ctx context.Context,
	groupID, topic string,
	brokers []string,
	numWorkers int,
	bufferSize int,
	processMessage func(msg []byte),
) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Error create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	messages := make(chan *sarama.ConsumerMessage, bufferSize)

	var counter uint64

	// Worker pool
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			for msg := range messages {
				atomic.AddUint64(&counter, 1)
				current := atomic.LoadUint64(&counter)

				processMessage(msg.Value)

				log.Printf("[Worker-%d] Topic: %s | Partition: %d | Offset: %d | Total messages processed: %d",
					workerID, msg.Topic, msg.Partition, msg.Offset, current)
			}
		}(w)
	}

	go func() {
		for err := range consumerGroup.Errors() {
			log.Printf("[Consumer] Error: %v", err)
		}
	}()

	handler := &ConsumerGroupHandler{
		Messages: messages,
		Counter:  &counter,
	}

	for {
		if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("[Consumer] Errot at consumation: %v", err)
		}
		if ctx.Err() != nil {
			log.Println("[Consumer] Context cancel")
			close(messages)
			return
		}
	}
}
