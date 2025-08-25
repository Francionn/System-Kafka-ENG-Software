package infra

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// new client
func NewClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no Firestore: %w", err)
	}
	return client, nil
}