package fxelasticsearch_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxelasticsearch"
	"github.com/ankorstore/yokai/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
)

func TestDefaultElasticsearchClientFactory(t *testing.T) {
	createConfig := func() (*config.Config, error) {
		return config.NewDefaultConfigFactory().Create(
			config.WithFilePaths("./testdata/config"),
		)
	}

	t.Run("create success", func(t *testing.T) {
		t.Setenv("ELASTICSEARCH_ADDRESS", "http://localhost:9200")

		cfg, err := createConfig()
		assert.NoError(t, err)

		factory := fxelasticsearch.NewDefaultElasticsearchClientFactory(cfg)

		client, err := factory.Create()
		assert.NoError(t, err)
		assert.IsType(t, &elasticsearch.Client{}, client)
	})

	t.Run("create with auth", func(t *testing.T) {
		t.Setenv("ELASTICSEARCH_ADDRESS", "http://localhost:9200")
		t.Setenv("ELASTICSEARCH_USERNAME", "elastic")
		t.Setenv("ELASTICSEARCH_PASSWORD", "changeme")

		cfg, err := createConfig()
		assert.NoError(t, err)

		factory := fxelasticsearch.NewDefaultElasticsearchClientFactory(cfg)

		client, err := factory.Create()
		assert.NoError(t, err)
		assert.IsType(t, &elasticsearch.Client{}, client)
	})
}
