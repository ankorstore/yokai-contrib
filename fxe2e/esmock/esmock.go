// Package esmock routes an application's Elasticsearch client through the
// end-to-end harness's HTTP mock transport.
//
// Yokai's fxelasticsearch module, in test mode, swaps in a mock client that
// returns a single canned empty response — which forces applications to bypass
// their real Elasticsearch-backed code with stubs. This package instead points
// the *elasticsearch.Client at the harness's shared mock transport, so
// Elasticsearch becomes just another external HTTP dependency: the genuine
// query-building and response-parsing code runs, and each case declares the
// Elasticsearch responses it expects in 3_mocks.json alongside every other
// outbound call.
//
// Wire it into a Boot function with the same transport passed in by the runner:
//
//	internal.RunE2ETest(tb,
//	    esmock.FxOption(transport),
//	    // ... other options
//	)
//
// A mock entry then answers the search, e.g.:
//
//	{
//	  "match":   { "pathPrefix": "/my_index/_search" },
//	  "respond": { "status": 200, "body": { "hits": { "total": {"value":1}, "hits": [ ... ] } } }
//	}
//
// The harness's mock responses carry the X-Elastic-Product header the
// go-elasticsearch client requires, so its product check passes.
package esmock

import (
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/fx"
)

// mockAddr is an arbitrary, syntactically-valid Elasticsearch address. Requests
// never leave the process — the mock transport intercepts them — so the host is
// irrelevant; match on path (or pathPrefix) in 3_mocks.json.
const mockAddr = "http://elasticsearch.e2e.mock:9200"

// FxOption returns an fx.Option that decorates the application's
// *elasticsearch.Client so all its requests go through transport.
func FxOption(transport http.RoundTripper) fx.Option {
	return fx.Decorate(func(*elasticsearch.Client) (*elasticsearch.Client, error) {
		return elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{mockAddr},
			Transport: transport,
		})
	})
}
