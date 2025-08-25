package producer

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	ctx     context.Context
	client  *firestore.Client
	writer  *kafka.Writer
	topic   string
	project string
}
