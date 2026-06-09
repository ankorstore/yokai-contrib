package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// mockTransport is an http.RoundTripper that answers outbound requests from the
// fixture's mock list. Any request that matches no mock is recorded as an
// unmatched call and answered with 599 so the caller fails fast; the harness
// then asserts that no unmatched calls occurred. It also counts how many times
// each mock was matched so a case can assert exact call counts.
type mockTransport struct {
	mocks []MockFixture

	mu        sync.Mutex
	hits      []int
	unmatched []string
}

func newMockTransport(mocks []MockFixture) *mockTransport {
	return &mockTransport{mocks: mocks, hits: make([]int, len(mocks))}
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body := drainBody(req)

	for i := range t.mocks {
		if matches(t.mocks[i].Match, req, body) {
			t.mu.Lock()
			t.hits[i]++
			t.mu.Unlock()

			return buildResponse(req, t.mocks[i].Respond), nil
		}
	}

	t.mu.Lock()
	t.unmatched = append(t.unmatched, fmt.Sprintf("%s %s", req.Method, req.URL.String()))
	t.mu.Unlock()

	return &http.Response{
		StatusCode: 599,
		Status:     "599 Unmatched E2E Mock",
		Body:       io.NopCloser(strings.NewReader(`{"error":"no e2e mock matched this request"}`)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Request:    req,
	}, nil
}

func (t *mockTransport) unmatchedCalls() []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	return append([]string(nil), t.unmatched...)
}

// expectationViolations returns a human-readable message for every mock whose
// ExpectedCalls constraint was not met. Mocks without an ExpectedCalls are
// unconstrained.
func (t *mockTransport) expectationViolations() []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	var out []string
	for i := range t.mocks {
		want := t.mocks[i].ExpectedCalls
		if want == nil {
			continue
		}
		if t.hits[i] != *want {
			out = append(out, fmt.Sprintf("mock #%d (%s): expected %d call(s), got %d",
				i, describeMatch(t.mocks[i].Match), *want, t.hits[i]))
		}
	}

	return out
}

// drainBody reads and replaces req.Body so it can be matched against and still
// be read again downstream. Returns nil when there is no body.
func drainBody(req *http.Request) []byte {
	if req.Body == nil {
		return nil
	}

	data, err := io.ReadAll(req.Body)
	_ = req.Body.Close()
	if err != nil {
		return nil
	}
	req.Body = io.NopCloser(bytes.NewReader(data))

	return data
}

func matches(m MockMatch, req *http.Request, body []byte) bool {
	if m.Method != "" && !strings.EqualFold(m.Method, req.Method) {
		return false
	}
	if m.Host != "" && m.Host != req.URL.Host {
		return false
	}
	if m.Path != "" && m.Path != req.URL.Path {
		return false
	}
	if m.PathPrefix != "" && !strings.HasPrefix(req.URL.Path, m.PathPrefix) {
		return false
	}
	if !queryMatches(m.Query, req) {
		return false
	}
	if !headerMatches(m.Headers, req) {
		return false
	}
	if !bodyMatches(m.Body, body) {
		return false
	}

	return true
}

func queryMatches(want map[string]string, req *http.Request) bool {
	if len(want) == 0 {
		return true
	}

	q := req.URL.Query()
	for k, v := range want {
		if q.Get(k) != v {
			return false
		}
	}

	return true
}

func headerMatches(want map[string]string, req *http.Request) bool {
	for k, v := range want {
		if req.Header.Get(k) != v {
			return false
		}
	}

	return true
}

// bodyMatches reports whether the outbound body satisfies the expected JSON
// subset. An empty expectation matches anything; a non-JSON body never matches a
// JSON expectation.
func bodyMatches(want json.RawMessage, body []byte) bool {
	if len(want) == 0 {
		return true
	}

	var wantRoot, gotRoot any
	if err := json.Unmarshal(want, &wantRoot); err != nil {
		return false
	}
	if err := json.Unmarshal(body, &gotRoot); err != nil {
		return false
	}

	return len(subsetMatch(wantRoot, gotRoot, "$")) == 0
}

func describeMatch(m MockMatch) string {
	var parts []string
	if m.Method != "" {
		parts = append(parts, m.Method)
	}
	switch {
	case m.Path != "":
		parts = append(parts, m.Path)
	case m.PathPrefix != "":
		parts = append(parts, m.PathPrefix+"*")
	}
	if len(parts) == 0 {
		return "any request"
	}

	return strings.Join(parts, " ")
}

func buildResponse(req *http.Request, r MockRespond) *http.Response {
	status := r.Status
	if status == 0 {
		status = http.StatusOK
	}

	contentType := r.ContentType
	if contentType == "" {
		contentType = "application/json"
	}

	var body []byte
	switch {
	case r.BodyString != "":
		body = []byte(r.BodyString)
	case len(r.Body) > 0:
		body = []byte(r.Body)
	default:
		body = []byte("{}")
	}

	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header: http.Header{
			"Content-Type": []string{contentType},
			// Satisfies the go-elasticsearch client's product check so the same
			// mock transport can answer Elasticsearch calls (see the esmock
			// helper). Harmless for non-ES responses.
			"X-Elastic-Product": []string{"Elasticsearch"},
		},
		Request: req,
	}
}
