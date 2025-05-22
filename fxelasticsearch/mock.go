package fxelasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/mock"
)

// ElasticsearchClientInterface defines the methods for the Elasticsearch client.
// Extend this interface as needed.
type ElasticsearchClientInterface interface {
	Search(indices []string, body interface{}) (*esapi.Response, error)
}

// ElasticsearchClientMock implements ElasticsearchClientInterface.
// Add methods as needed for testing.
type ElasticsearchClientMock struct {
	mock.Mock
}

func (m *ElasticsearchClientMock) Search(indices []string, body interface{}) (*esapi.Response, error) {
	args := m.Called(indices, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	resp, ok := args.Get(0).(*esapi.Response)
	if !ok {
		return nil, args.Error(1)
	}

	return resp, args.Error(1)
}

// Ensure ElasticsearchClientMock implements ElasticsearchClientInterface.
var _ ElasticsearchClientInterface = (*ElasticsearchClientMock)(nil)
