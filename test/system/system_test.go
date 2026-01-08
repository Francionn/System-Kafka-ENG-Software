package infra_test

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"kafka_system/infra"

	"cloud.google.com/go/firestore"
)

func TestSystemEndToEnd(t *testing.T) {
	ctx := context.Background()
	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	if projectID == "" {
		t.Skip("FIRESTORE_PROJECT_ID no defined")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("Error create client Firestore: %v", err)
	}
	defer client.Close()

	collection := client.Collection("test-system-collection")

	handler := infra.NewFirestoreHandler(ctx, collection)

	msg := map[string]interface{}{
		"sensor_id": "abc123",
		"temp":      26.5,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	jsonMsg, _ := json.Marshal(msg)

	handler(jsonMsg)

	time.Sleep(2 * time.Second)

	iter := collection.Where("sensor_id", "==", "abc123").Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		t.Fatalf("Error searching for documents: %v", err)
	}

	if len(docs) == 0 {
		t.Errorf("documents not found in Firestore")
	} else {
		log.Printf("System test successful — document saved: %+v", docs[0].Data())
	}
}
