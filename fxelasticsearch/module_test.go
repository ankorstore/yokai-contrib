package fxelasticsearch_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxelasticsearch"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	var mockClient fxelasticsearch.ElasticsearchClientInterface

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxelasticsearch.FxElasticsearchModule,
		fx.Populate(&conf, &client, &mockClient),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create elasticsearch.Client")
	assert.NotNil(t, client)
	assert.Nil(t, mockClient)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close elasticsearch.Client")
}

func TestFxElasticsearchTestClient(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("ELASTICSEARCH_ADDRESS", "http://localhost:9200")
	t.Setenv("ELASTICSEARCH_USERNAME", "elastic")
	t.Setenv("ELASTICSEARCH_PASSWORD", "changeme")

	var conf *config.Config
	var client *elasticsearch.Client
	var mockClient fxelasticsearch.ElasticsearchClientInterface

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxelasticsearch.FxElasticsearchModule,
		fx.Populate(&conf, &client, &mockClient),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create test elasticsearch.Client")
	assert.NotNil(t, client)
	assert.NotNil(t, mockClient)

	// Test mock functionality
	mockEs, ok := mockClient.(*fxelasticsearch.ElasticsearchClientMock)
	assert.True(t, ok, "mockClient should be of type *ElasticsearchClientMock")

	indices := []string{"test-index"}
	mockEs.On("Search", indices, mock.Anything).Return(&esapi.Response{StatusCode: 200}, nil)

	// Use the mock
	response, err := mockEs.Search(indices, map[string]interface{}{"query": "test"})
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	// Verify expectations were met
	mockEs.AssertExpectations(t)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close test elasticsearch.Client")
}
