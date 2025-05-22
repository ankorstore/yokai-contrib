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

// NewFxElasticsearchClient returns a [elasticsearch.Client] and a general mock in test mode.
func NewFxElasticsearchClient(p FxElasticsearchClientParam) (*elasticsearch.Client, ElasticsearchClientInterface, error) {
	if p.Config.IsTestEnv() {
		mockClient := &elasticsearch.Client{}
		generalMock := &ElasticsearchClientMock{}

		return mockClient, generalMock, nil
	} else {
		client, err := p.Factory.Create()

		return client, nil, err
	}
}
