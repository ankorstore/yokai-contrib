package fxelasticsearch

import (
	"github.com/ankorstore/yokai/config"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "elasticsearch"

// FxElasticsearchModule is the [Fx] elasticsearch module.
//
// [Fx]: https://github.com/uber-go/fx
var FxElasticsearchModule = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(NewDefaultElasticsearchClientFactory, fx.As(new(ElasticsearchClientFactory))),
		NewFxElasticsearchClient,
	),
)

// FxElasticsearchClientParam allows injection of the required dependencies in [NewElasticsearchClient].
type FxElasticsearchClientParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
	Factory   ElasticsearchClientFactory
}

// NewFxElasticsearchClient returns a [elasticsearch.Client].
// In test environment, it returns a default mock client with basic functionality.
// In production, it returns a real client connected to Elasticsearch.
//
// For advanced testing scenarios, use NewMockESClient, NewMockESClientWithSingleResponse,
// or NewMockESClientWithError from this package to create custom mock clients.
func NewFxElasticsearchClient(p FxElasticsearchClientParam) (*elasticsearch.Client, error) {
	if p.Config.IsTestEnv() {
		// In test environment, provide a default mock client that returns empty successful responses
		// This allows basic functionality to work out of the box in tests
		defaultResponse := `{"took":1,"timed_out":false,"hits":{"total":{"value":0},"hits":[]}}`

		return NewMockESClientWithSingleResponse(defaultResponse, 200)
	}
	// In production, use the factory to create a real client
	client, err := p.Factory.Create()

	return client, err
}
