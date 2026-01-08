package integration

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"testing"
	"time"

	"kafka_system/consumer"

	"cloud.google.com/go/firestore"
)

type mockFirestoreClient struct{}

func (m *mockFirestoreClient) Collection(path string) *firestore.CollectionRef {
	return &firestore.CollectionRef{}
}

func mockProcessMessage(t *testing.T, wg *sync.WaitGroup, processed *bool) func([]byte) {
	return func(msg []byte) {
		defer wg.Done()

		data := map[string]interface{}{}
		if err := json.Unmarshal(msg, &data); err != nil {
			t.Errorf("Erro ao decodificar JSON: %v", err)
			return
		}
		log.Printf("[Integration Test] Mensagem recebida: %+v", data)
		*processed = true
	}
}

func TestKafkaToFirestoreIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	processed := false

	process := mockProcessMessage(t, &wg, &processed)

	mockClient := &firestore.Client{}

	go func() {
		brokers := []string{"localhost:9094"}
		groupID := "integration-group-test"
		topic := "my-topic-test-new"
		firestoreCollection := "messages"

		consumer.RunConsumerWithHandler(
			ctx,
			groupID,
			topic,
			brokers,
			1,
			10,
			process,
			mockClient,
			firestoreCollection,
		)
	}()

	wg.Wait()

	if !processed {
		t.Errorf("Integration failure: the consumer did not process the message")
	} else {
		t.Log("The consumer integrated correctly with the handler (mock)")
	}
}
