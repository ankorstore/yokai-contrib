package fxelasticsearch

import (
	"github.com/ankorstore/yokai/config"
	"github.com/elastic/go-elasticsearch/v8"
)

var _ ElasticsearchClientFactory = (*DefaultElasticsearchClientFactory)(nil)

type ElasticsearchClientFactory interface {
	Create() (*elasticsearch.Client, error)
}

type DefaultElasticsearchClientFactory struct {
	config *config.Config
}

func NewDefaultElasticsearchClientFactory(config *config.Config) *DefaultElasticsearchClientFactory {
	return &DefaultElasticsearchClientFactory{
		config: config,
	}
}

func (f *DefaultElasticsearchClientFactory) Create() (*elasticsearch.Client, error) {
	// Create Elasticsearch config
	cfg := elasticsearch.Config{
		Addresses: []string{f.config.GetString("modules.elasticsearch.address")},
	}

	// Add authentication if provided
	if f.config.IsSet("modules.elasticsearch.username") && f.config.IsSet("modules.elasticsearch.password") {
		cfg.Username = f.config.GetString("modules.elasticsearch.username")
		cfg.Password = f.config.GetString("modules.elasticsearch.password")
	}

	return elasticsearch.NewClient(cfg)
}
