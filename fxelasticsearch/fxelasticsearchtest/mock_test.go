package fxelasticsearchtest_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxelasticsearch/fxelasticsearchtest"
	"github.com/stretchr/testify/assert"
)

// Test the new HTTP transport-level mocking.
func TestMockTransport(t *testing.T) {
	t.Run("single response", func(t *testing.T) {
		responses := []fxelasticsearchtest.MockResponse{
			{
				StatusCode:   200,
				ResponseBody: `{"status":"ok"}`,
			},
		}

		transport := fxelasticsearchtest.NewMockTransport(responses)

		// Create a mock request
		req, err := http.NewRequest(http.MethodGet, "http://localhost:9200", nil)
		assert.NoError(t, err)

		// Test the transport
		resp, err := transport.RoundTrip(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		assert.Equal(t, "Elasticsearch", resp.Header.Get("X-Elastic-Product"))

		// Read body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"status":"ok"}`, string(body))
	})

	t.Run("multiple responses", func(t *testing.T) {
		responses := []fxelasticsearchtest.MockResponse{
			{StatusCode: 200, ResponseBody: `{"response":1}`},
			{StatusCode: 201, ResponseBody: `{"response":2}`},
			{StatusCode: 202, ResponseBody: `{"response":3}`},
		}

		transport := fxelasticsearchtest.NewMockTransport(responses)
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:9200", nil)

		// First request
		resp1, err := transport.RoundTrip(req)
		assert.NoError(t, err)
		resp1.Body.Close()
		assert.Equal(t, 200, resp1.StatusCode)

		// Second request
		resp2, err := transport.RoundTrip(req)
		assert.NoError(t, err)
		resp2.Body.Close()
		assert.Equal(t, 201, resp2.StatusCode)

		// Third request
		resp3, err := transport.RoundTrip(req)
		assert.NoError(t, err)
		resp3.Body.Close()
		assert.Equal(t, 202, resp3.StatusCode)

		// Fourth request (should repeat last response)
		resp4, err := transport.RoundTrip(req)
		assert.NoError(t, err)
		resp4.Body.Close()
		assert.Equal(t, 202, resp4.StatusCode)
	})

	t.Run("error response", func(t *testing.T) {
		expectedErr := errors.New("connection failed")
		responses := []fxelasticsearchtest.MockResponse{
			{Error: expectedErr},
		}

		transport := fxelasticsearchtest.NewMockTransport(responses)
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:9200", nil)

		resp, err := transport.RoundTrip(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, resp)
	})
}

func TestNewMockESClient(t *testing.T) {
	t.Run("successful search with single response", func(t *testing.T) {
		mockResponse := `{
			"took": 5,
			"hits": {
				"total": {"value": 1},
				"hits": [
					{
						"_source": {"title": "test document"}
					}
				]
			}
		}`

		client, err := fxelasticsearchtest.NewMockESClientWithSingleResponse(mockResponse, 200)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		// Perform a search
		res, err := client.Search(
			client.Search.WithContext(context.Background()),
			client.Search.WithIndex("test-index"),
			client.Search.WithBody(strings.NewReader(`{"query":{"match_all":{}}}`)),
		)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
		assert.False(t, res.IsError())

		// Read and verify response body
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(body), "test document")
	})

	t.Run("multiple requests", func(t *testing.T) {
		responses := []fxelasticsearchtest.MockResponse{
			{StatusCode: 200, ResponseBody: `{"hits":{"total":{"value":1}}}`},
			{StatusCode: 200, ResponseBody: `{"hits":{"total":{"value":2}}}`},
		}

		client, err := fxelasticsearchtest.NewMockESClient(responses)
		assert.NoError(t, err)

		// First search
		res1, err := client.Search()
		assert.NoError(t, err)
		defer res1.Body.Close()
		body1, _ := io.ReadAll(res1.Body)
		assert.Contains(t, string(body1), `"value":1`)

		// Second search
		res2, err := client.Search()
		assert.NoError(t, err)
		defer res2.Body.Close()
		body2, _ := io.ReadAll(res2.Body)
		assert.Contains(t, string(body2), `"value":2`)
	})

	t.Run("error handling", func(t *testing.T) {
		expectedErr := errors.New("network error")
		client, err := fxelasticsearchtest.NewMockESClientWithError(expectedErr)
		assert.NoError(t, err)

		res, err := client.Search()
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, res)
	})

	t.Run("elasticsearch error response", func(t *testing.T) {
		errorResponse := `{
			"error": {
				"type": "index_not_found_exception",
				"reason": "no such index [missing]"
			}
		}`

		client, err := fxelasticsearchtest.NewMockESClientWithSingleResponse(errorResponse, 404)
		assert.NoError(t, err)

		res, err := client.Search(client.Search.WithIndex("missing"))
		assert.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, 404, res.StatusCode)
		assert.True(t, res.IsError())
	})
}

func TestMockTransportEdgeCases(t *testing.T) {
	t.Run("empty responses defaults to success", func(t *testing.T) {
		transport := fxelasticsearchtest.NewMockTransport([]fxelasticsearchtest.MockResponse{})
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:9200", nil)

		resp, err := transport.RoundTrip(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "{}", string(body))
	})
}
