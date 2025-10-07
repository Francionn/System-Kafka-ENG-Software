// Não faz parte do sistema no momento

package producer

import (
	"context"
	"fmt"
	"log"

	"kafka_system/infra"

	"cloud.google.com/go/firestore"
	"github.com/segmentio/kafka-go"
)

// NewProducer create new producer
func NewProducer(ctx context.Context, projectID, topic string, brokers []string) (*Producer, error) {
	client, err := infra.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// create writer to Kafka
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &Producer{
		ctx:     ctx,
		client:  client,
		writer:  writer,
		topic:   topic,
		project: projectID,
	}, nil
}

// run firestore and send to Kafka
func (p *Producer) Run(collection string) {
	defer p.client.Close()
	defer p.writer.Close()

	fmt.Println("Producer started")

	colRef := p.client.Collection(collection)
	snapshots := colRef.Snapshots(p.ctx)
	defer snapshots.Stop()

	firstSnapshot := true

	for {
		snap, err := snapshots.Next()
		if err != nil {
			log.Printf("Error in snapshot: %v", err)
			continue
		}

		// real time
		if firstSnapshot {
			firstSnapshot = false
			continue
		}

		if len(snap.Changes) == 0 {
			p.sendMessage("none")
			continue
		}

		for _, change := range snap.Changes {
			if change.Kind == firestore.DocumentAdded || change.Kind == firestore.DocumentModified {
				value := fmt.Sprintf("%v", change.Doc.Data())
				p.sendMessage(value)
			}
		}
	}
}

func (p *Producer) sendMessage(msg string) {
	err := p.writer.WriteMessages(p.ctx, kafka.Message{
		Value: []byte(msg),
	})
	if err != nil {
		log.Printf("Error send message: %v", err)
	} else {
		fmt.Println("send:", msg)
	}
}

// fmt.Println("Enviado:", msg)
