package fxelasticsearch_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxelasticsearch"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestElasticsearchClientMock(t *testing.T) {
	mockClient := new(fxelasticsearch.ElasticsearchClientMock)

	// Set up expectations
	indices := []string{"test-index"}
	mockClient.On("Search", indices, mock.Anything).Return(&esapi.Response{StatusCode: 200}, nil)

	// Call the method
	res, err := mockClient.Search(indices, nil)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 200, res.StatusCode)
	mockClient.AssertExpectations(t)
}
