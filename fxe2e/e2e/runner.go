package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ankorstore/yokai-contrib/fxe2e/auth"
	"github.com/ankorstore/yokai-contrib/fxe2e/await"
	"github.com/onsi/gomega"
)

// Defaults for an AwaitSpec that omits them.
const (
	defaultAwaitTimeout  = 10 * time.Second
	defaultAwaitInterval = 25 * time.Millisecond
)

// RecordedEvent is a captured published event: its schema namespace and the
// JSON rendering of its payload. Consumers feed these in via BootResult.Events;
// the recorder helper package provides a ready-made sink.
type RecordedEvent struct {
	Schema  string
	Payload json.RawMessage
}

// JobRunner runs a named background job (instead of an HTTP request) for a case
// whose request fixture sets "job". A registry that runs jobs by name — such as
// a yokai/cron registry with a Run(ctx, name, args...) method — satisfies this
// interface directly, so a Boot can expose it without per-job wiring.
type JobRunner interface {
	Run(ctx context.Context, name string, args ...string) error
}

// BootResult is what a consumer's Boot returns: the booted HTTP handler, the
// live database, an accessor over the events recorded during the run, and the
// app-specific job runner (produced here because it captures services resolved
// from the freshly-booted DI container).
type BootResult struct {
	Server http.Handler
	DB     *sql.DB
	Events func() []RecordedEvent

	// Jobs runs the background job a case names via "job". Expose your app's job
	// registry here (anything with Run(ctx, name, args...) error). Leave nil if
	// no case uses "job".
	Jobs JobRunner

	// Redis, when set, lets a case assert Redis state via 7_redis.json. A Boot
	// using the redismem helper can pass its *miniredis.Miniredis here (it
	// satisfies RedisInspector). Leave nil if the app does not use Redis.
	Redis RedisInspector
}

// Runner orchestrates file-driven cases. Boot is required; Auth is optional.
type Runner struct {
	// Boot boots the app under test. It receives the loaded fixture (so it can
	// read reserved datasets to configure stubs) and the mock transport to wire
	// into the app's outbound HTTP client.
	Boot func(tb testing.TB, fix *Fixture, transport http.RoundTripper) BootResult

	// Auth applies the authentication described by a request's raw auth fixture.
	// When nil it defaults to auth.Ankorstore(), which dispatches on the fixture's
	// "type" field to the matching Ankorstore principal — so cases authenticate
	// purely by declaring a "type" in 2_request.json, with no wiring here.
	Auth func(tb testing.TB, req *http.Request, raw json.RawMessage)
}

// Run executes a single file-driven end-to-end test case located at caseDir.
//
// It boots the application with the case's initial data, routes outbound HTTP
// through the case's mock list, fires the case's request (or trigger), then
// asserts the response, the resulting database state and the published events —
// all declared in the case's JSON files.
func (r *Runner) Run(tb testing.TB, caseDir string) {
	tb.Helper()

	g := gomega.NewWithT(tb)
	fix := loadFixture(tb, caseDir)

	transport := newMockTransport(fix.Mocks)
	boot := r.Boot(tb, &fix, transport)

	seedDatabase(tb, boot.DB, fix.Database)

	rec := r.act(tb, fix.Request, boot)

	if fix.Request.Await != nil {
		awaitEffect(tb, boot.DB, *fix.Request.Await)
	}

	if rec != nil {
		assertResponse(g, fix.Response, rec.Code, rec.Result().Header, rec.Body.Bytes())
	}
	assertDatabase(tb, g, boot.DB, fix.Database, fix.DatabaseOut)
	assertEvents(g, boot.Events(), fix.Events)
	assertRedis(g, boot.Redis, fix.Redis)
	g.Expect(transport.unmatchedCalls()).To(gomega.BeEmpty(),
		"the request triggered outbound HTTP calls that no mock matched")
	g.Expect(transport.expectationViolations()).To(gomega.BeEmpty(),
		"one or more mocks were called a different number of times than expected")
}

// awaitEffect polls the database until the case's await query reaches its target
// count, so assertions run only after the background work has settled.
func awaitEffect(tb testing.TB, db *sql.DB, spec AwaitSpec) {
	tb.Helper()

	if spec.Query == "" {
		tb.Fatalf(`e2e: "await" must set a "query"`)
	}
	if db == nil {
		tb.Fatalf(`e2e: "await" needs a database, but Boot exposed none`)
	}

	want := 0
	if spec.Equals != nil {
		want = *spec.Equals
	}

	timeout := parseAwaitDuration(tb, "timeout", spec.Timeout, defaultAwaitTimeout)
	interval := parseAwaitDuration(tb, "interval", spec.Interval, defaultAwaitInterval)

	await.Until(tb, timeout, interval,
		fmt.Sprintf("query %q to return %d", spec.Query, want),
		await.Count(db, want, spec.Query))
}

func parseAwaitDuration(tb testing.TB, name, value string, fallback time.Duration) time.Duration {
	tb.Helper()

	if value == "" {
		return fallback
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		tb.Fatalf("e2e: invalid await %s %q: %v", name, value, err)
	}

	return d
}

// act performs the case's action: an HTTP request (default) or, when the fixture
// names a "job", that background job. It returns the HTTP recorder for a request,
// or nil for a job.
func (r *Runner) act(tb testing.TB, req RequestFixture, boot BootResult) *httptest.ResponseRecorder {
	tb.Helper()

	if req.Job == "" {
		httpReq := r.buildRequest(tb, req)
		rec := httptest.NewRecorder()
		boot.Server.ServeHTTP(rec, httpReq)

		return rec
	}

	if boot.Jobs == nil {
		tb.Fatalf("e2e: case names job %q but Boot exposed no JobRunner (BootResult.Jobs)", req.Job)
	}
	if err := boot.Jobs.Run(context.Background(), req.Job, req.Args...); err != nil {
		tb.Fatalf("e2e: job %q failed: %v", req.Job, err)
	}

	return nil
}

func (r *Runner) buildRequest(tb testing.TB, req RequestFixture) *http.Request {
	tb.Helper()

	var body *bytes.Reader
	switch {
	case req.BodyString != "":
		body = bytes.NewReader([]byte(req.BodyString))
	case len(req.Body) > 0:
		body = bytes.NewReader(req.Body)
	default:
		body = bytes.NewReader(nil)
	}

	method := req.Method
	if method == "" {
		method = http.MethodGet
	}

	httpReq := httptest.NewRequest(method, req.Path, body)

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	authFn := r.Auth
	if authFn == nil {
		authFn = auth.Ankorstore()
	}
	authFn(tb, httpReq, req.Auth)

	return httpReq
}
