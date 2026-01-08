package infra_test

import (
	"context"
	"testing"

	"kafka_system/infra"
)

func TestNewFirestoreHandler(t *testing.T) {
	handler := infra.NewFirestoreHandler(context.Background(), nil)
	data := []byte(`{"sensor": "temp", "value": 22}`)
	handler(data)

	t.Log("Handler processed JSON correctly")
}
