package e2e

import (
	"database/sql"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/onsi/gomega"
)

// assertResponse checks the HTTP status, response headers (subset) and (if
// provided) the response body as a subset: every field present in the expected
// body must match the actual body; extra actual fields are ignored. Paths listed
// in Ignore are dropped from both sides before comparison.
//
// Status is mandatory: a response fixture that omits it fails, so a case can
// never silently assert nothing about an HTTP response.
func assertResponse(g gomega.Gomega, expected ResponseFixture, status int, header http.Header, body []byte) {
	g.Expect(expected.Status).NotTo(gomega.BeZero(),
		"6_response.json must declare a non-zero \"status\" for an HTTP case")
	g.Expect(status).To(gomega.Equal(expected.Status),
		"unexpected HTTP status; body: %s", string(body))

	assertResponseHeaders(g, expected.Headers, header)

	// A raw, non-JSON body is matched by exact string equality.
	if expected.BodyString != "" {
		g.Expect(strings.TrimSpace(string(body))).To(gomega.Equal(strings.TrimSpace(expected.BodyString)),
			"response body did not match expected raw string")

		return
	}

	if len(expected.Body) == 0 {
		return
	}

	wantRoot := parseJSON(g, expected.Body, "expected response body")
	gotRoot := parseJSON(g, body, "actual response body")

	for _, p := range expected.Ignore {
		removePath(wantRoot, p)
		removePath(gotRoot, p)
	}

	mismatches := subsetMatchUnordered(wantRoot, gotRoot, "$", expected.Unordered)
	g.Expect(mismatches).To(gomega.BeEmpty(),
		"response body did not match expectation:\n  %s\nfull body: %s",
		strings.Join(mismatches, "\n  "), string(body))
}

// RedisInspector is the read-only subset of an in-memory Redis the engine needs
// to assert state from 7_redis.json. *miniredis.Miniredis satisfies it, so a
// Boot using the redismem helper can expose its server directly via
// BootResult.Redis.
type RedisInspector interface {
	Get(key string) (string, error)
	Exists(key string) bool
}

// assertRedis verifies the expected Redis state against the Boot-provided
// inspector. It fails the case if a 7_redis.json is present but no inspector was
// exposed, so the assertion can never be silently skipped.
func assertRedis(g gomega.Gomega, inspector RedisInspector, expected *RedisFixture) {
	if expected.isEmpty() {
		return
	}

	g.Expect(inspector).NotTo(gomega.BeNil(),
		"7_redis.json declares Redis expectations but Boot did not expose a Redis inspector via BootResult.Redis")

	for _, key := range sortedHeaderNames(expected.Values) {
		got, err := inspector.Get(key)
		g.Expect(err).NotTo(gomega.HaveOccurred(),
			"redis: GET %q (key missing or not a string)", key)
		g.Expect(got).To(gomega.Equal(expected.Values[key]),
			"redis: unexpected value for key %q", key)
	}

	for _, key := range expected.Exists {
		g.Expect(inspector.Exists(key)).To(gomega.BeTrue(),
			"redis: expected key %q to exist", key)
	}

	for _, key := range expected.Absent {
		g.Expect(inspector.Exists(key)).To(gomega.BeFalse(),
			"redis: expected key %q to be absent", key)
	}
}

// assertResponseHeaders checks that every expected header is present with the
// expected value. Header names are case-insensitive (canonicalised by
// http.Header); response headers not named in want are ignored.
func assertResponseHeaders(g gomega.Gomega, want map[string]string, got http.Header) {
	for _, name := range sortedHeaderNames(want) {
		g.Expect(got.Get(name)).To(gomega.Equal(want[name]),
			"unexpected value for response header %q", name)
	}
}

func sortedHeaderNames(m map[string]string) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}

// assertDatabase verifies the created/updated/deleted rows declared per table
// against the diff between the initial state (initial) and the state after the
// request. Only the columns named in an expected row are compared (subset of
// columns); IgnoreColumns are skipped entirely.
func assertDatabase(tb testing.TB, g gomega.Gomega, db *sql.DB, initial map[string][]Row, out DatabaseOutFixture) {
	tb.Helper()

	for table, exp := range out.Tables {
		keyCols := exp.Key
		if len(keyCols) == 0 {
			keyCols = []string{"id"}
		}

		finalRows := queryRows(tb, db, table)
		initialKeys := initialKeySet(initial[table], keyCols)
		finalKeys := finalKeySet(finalRows, keyCols)

		for i, want := range exp.Created {
			ok, detail := findRowFailure(want, finalRows, exp.IgnoreColumns, keyCols, initialKeys, false)
			g.Expect(ok).To(gomega.BeTrue(),
				"table %q: created row #%d not found as a NEW row (key absent before, present after): %v\n%s",
				table, i, want, detail)
		}

		for i, want := range exp.Updated {
			ok, detail := findRowFailure(want, finalRows, exp.IgnoreColumns, keyCols, initialKeys, true)
			g.Expect(ok).To(gomega.BeTrue(),
				"table %q: updated row #%d not found as an EXISTING row holding these values: %v\n%s",
				table, i, want, detail)
		}

		for i, want := range exp.Deleted {
			k := jsonRowKey(want, keyCols)
			g.Expect(initialKeys).To(gomega.HaveKey(k),
				"table %q: deleted row #%d was not present in the initial state: %v", table, i, want)
			g.Expect(finalKeys).NotTo(gomega.HaveKey(k),
				"table %q: deleted row #%d is still present after the request: %v", table, i, want)
		}

		if exp.Exact {
			assertExactKeyChanges(g, table, exp, keyCols, initialKeys, finalKeys)
		}
	}
}

// assertExactKeyChanges fails if any row was created or deleted in the table
// beyond those declared. A "created" row is one whose key is present after the
// request but was absent before; a "deleted" row is the reverse. Only key-set
// changes are constrained — updates leave the key set unchanged.
func assertExactKeyChanges(g gomega.Gomega, table string, exp TableExpectation, keyCols []string, initialKeys, finalKeys map[string]struct{}) {
	declaredCreated := initialKeySet(exp.Created, keyCols)
	declaredDeleted := initialKeySet(exp.Deleted, keyCols)

	var unexpectedCreated []string
	for k := range finalKeys {
		if _, preexisting := initialKeys[k]; preexisting {
			continue
		}
		if _, declared := declaredCreated[k]; !declared {
			unexpectedCreated = append(unexpectedCreated, prettyKey(k))
		}
	}

	var unexpectedDeleted []string
	for k := range initialKeys {
		if _, stillPresent := finalKeys[k]; stillPresent {
			continue
		}
		if _, declared := declaredDeleted[k]; !declared {
			unexpectedDeleted = append(unexpectedDeleted, prettyKey(k))
		}
	}

	sort.Strings(unexpectedCreated)
	sort.Strings(unexpectedDeleted)

	g.Expect(unexpectedCreated).To(gomega.BeEmpty(),
		"table %q (exact): rows were created but not declared (key columns %v): %v",
		table, keyCols, unexpectedCreated)
	g.Expect(unexpectedDeleted).To(gomega.BeEmpty(),
		"table %q (exact): rows were deleted but not declared (key columns %v): %v",
		table, keyCols, unexpectedDeleted)
}

// prettyKey renders a composite row key (NUL-joined internally) for display.
func prettyKey(k string) string {
	return strings.ReplaceAll(k, "\x00", ",")
}

// assertEvents verifies that exactly the expected events were published.
// Order-independent: each expected event must match a recorded event by schema
// and payload subset, and the counts must be equal.
func assertEvents(g gomega.Gomega, recorded []RecordedEvent, expected EventsFixture) {
	g.Expect(recorded).To(gomega.HaveLen(len(expected.Events)),
		"unexpected number of published events; recorded: %v", schemasOf(recorded))

	used := make([]bool, len(recorded))
	for i, want := range expected.Events {
		matched := false
		for j, got := range recorded {
			if used[j] || got.Schema != want.Schema {
				continue
			}
			if eventPayloadMatches(g, want, got) {
				used[j] = true
				matched = true

				break
			}
		}
		g.Expect(matched).To(gomega.BeTrue(),
			"no published event matched expected event #%d (schema %q)", i, want.Schema)
	}
}

func eventPayloadMatches(g gomega.Gomega, want EventExpectation, got RecordedEvent) bool {
	if len(want.Payload) == 0 {
		return true
	}

	wantRoot := parseJSON(g, want.Payload, "expected event payload")
	gotRoot := parseJSON(g, got.Payload, "actual event payload")
	for _, p := range want.Ignore {
		removePath(wantRoot, p)
		removePath(gotRoot, p)
	}

	return len(subsetMatch(wantRoot, gotRoot, "$")) == 0
}

func schemasOf(events []RecordedEvent) []string {
	out := make([]string, len(events))
	for i, e := range events {
		out[i] = e.Schema
	}

	return out
}
