package fxelasticsearch_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxelasticsearch"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestElasticsearchClientMock(t *testing.T) {
	t.Run("Search success", func(t *testing.T) {
		mockClient := new(fxelasticsearch.ElasticsearchClientMock)

		// Set up expectations - success case
		indices := []string{"test-index"}
		mockClient.On("Search", indices, mock.Anything).Return(&esapi.Response{StatusCode: 200}, nil)

		// Call the method
		res, err := mockClient.Search(indices, nil)

		// Assert expectations
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 200, res.StatusCode)
		mockClient.AssertExpectations(t)
	})

	t.Run("Search with nil response", func(t *testing.T) {
		mockClient := new(fxelasticsearch.ElasticsearchClientMock)

		// Set up expectations - nil response case
		indices := []string{"test-index"}
		expectedErr := assert.AnError
		mockClient.On("Search", indices, mock.Anything).Return(nil, expectedErr)

		// Call the method
		res, err := mockClient.Search(indices, nil)

		// Assert expectations
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, res)
		mockClient.AssertExpectations(t)
	})

	t.Run("Search with invalid response type", func(t *testing.T) {
		mockClient := new(fxelasticsearch.ElasticsearchClientMock)

		// Set up expectations - invalid response type
		indices := []string{"test-index"}
		expectedErr := assert.AnError
		// Return a string instead of *esapi.Response to trigger type assertion failure
		mockClient.On("Search", indices, mock.Anything).Return("not a response", expectedErr)

		// Call the method
		res, err := mockClient.Search(indices, nil)

		// Assert expectations
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, res)
		mockClient.AssertExpectations(t)
	})

	t.Run("ExampleMethod", func(t *testing.T) {
		mockClient := new(fxelasticsearch.ElasticsearchClientMock)

		// Set up expectations
		expectedErr := assert.AnError
		mockClient.On("ExampleMethod").Return(expectedErr)

		// Call the method
		err := mockClient.ExampleMethod()

		// Assert expectations
		assert.Equal(t, expectedErr, err)
		mockClient.AssertExpectations(t)
	})
}
