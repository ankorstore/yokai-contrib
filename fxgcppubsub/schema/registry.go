package schema

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
)

var _ SchemaConfigRegistry = (*DefaultSchemaConfigRegistry)(nil)

// SchemaConfigRegistry is the interface for schema config registries.
type SchemaConfigRegistry interface {
	Get(ctx context.Context, schemaID string) (*pubsub.SchemaConfig, error)
}

// DefaultSchemaConfigRegistry is the default SchemaConfigRegistry implementation.
type DefaultSchemaConfigRegistry struct {
	client  *pubsub.SchemaClient
	schemas map[string]*pubsub.SchemaConfig
	mutex   sync.RWMutex
}

// NewDefaultSchemaConfigRegistry returns a new DefaultSchemaConfigRegistry instance.
func NewDefaultSchemaConfigRegistry(client *pubsub.SchemaClient) *DefaultSchemaConfigRegistry {
	return &DefaultSchemaConfigRegistry{
		client:  client,
		schemas: make(map[string]*pubsub.SchemaConfig),
	}
}

// Get gets a pubsub.SchemaConfig for a provided schemaID.
func (r *DefaultSchemaConfigRegistry) Get(ctx context.Context, schemaID string) (*pubsub.SchemaConfig, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	schemaID = NormalizeSchemaID(schemaID)

	if _, found := r.schemas[schemaID]; found {
		return r.schemas[schemaID], nil
	}

	schema, err := r.client.Schema(ctx, schemaID, pubsub.SchemaViewFull)
	if err != nil {
		return nil, fmt.Errorf("cannot get schema configuration: %w", err)
	}

	r.schemas[schemaID] = schema

	return schema, nil
}
