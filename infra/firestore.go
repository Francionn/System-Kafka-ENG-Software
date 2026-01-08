package infra

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/firestore"
)

func NewClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewFirestoreHandler(ctx context.Context, collection *firestore.CollectionRef) func(msg []byte) {
	return func(msg []byte) {
		data := map[string]interface{}{}
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Printf("Error when decoding JSON: %v", err)
			return
		}
		if _, _, err := collection.Add(ctx, data); err != nil {
			log.Printf("Error saving to Firestore: %v", err)
		}
	}
}
