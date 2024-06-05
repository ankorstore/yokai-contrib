package schema

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"
)

type SchemaRegistry struct {
	client  *pubsub.SchemaClient
	schemas map[string]*pubsub.SchemaConfig
	mutex   sync.RWMutex
}

func NewSchemaRegistry(client *pubsub.SchemaClient) *SchemaRegistry {
	return &SchemaRegistry{
		client:  client,
		schemas: make(map[string]*pubsub.SchemaConfig),
	}
}

func (r *SchemaRegistry) Get(ctx context.Context, schemaID string) (*pubsub.SchemaConfig, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.schemas[schemaID]; found {
		return r.schemas[schemaID], nil
	}

	schema, err := r.client.Schema(ctx, schemaID, pubsub.SchemaViewFull)
	if err != nil {
		return nil, err
	}

	r.schemas[schemaID] = schema

	return schema, nil
}
