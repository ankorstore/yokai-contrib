# fxe2e

A declarative, file-driven **end-to-end test harness** for [Yokai](https://github.com/ankorstore/yokai)-based services.

Each test case is a directory of JSON files — no Go code. The harness boots the
**real** application (real DB, real domain services, real workers) against an
in-memory MySQL and only stubs the outside world:

- **outbound HTTP** (including Elasticsearch) is answered from the case's mock list,
- **Redis** is a real in-process in-memory server (miniredis),
- **published events** are recorded for assertion.



---

## The fixture format

Cases live under `testdata/<route>/<method>/<name>/`. Every case directory must
contain **all seven** files (a missing file fails the case before it runs):

| File | Purpose |
|------|---------|
| `1_database.json` | initial DB rows, keyed by table |
| `2_request.json` | the HTTP request (or a named `job`) + auth |
| `3_mocks.json` | outbound HTTP mocks (match → canned response) |
| `4_database_out.json` | expected created / updated / deleted rows |
| `5_job.json` | expected published events |
| `6_response.json` | expected response (status + body subset + ignore paths) |
| `7_redis.json` | expected Redis state — a bare `{}` when the case doesn't touch Redis |

Assertions use subset matching: an expected file only needs to describe the
fields it cares about, and volatile values (timestamps, generated IDs) can be
dropped with an `ignore` list. Add a case by dropping a new folder under
`testdata/`.

Fixture files are decoded **strictly**: an unknown or misspelled key (e.g.
`"mehtod"`) fails the case rather than being silently ignored.

### A complete example

`testdata/checkout/post/create_order/` — a retailer checks out their cart. It
seeds the cart, its line items and the referenced products, sends an
authenticated `POST`, mocks the outbound call to the payment provider, and then
asserts the response, the created **and** updated DB rows, and the published
event. Every one of the six files carries weight here.

**`1_database.json`** — the open cart, its line items, and the products in stock:

```json
{
  "products": [
    { "id": 1, "uuid": "a1111111-1111-4111-8111-111111111111", "brand_uuid": "c90a9d43-e47a-6062-8ed6-0000214bf5f5", "name": "Lavender Soap",  "price_cents": 450,  "stock": 20 },
    { "id": 2, "uuid": "b2222222-2222-4222-8222-222222222222", "brand_uuid": "c90a9d43-e47a-6062-8ed6-0000214bf5f5", "name": "Beeswax Candle", "price_cents": 1200, "stock": 5  }
  ],
  "carts": [
    { "id": 1, "uuid": "c1a2b3c4-d5e6-4f70-8a1b-2c3d4e5f6a7b", "retailer_uuid": "5e9d2f10-3b4a-4c6d-8e1f-9a0b1c2d3e4f", "status": "open" }
  ],
  "cart_items": [
    { "id": 1, "cart_id": 1, "product_id": 1, "quantity": 3 },
    { "id": 2, "cart_id": 1, "product_id": 2, "quantity": 2 }
  ]
}
```

**`2_request.json`** — an authenticated retailer request (the body is sent
verbatim; `auth` is decoded by the `auth.Func` you registered for `"retailer"`):

```json
{
  "method": "POST",
  "path": "/api/checkout/v1/orders",
  "headers": { "Content-Type": "application/vnd.api+json" },
  "auth": { "type": "retailer", "retailerUuid": "5e9d2f10-3b4a-4c6d-8e1f-9a0b1c2d3e4f" },
  "body": {
    "data": {
      "type": "orders",
      "id": "0d1e2f30-4a5b-4c6d-8e9f-0a1b2c3d4e5f",
      "attributes": {
        "cartUuid": "c1a2b3c4-d5e6-4f70-8a1b-2c3d4e5f6a7b",
        "paymentMethod": "card"
      }
    }
  }
}
```

**`3_mocks.json`** — checkout calls the payment provider to open a payment intent
for the cart total (3 × 450 + 2 × 1200 = **3750** cents); answer it from a canned
response (an unmatched outbound call fails the test):

```json
[
  {
    "match":   { "method": "POST", "pathPrefix": "/v1/payment_intents" },
    "respond": {
      "status": 200,
      "body": {
        "id": "pi_3NxAbc123",
        "status": "requires_capture",
        "amount": 3750,
        "currency": "eur"
      }
    }
  }
]
```

**`4_database_out.json`** — assert the order was *created*, the cart *updated*
(`open` → `checked_out`, payment intent recorded), and each product's stock
*decremented*. `created` and `updated` rows are both matched by `key`; only the
listed columns are compared, volatile ones are ignored.

By default the check is a **subset**: rows you don't declare are ignored, so an
unexpected insert or delete (e.g. a cascade you didn't anticipate) goes
unnoticed. Set `"exact": true` on a table to additionally assert that *no*
undeclared rows were created or deleted there — the set of newly-keyed rows must
be exactly the declared `created`, and the set of vanished keys exactly the
declared `deleted`. (Updates leave the key set unchanged, so `exact` does not
constrain them.) A failing created/updated assertion now names the offending
columns of the row that shares the expected key, rather than dumping every row.

```json
{
  "tables": {
    "orders": {
      "key": ["uuid"],
      "created": [
        {
          "uuid": "0d1e2f30-4a5b-4c6d-8e9f-0a1b2c3d4e5f",
          "retailer_uuid": "5e9d2f10-3b4a-4c6d-8e1f-9a0b1c2d3e4f",
          "cart_uuid": "c1a2b3c4-d5e6-4f70-8a1b-2c3d4e5f6a7b",
          "status": "pending_payment",
          "total_cents": 3750,
          "payment_intent_id": "pi_3NxAbc123"
        }
      ],
      "ignoreColumns": ["id", "created_at", "updated_at"]
    },
    "carts": {
      "key": ["uuid"],
      "updated": [
        { "uuid": "c1a2b3c4-d5e6-4f70-8a1b-2c3d4e5f6a7b", "status": "checked_out" }
      ],
      "ignoreColumns": ["updated_at"]
    },
    "products": {
      "key": ["uuid"],
      "updated": [
        { "uuid": "a1111111-1111-4111-8111-111111111111", "stock": 17 },
        { "uuid": "b2222222-2222-4222-8222-222222222222", "stock": 3  }
      ],
      "ignoreColumns": ["updated_at"]
    }
  }
}
```

**`5_job.json`** — assert the checkout published exactly one `order.placed.v1`
event. Matching is **count-exact** (the number of recorded events must equal the
number listed) and **order-independent** (each expected event is matched to a
recorded one by `schema` + payload subset). Each entry has three keys:

- `schema` — the event's schema namespace (what your publisher reports via
  `SchemaNamespace()`); must match exactly.
- `payload` — a **subset** of the published payload. Only the fields you list are
  compared; extra fields in the actual event are allowed. The payload is the JSON
  rendering of your event struct, so the keys are whatever that struct marshals to.
- `ignore` — dotted paths dropped from both sides before comparison; use it for
  generated IDs, timestamps, and other volatile values.

```json
{
  "events": [
    {
      "schema": "order.placed.v1",
      "payload": {
        "OrderUUID": "0d1e2f30-4a5b-4c6d-8e9f-0a1b2c3d4e5f",
        "RetailerUUID": "5e9d2f10-3b4a-4c6d-8e1f-9a0b1c2d3e4f",
        "TotalCents": 3750,
        "Currency": "EUR"
      },
      "ignore": ["EventID", "PlacedAt"]
    }
  ]
}
```

(For a flow that publishes nothing, this file is just `{ "events": [] }` — an
empty list still asserts that *zero* events were recorded.)

**`6_response.json`** — `201` with a body subset; volatile/derived fields are
dropped via `ignore`. `status` is **mandatory** for an HTTP case — omitting it
fails the test rather than silently asserting nothing. `headers` is optional and
matched as a subset: every header you list must be present with that value (names
are case-insensitive); response headers you don't list are ignored.

```json
{
  "status": 201,
  "headers": { "Content-Type": "application/vnd.api+json" },
  "body": {
    "data": {
      "type": "orders",
      "id": "0d1e2f30-4a5b-4c6d-8e9f-0a1b2c3d4e5f",
      "attributes": {
        "status": "pending_payment",
        "totalCents": 3750,
        "paymentStatus": "requires_capture"
      }
    }
  },
  "ignore": ["data.attributes.createdAt", "data.attributes.updatedAt"]
}
```

That's the entire test. The harness seeds the rows, fires the request (serving the
payment-provider call from `3_mocks.json`), then checks the response, the database
mutations, and the published events — all declaratively.

Each outbound request is matched against `3_mocks.json` in order; an unmatched
call fails the test. The same `match`/`respond` shape covers Elasticsearch — match
the index path and answer with a canned `_search` body:

```json
[
  {
    "match":   { "pathPrefix": "/my_index/_search" },
    "respond": {
      "status": 200,
      "body": { "hits": { "total": { "value": 1 }, "hits": [ { "_source": { "uuid": "…" } } ] } }
    }
  }
]
```

---

## Matching beyond the basics

The defaults above cover most cases; these knobs handle the rest.

**Match a mock on more than the path.** Besides `method`/`host`/`path`/`pathPrefix`,
a `match` can also require query parameters, headers, and a JSON **body subset** —
so two calls to the same path are told apart by what they send:

```json
{
  "match": {
    "method": "POST",
    "pathPrefix": "/v1/payment_intents",
    "query":   { "expand": "charges" },
    "headers": { "Idempotency-Key": "abc" },
    "body":    { "amount": 3750, "currency": "eur" }
  },
  "respond": { "status": 200, "body": { "id": "pi_123" } },
  "expectedCalls": 1
}
```

**Assert how often a mock is hit.** `expectedCalls` (above) asserts an exact
count; use `0` to assert a mock was *never* called. Omit it for no constraint.

**Non-JSON bodies.** Use `bodyString` (+ `contentType`) for raw payloads on a
mock response, the request, or the expected response (the latter is matched by
exact string equality):

```json
"respond": { "status": 200, "contentType": "text/csv", "bodyString": "id,name\n1,foo" }
```

**Order-independent arrays.** List the dotted paths whose arrays should match
regardless of order in `6_response.json` — lengths must still be equal, each
expected element must match a distinct actual one:

```json
{ "status": 200, "body": { "data": [ … ] }, "unordered": ["data"] }
```

**Redis state.** `7_redis.json` is mandatory — a bare `{}` for cases that don't
touch Redis. When it declares expectations, the harness asserts the final Redis
state, provided the `Boot` exposes the in-memory server via `BootResult.Redis`
(a `*miniredis.Miniredis` from `redismem.Server` satisfies it):

```json
{
  "values": { "lease:brand:xyz": "locked" },
  "exists": ["oauth:state:abc"],
  "absent": ["stale:key"]
}
```

**Deterministic clock.** A request may declare `"now": "2030-01-01T00:00:00Z"`;
`Boot` reads it via `fix.FixedTime()` and decorates the app's clock so
time-dependent logic runs reproducibly.

---

## Authentication

A case authenticates by declaring a principal under `auth` in `2_request.json`.
The harness knows the Ankorstore principal types out of the box and applies the
matching token automatically — **no wiring in your `Boot` or test entrypoint**.
An absent `auth` (or `"type": "none"`) sends the request anonymously.

```json
"auth": { "type": "brand", "accountUuid": "c90a9d43-e47a-6062-8ed6-0000214bf5f5" }
```

Every field is optional — `go-modules` fills sensible defaults (e.g. a default
account UUID and email), so `{ "type": "brand" }` already yields a valid token.

| `type` | Fields |
|--------|--------|
| `guest` | `clientId` |
| `machine` | `clientId` |
| `admin` | `clientId`, `entityUuid`, `roles[]`, `permissions[]` |
| `brand` | `clientId`, `entityUuid`, `accountUuid`, `accountEmail` |
| `retailer` | `clientId`, `entityUuid`, `accountUuid`, `accountEmail` |
| `impersonation` | `accountType` (`brand`\|`retailer`), `accountUuid`, `accountEmail`, `entityUuid`, `clientId`, `impersonatorUuid`, `impersonatorRoles[]`, `impersonatorPermissions[]` |

An admin impersonating a brand:

```json
"auth": {
  "type": "impersonation",
  "accountType": "brand",
  "accountUuid": "c90a9d43-e47a-6062-8ed6-0000214bf5f5",
  "impersonatorRoles": ["customer-care"]
}
```

An optional `origin` (`internal` or `external`) sets the request's SPIFFE origin
on top of any principal — including anonymous — for routes that gate on it:

```json
"auth": { "type": "brand", "accountUuid": "…", "origin": "internal" }
```

This is wired through `e2e.Runner.Auth`, which defaults to `auth.Ankorstore()`.
Set that field only if you need to override the default dispatch.

---

## Integrating it into a Yokai application

### 1. Add the dependency

```bash
go get github.com/ankorstore/yokai-contrib/fxe2e
```



### 2. Write a `Boot` function

This is the only real wiring. `Boot` boots your app in test mode and returns a
`BootResult`. It is where you connect the harness's seams to your app's Fx graph:
route outbound HTTP through the supplied `transport`, record published events,
point Redis/Elasticsearch at the in-memory equivalents, pin the clock, and
expose your job runner (produced *here* because it captures services
resolved from the freshly-booted container). The framework helpers
(`pubsubmem`, `redismem`, `esmock`, `clockmem`) keep this to a few lines.

```go
package e2e

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/myorg/myapp/internal" // your app's RunE2ETest bootstrapper
	"github.com/ankorstore/yokai-contrib/fxe2e/clockmem"
	"github.com/ankorstore/yokai-contrib/fxe2e/e2e"
	"github.com/ankorstore/yokai-contrib/fxe2e/esmock"
	"github.com/ankorstore/yokai-contrib/fxe2e/pubsubmem"
	"github.com/ankorstore/yokai-contrib/fxe2e/redismem"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func Boot(tb testing.TB, fix *e2e.Fixture, transport http.RoundTripper) e2e.BootResult {
	tb.Helper()

	// Capture published events (decorates the shared gcppubsub.Publisher).
	events, recordEvents := pubsubmem.Recorder()

	var (
		httpServer *echo.Echo
		db         *sql.DB
		jobs       *mycron.Registry // your job registry: Run(ctx, name, args...) error
	)

	internal.RunE2ETest(
		tb,
		// 1. Route all outbound HTTP through the case's mocks.
		fx.Decorate(func(c *http.Client) *http.Client {
			return &http.Client{Transport: transport, Timeout: c.Timeout}
		}),
		// 2. (optional) Record published events for 5_job.json assertions.
		recordEvents,
		// 3. (optional) Real in-memory Redis — your Redis-backed code runs unchanged.
		redismem.FxOption(tb),
		// 4. (optional) Elasticsearch over the same mock transport.
		esmock.FxOption(transport),
		// 5. (optional) Pin the clock when a case fixes "now"; no-op otherwise.
		clockmem.FxOption(fix),

		fx.Populate(&httpServer, &db, &jobs),
	)

	return e2e.BootResult{
		Server: httpServer,        // anything implementing http.Handler
		DB:     db,
		Events: events.Recorded,   // func() []e2e.RecordedEvent
		Jobs:   jobs,              // (optional) runs a case's "job" by name
	}
}
```

> **Note on `RunE2ETest`.** Most Yokai apps already have a `RunTest` helper that
> boots the Fx app against an in-memory MySQL and runs migrations. For E2E you
> want the same thing **without** loading Go seed data (each case owns its
> initial state via `1_database.json`). Add a `RunE2ETest` variant that runs
> migrations but skips seeds.

### 3. Write the test entrypoint

```go
package e2e_test

import (
	"testing"

	"github.com/myorg/myapp/internal/e2e"
	libe2e "github.com/ankorstore/yokai-contrib/fxe2e/e2e"
)

func TestE2E(t *testing.T) {
	// No Auth wiring needed: the runner authenticates each case from the "type"
	// in its 2_request.json (see "Authentication" below).
	runner := &libe2e.Runner{Boot: e2e.Boot}

	cases := libe2e.DiscoverCases(t, "testdata")
	if len(cases) == 0 {
		t.Fatal("e2e: no test cases discovered under testdata/")
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			runner.Run(t, c.Dir)
		})
	}
}
```

### 4. Add your first case

Create `internal/e2e/testdata/<route>/<method>/<name>/` with the six JSON files.
That's it — `TestE2E` discovers it automatically.

---

## Running the tests

The whole harness is plain `go test`, exposed as a single `TestE2E` with one
subtest per case (named by its path under `testdata/`). A `Makefile` target keeps
it ergonomic:

```makefile
TESTARGS ?=

# Run only end-to-end tests.            make test-e2e
# Run one case / route / method:        make test-e2e CASE=sync/post/start_fetch
test-e2e:
	@c='$(CASE)'; c="$${c#./}"; c="$${c#internal/e2e/testdata/}"; c="$${c#testdata/}"; c="$${c%/}"; \
	go test -failfast -race -v $${c:+-run "TestE2E/$$c"} $(TESTARGS) ./internal/e2e/...
```

### Only the E2E suite

```bash
make test-e2e
# equivalent to:
go test -v ./internal/e2e/...
```

### A specific case (or a whole route/method)

`CASE` is the subtest path — everything after `testdata/`. A leading
`internal/e2e/testdata/` or `testdata/` prefix and trailing slash are stripped,
so a shell-completed path works too.

```bash
make test-e2e CASE=sync/post/start_fetch   # one case
make test-e2e CASE=sync/post               # every case under a method
make test-e2e CASE=sync                     # every case for a route
```

Under the hood this is just Go's `-run` regex against the subtest tree:

```bash
go test -run 'TestE2E/sync/post/start_fetch' ./internal/e2e/...
go test -run 'TestE2E/sync'                   ./internal/e2e/...   # prefix match
```

Extra flags go through `TESTARGS`, e.g. `make test-e2e TESTARGS=-count=1`.

> If your repo's default `make test` keeps `PKG=./...`, don't pass the e2e path
> through it (the suite still runs in full). Use `test-e2e`.

---

## What the library provides

| Package | What it does | App-specific? |
|---------|--------------|---------------|
| `e2e` | the engine: case discovery, fixture loading, DB seeding, HTTP mock transport, DB diffing, JSON subset assertions, `Runner`/`BootResult` | generic |
| `pubsubmem` | records published events for `5_job.json` (decorates `gcppubsub.Publisher`) | for apps on `gcppubsub` |
| `recorder` | concurrency-safe event sink (returned by `pubsubmem.Recorder`) | generic |
| `await` | `Until`/`Count`/`NoRows` poll helpers; powers a case's `"await"` query | generic |
| `redismem` | real in-process in-memory Redis (miniredis), repoints `*redis.Client` | for apps on `fxredis` |
| `esmock` | routes `*elasticsearch.Client` through the mock transport | for apps on `fxelasticsearch` |
| `clockmem` | pins `clockwork.Clock` to a case's `"now"` | for apps on `fxclock` |
| `auth` | `Ankorstore()` dispatch over all principal types (guest/machine/admin/brand/retailer/impersonation), applied from each case's `auth.type` | Ankorstore-specific; used by `Runner` automatically |

`pubsubmem`, `redismem`, `esmock` and `clockmem` are **opt-in** fx-option
helpers — import only what your app uses. `auth` is applied automatically by the
`Runner` (see [Authentication](#authentication)).

### Running a background job instead of a request

When a case sets `"job": "<name>"` in `2_request.json`, the harness runs that
named background job instead of sending an HTTP request (the HTTP fields are
ignored and `6_response.json` is not asserted). Expose your app's job registry —
anything with `Run(ctx, name, args...) error`, e.g. a yokai/cron registry — as
`BootResult.Jobs`, and every registered job is runnable by name with no per-job
wiring. Optional `"args": ["…"]` are passed through to the job.

```json
{ "job": "run-sync", "await": { "query": "SELECT COUNT(*) FROM sync_runs WHERE status = 'running'" } }
```

### Waiting for background work

When a journey finishes its work in a background goroutine (a job, or an HTTP
handler that returns early and keeps working), the action returns before the
effect lands. Set `"await"` to a query that observes the effect, and the harness
polls the database until it settles before asserting:

```json
"await": { "query": "SELECT COUNT(*) FROM sync_runs WHERE status = 'running'" }
```

The query is a `COUNT(*)`; the harness waits until it returns `equals` (default
`0`), or fails after `timeout` (default `10s`, polled every `interval`, default
`25ms`):

```json
"await": { "query": "SELECT COUNT(*) FROM outbox WHERE sent = 0", "equals": 0, "timeout": "5s", "interval": "20ms" }
```

This is deliberately effect-based rather than goroutine-based: you can't reliably
join a fire-and-forget `go f()` from outside (polling the process-wide goroutine
count is defeated by pools, tickers and per-case boots), so the harness waits on
the *observable result* instead. Crucially, **this needs no change to your
application code** — the wait condition is test data, so the harness drops into
any project unmodified.

---

## How heavy is it?

Same infrastructure cost as ordinary Yokai handler/contract tests: each case
boots the app against the same in-memory MySQL via `RunTestApp` (one fresh
migrated DB per case), with Redis as in-memory miniredis and all outbound HTTP
(incl. Elasticsearch) served in-process. **No containers, no network.**

It exercises *more real code* — the full app boot and the end-to-end journey
through every layer — which is the point, and where the extra wall-clock comes
from. It's the right tool for complex multi-layer flows (especially ones with a
background job in the journey); contract tests remain the right tool for breadth
across the API surface.
