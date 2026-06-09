// Package e2e provides a declarative, file-driven end-to-end test harness for
// Yokai-based HTTP services.
//
// Test cases are organised as testdata/<route>/<method>/<test_name>/, and every
// case directory MUST contain all seven numbered JSON files — a case missing any
// of them fails rather than running partially:
//
//	1_database.json      initial DB state: {"<table>": [ {col: val, ...}, ... ]}
//	2_request.json       the request (method, path, headers, auth, body) or job
//	3_mocks.json         external HTTP call mocks (match -> canned response)
//	4_database_out.json  expected DB state after the request (subset + ignore)
//	5_job.json           events expected to be published
//	6_response.json      expected response (status, body subset, ignore paths)
//	7_redis.json         expected Redis state (a bare {} when none)
//
// The harness boots the real application via a consumer-supplied Boot function
// and only stubs the outside world: outbound HTTP is routed through a mock
// transport, and published events are recorded. Assertions use Gomega with
// subset matching — an expected file only needs to describe the fields it cares
// about, and volatile values (timestamps, generated IDs) can be dropped via an
// ignore list.
//
// Everything application-specific (how to boot the app, how to authenticate a
// request, how to run a named background job) is injected through the Runner and
// the BootResult it receives, so the engine itself stays free of any application
// or domain type. Waiting for background work is declared as data via a case's
// "await" query — no per-app wiring.
package e2e

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Fixture is the fully-loaded set of files for a single test case.
type Fixture struct {
	Database    map[string][]Row // initial DB state, keyed by table
	Request     RequestFixture
	Response    ResponseFixture
	Mocks       []MockFixture
	DatabaseOut DatabaseOutFixture
	Events      EventsFixture
	Redis       *RedisFixture // expected Redis state; nil when 7_redis.json is absent
}

// FixedTime returns the deterministic "now" a case declared via the request's
// optional "now" field (RFC 3339), so a Boot can decorate the app's clock for a
// reproducible run. The bool is false when the case declares no fixed time.
func (f *Fixture) FixedTime() (time.Time, bool) {
	if f.Request.Now == "" {
		return time.Time{}, false
	}

	t, err := time.Parse(time.RFC3339, f.Request.Now)
	if err != nil {
		return time.Time{}, false
	}

	return t, true
}

// Row is a single table row: column name -> value. JSON numbers are decoded as
// json.Number so integer columns are not corrupted into floats.
type Row = map[string]any

// ReservedKeyPrefix marks keys in 1_database.json that are NOT SQL tables but
// harness datasets a consumer's Boot may read to configure stubs (e.g. a
// "__contacts__" dataset feeding a stubbed search client). The seeder skips any
// table whose name carries this prefix.
const ReservedKeyPrefix = "__"

// RequestFixture describes how a case acts on the system.
//
// Job selects the action:
//
//	""    — send the HTTP request described by Method/Path/Headers/Auth/Body.
//	other — run the named background job via the Boot-provided JobRunner
//	        (BootResult.Jobs), passing Args. The HTTP fields are then ignored
//	        and 6_response.json is not asserted.
//
// Await, when set, makes the harness poll the database after the action until a
// COUNT(*) query reaches its target — for journeys whose effects complete in a
// background goroutine. It needs no production-code change: the wait condition
// is declared as data here, not wired in the app.
type RequestFixture struct {
	Job   string     `json:"job"`
	Args  []string   `json:"args"`
	Await *AwaitSpec `json:"await"`

	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Auth    json.RawMessage   `json:"auth"` // opaque to the engine; decoded by the consumer's Auth func
	Body    json.RawMessage   `json:"body"`

	// BodyString is a raw, non-JSON request body (e.g. form-encoded or plain
	// text). When set it is sent verbatim and Body is ignored; pair it with a
	// Content-Type in Headers.
	BodyString string `json:"bodyString"`

	// Now is an optional fixed wall-clock time (RFC 3339) for the case, surfaced
	// to Boot via Fixture.FixedTime so the app's clock can run deterministically.
	Now string `json:"now"`

	// Params carries job-specific data for custom (non-HTTP) jobs, kept here so
	// the rest of the fixture can be strictly decoded.
	Params json.RawMessage `json:"params"`

	// Raw is the verbatim content of 2_request.json, made available to custom
	// jobs that need fields beyond the ones the engine decodes.
	Raw json.RawMessage `json:"-"`
}

// AwaitSpec declares how to wait for a background effect: poll Query (a COUNT(*)
// statement) until it returns Equals (default 0), or fail after Timeout. It
// observes the database the action affects, so no production-code change is
// needed — the wait condition is test data, not app wiring.
type AwaitSpec struct {
	Query    string `json:"query"`    // a COUNT(*) statement, e.g. "SELECT COUNT(*) FROM jobs WHERE status = 'running'"
	Equals   *int   `json:"equals"`   // target count; defaults to 0
	Timeout  string `json:"timeout"`  // e.g. "10s"; defaults to 10s
	Interval string `json:"interval"` // e.g. "25ms"; defaults to 25ms
}

type ResponseFixture struct {
	// Status is the expected HTTP status code. It is mandatory for HTTP cases:
	// a response fixture that omits it (or sets it to 0) fails rather than
	// silently skipping the status assertion.
	Status int `json:"status"`
	// Headers are expected response headers, matched as a subset: every header
	// named here must be present with the given value (names are
	// case-insensitive); response headers not listed here are ignored.
	Headers map[string]string `json:"headers"`
	Body    json.RawMessage   `json:"body"`
	// BodyString asserts a raw, non-JSON response body by exact string equality
	// (after trimming surrounding whitespace). When set, Body/Ignore/Unordered —
	// which are JSON-aware — are not used.
	BodyString string `json:"bodyString"`
	// Ignore lists dotted JSON paths to drop from both expected and actual
	// body before comparison, e.g. "data.attributes.createdAt".
	Ignore []string `json:"ignore"`
	// Unordered lists dotted JSON paths whose arrays are matched
	// order-independently (length still must be equal; each expected element
	// must match some distinct actual element). Use for list endpoints whose
	// ordering is not guaranteed.
	Unordered []string `json:"unordered"`
}

type MockFixture struct {
	Match   MockMatch   `json:"match"`
	Respond MockRespond `json:"respond"`
	// ExpectedCalls, when set, asserts how many times this mock was matched over
	// the case. Use 0 to assert a mock was never called. When nil there is no
	// constraint on this mock's call count.
	ExpectedCalls *int `json:"expectedCalls"`
}

type MockMatch struct {
	Method     string `json:"method"`     // optional; matches any method if empty
	Host       string `json:"host"`       // optional; matches any host if empty
	Path       string `json:"path"`       // optional; exact path match
	PathPrefix string `json:"pathPrefix"` // optional; path prefix match
	// Query, when set, requires each named query parameter to be present on the
	// request with the given value. Extra query parameters are ignored.
	Query map[string]string `json:"query"`
	// Headers, when set, requires each named header to be present with the given
	// value (case-insensitive name). Extra headers are ignored.
	Headers map[string]string `json:"headers"`
	// Body, when set, is matched as a JSON subset against the outbound request
	// body: every field listed must be present with the given value; extra
	// fields in the actual body are ignored.
	Body json.RawMessage `json:"body"`
}

type MockRespond struct {
	Status int             `json:"status"`
	Body   json.RawMessage `json:"body"`
	// BodyString is a raw, non-JSON response body returned verbatim. When set it
	// takes precedence over Body; pair it with ContentType.
	BodyString string `json:"bodyString"`
	// ContentType overrides the response Content-Type (default
	// "application/json").
	ContentType string `json:"contentType"`
}

// DatabaseOutFixture is the expected effect of the request on the database. Each
// table declares which rows were created, updated and/or deleted — the harness
// verifies the operation against the initial state seeded from 1_database.json.
type DatabaseOutFixture struct {
	Tables map[string]TableExpectation `json:"tables"`
}

// TableExpectation describes the mutations expected on a single table.
//
//	Created  rows that exist after the request but did NOT exist before
//	         (their key was absent from the initial state).
//	Updated  rows that existed before (by key) and now hold these values.
//	Deleted  rows that existed before (by key) and are now gone.
//
// Key is the column set used to identify a row across the before/after states
// (defaults to ["id"]). IgnoreColumns are skipped when comparing row values.
//
// Exact, when true, additionally asserts that NO undeclared rows were created
// or deleted in this table: the set of rows whose key is new after the request
// must be exactly the declared Created rows, and the set of rows whose key
// disappeared must be exactly the declared Deleted rows. (Updates do not change
// the key set, so Exact does not constrain them.)
type TableExpectation struct {
	Key           []string `json:"key"`
	Created       []Row    `json:"created"`
	Updated       []Row    `json:"updated"`
	Deleted       []Row    `json:"deleted"`
	IgnoreColumns []string `json:"ignoreColumns"`
	Exact         bool     `json:"exact"`
}

type EventsFixture struct {
	Events []EventExpectation `json:"events"`
}

type EventExpectation struct {
	Schema  string          `json:"schema"`
	Payload json.RawMessage `json:"payload"`
	Ignore  []string        `json:"ignore"`
}

// RedisFixture is the expected Redis state after the request, loaded from the
// mandatory 7_redis.json. An empty fixture ({}) asserts nothing — use it for
// cases that don't touch Redis. A non-empty fixture requires the Boot to expose
// a Redis inspector via BootResult.Redis.
type RedisFixture struct {
	// Values maps a key to its expected string value (compared via GET).
	Values map[string]string `json:"values"`
	// Exists lists keys that must be present (any type/value).
	Exists []string `json:"exists"`
	// Absent lists keys that must not be present.
	Absent []string `json:"absent"`
}

// isEmpty reports whether the fixture declares no Redis expectations, so the
// (still mandatory) file can be a bare {} for cases that don't touch Redis.
func (r *RedisFixture) isEmpty() bool {
	return r == nil || (len(r.Values) == 0 && len(r.Exists) == 0 && len(r.Absent) == 0)
}

// Fixture file names. All are mandatory in every case directory.
const (
	fileDatabase    = "1_database.json"
	fileRequest     = "2_request.json"
	fileMocks       = "3_mocks.json"
	fileDatabaseOut = "4_database_out.json"
	fileJob         = "5_job.json"
	fileResponse    = "6_response.json"
	fileRedis       = "7_redis.json"
)

// requiredFiles is the complete set every case directory must contain.
var requiredFiles = []string{
	fileDatabase, fileRequest, fileMocks, fileDatabaseOut, fileJob, fileResponse, fileRedis,
}

// loadFixture reads every file for a case directory. All seven files are
// mandatory: a directory missing any of them fails before the request runs.
// 7_redis.json may be a bare {} when the case asserts no Redis state.
func loadFixture(tb testing.TB, caseDir string) Fixture {
	tb.Helper()

	requireAllFiles(tb, caseDir)

	var f Fixture

	f.Request.Raw = loadJSON(tb, filepath.Join(caseDir, fileRequest), &f.Request)
	loadJSON(tb, filepath.Join(caseDir, fileMocks), &f.Mocks)
	loadJSON(tb, filepath.Join(caseDir, fileResponse), &f.Response)
	loadJSON(tb, filepath.Join(caseDir, fileJob), &f.Events)

	f.Redis = &RedisFixture{}
	loadJSON(tb, filepath.Join(caseDir, fileRedis), f.Redis)

	// Database fixtures need json.Number so integer columns are not corrupted.
	loadJSONWithNumbers(tb, filepath.Join(caseDir, fileDatabase), &f.Database)
	loadJSONWithNumbers(tb, filepath.Join(caseDir, fileDatabaseOut), &f.DatabaseOut)

	return f
}

// requireAllFiles fails the test if any mandatory file is missing, listing every
// absent file so the case can be completed in one pass.
func requireAllFiles(tb testing.TB, caseDir string) {
	tb.Helper()

	var missing []string
	for _, name := range requiredFiles {
		if !fileExists(filepath.Join(caseDir, name)) {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		tb.Fatalf("e2e: case %q is incomplete, missing required file(s): %v", caseDir, missing)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// loadJSON decodes path into dst and returns the raw bytes read. Decoding is
// strict: an unknown key in the file (e.g. a misspelled field) fails the case
// rather than being silently ignored.
func loadJSON(tb testing.TB, path string, dst any) json.RawMessage {
	tb.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		tb.Fatalf("e2e: read %s: %v", path, err)
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		tb.Fatalf("e2e: parse %s: %v", path, err)
	}

	return data
}

// loadJSONWithNumbers decodes using json.Number so integer column values are not
// converted to float64 (which the SQL layer would render as e.g. "1"). Decoding
// is strict, as in loadJSON.
func loadJSONWithNumbers(tb testing.TB, path string, dst any) {
	tb.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		tb.Fatalf("e2e: read %s: %v", path, err)
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		tb.Fatalf("e2e: parse %s: %v", path, err)
	}
}
