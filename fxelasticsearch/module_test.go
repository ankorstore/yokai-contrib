package fxelasticsearch_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxelasticsearch"
	"github.com/ankorstore/yokai-contrib/fxelasticsearch/fxelasticsearchtest"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxElasticsearchModule(t *testing.T) {
	app := fxtest.New(
		t,
		fx.NopLogger,
		fxelasticsearch.FxElasticsearchModule,
		fxconfig.FxConfigModule,
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

func TestFxElasticsearchClient(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvDev)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("ELASTICSEARCH_ADDRESS", "http://localhost:9200")
	t.Setenv("ELASTICSEARCH_USERNAME", "elastic")
	t.Setenv("ELASTICSEARCH_PASSWORD", "changeme")

	var conf *config.Config
	var client *elasticsearch.Client

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxelasticsearch.FxElasticsearchModule,
		fx.Populate(&conf, &client),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create elasticsearch.Client")
	assert.NotNil(t, client)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close elasticsearch.Client")
}

func TestFxElasticsearchTestClient(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var conf *config.Config
	var client *elasticsearch.Client

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxelasticsearch.FxElasticsearchModule,
		fx.Populate(&conf, &client),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create test elasticsearch.Client")
	assert.NotNil(t, client)

	// Test that we can create mock clients independently
	mockResponse := `{"hits":{"total":{"value":1}}}`
	mockClient, err := fxelasticsearchtest.NewMockESClientWithSingleResponse(mockResponse, 200)
	assert.NoError(t, err)
	assert.NotNil(t, mockClient)

	// Verify the mock client works
	res, err := mockClient.Search()
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close test elasticsearch.Client")
}
