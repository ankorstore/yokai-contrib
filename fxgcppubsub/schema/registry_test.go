package schema_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/schema"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestDefaultSchemaConfigRegistry(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var client *pubsub.SchemaClient
	var registry schema.SchemaConfigRegistry

	ctx := context.Background()

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxgcppubsub.FxGcpPubSubModule,
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
		fx.Populate(&client, &registry),
	).RequireStart().RequireStop()

	t.Run("get non existing schema", func(t *testing.T) {
		_, err := registry.Get(ctx, "test-schema")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), `schema("projects/test-project/schemas/test-schema") not found`)
	})

	t.Run("get schema from schema client", func(t *testing.T) {
		_, err := client.CreateSchema(ctx, "test-schema", pubsub.SchemaConfig{
			Name: "test-schema",
		})
		assert.NoError(t, err)

		schemaConfig, err := registry.Get(ctx, "test-schema")
		assert.NoError(t, err)
		assert.Equal(t, "projects/test-project/schemas/test-schema", schemaConfig.Name)
	})

	t.Run("get schema from cache", func(t *testing.T) {
		_, err := client.CreateSchema(ctx, "test-schema", pubsub.SchemaConfig{
			Name: "test-schema",
		})
		assert.NoError(t, err)

		schemaConfig, err := registry.Get(ctx, "test-schema")
		assert.NoError(t, err)
		assert.Equal(t, "projects/test-project/schemas/test-schema", schemaConfig.Name)

		err = client.DeleteSchema(ctx, "test-schema")
		assert.NoError(t, err)

		schemaConfig2, err := registry.Get(ctx, "test-schema")
		assert.NoError(t, err)
		assert.Equal(t, "projects/test-project/schemas/test-schema", schemaConfig2.Name)

		assert.Equal(t, schemaConfig, schemaConfig2)
	})
}
