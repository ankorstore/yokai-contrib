package fxelasticsearchtest

import (
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

// MockResponse represents a single mock HTTP response.
type MockResponse struct {
	StatusCode   int
	ResponseBody string
	Error        error
}

// MockTransport provides HTTP transport mock for testing Elasticsearch client.
// This allows testing with the real elasticsearch.Client API without interface constraints.
type MockTransport struct {
	responses []MockResponse
	index     int
}

// NewMockTransport creates a new MockTransport with the given responses.
// Responses are returned in order. If more requests are made than responses provided,
// the last response is repeated.
func NewMockTransport(responses []MockResponse) *MockTransport {
	if len(responses) == 0 {
		responses = []MockResponse{{StatusCode: 200, ResponseBody: "{}"}}
	}

	return &MockTransport{
		responses: responses,
		index:     0,
	}
}

// RoundTrip implements the http.RoundTripper interface.
func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Get current response (or last one if we've exceeded the list)
	responseIndex := t.index
	if responseIndex >= len(t.responses) {
		responseIndex = len(t.responses) - 1
	} else {
		t.index++
	}

	mockResp := t.responses[responseIndex]

	if mockResp.Error != nil {
		return nil, mockResp.Error
	}

	// Create HTTP response
	response := &http.Response{
		StatusCode: mockResp.StatusCode,
		Body:       io.NopCloser(strings.NewReader(mockResp.ResponseBody)),
		Header:     make(http.Header),
		Request:    req,
	}

	// Add Elasticsearch-specific headers that the client expects
	response.Header.Set("Content-Type", "application/json")
	response.Header.Set("X-Elastic-Product", "Elasticsearch") // Critical header for client validation

	return response, nil
}

// NewMockESClient creates an Elasticsearch client with mocked HTTP transport.
// This allows you to test with the real elasticsearch.Client API while controlling responses.
func NewMockESClient(responses []MockResponse) (*elasticsearch.Client, error) {
	transport := NewMockTransport(responses)

	cfg := elasticsearch.Config{
		Transport: transport,
	}

	return elasticsearch.NewClient(cfg)
}

// NewMockESClientWithSingleResponse creates an Elasticsearch client with a single mock response.
// This is a convenience function for simple test cases.
func NewMockESClientWithSingleResponse(responseBody string, statusCode int) (*elasticsearch.Client, error) {
	responses := []MockResponse{
		{
			StatusCode:   statusCode,
			ResponseBody: responseBody,
		},
	}

	return NewMockESClient(responses)
}

// NewMockESClientWithError creates an Elasticsearch client that returns an error on requests.
// This is useful for testing error handling.
func NewMockESClientWithError(err error) (*elasticsearch.Client, error) {
	responses := []MockResponse{
		{
			Error: err,
		},
	}

	return NewMockESClient(responses)
}
