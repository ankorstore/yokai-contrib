# fxe2e

A declarative, file-driven **end-to-end test harness** for [Yokai](https://github.com/ankorstore/yokai) HTTP services.

Each test case is a directory of JSON files — **no Go code**. The harness boots the
**real** application (real DB, domain services, workers) and only stubs the
outside world:

- **outbound HTTP** (including Elasticsearch) is answered from the case's mock list,
- **MySQL** is an in-process [go-mysql-server](https://github.com/dolthub/go-mysql-server), freshly migrated per case,
- **Redis** is a real in-process server ([miniredis](https://github.com/alicebob/miniredis)),
- **published events** are recorded for assertion.

**No containers, no network** — it's plain `go test`.

---

## Case layout

Cases live under `testdata/<route>/<method>/<name>/`. Every case directory must
contain **all seven** files — a missing one fails the case before it runs:

| File | Purpose |
|------|---------|
| `1_database.json` | initial DB rows, keyed by table |
| `2_request.json` | the HTTP request (or a named `job`) + auth |
| `3_mocks.json` | outbound HTTP mocks (match → canned response) |
| `4_database_out.json` | expected created / updated / deleted rows |
| `5_job.json` | expected published events |
| `6_response.json` | expected response (status + headers + body) |
| `7_redis.json` | expected Redis state (a bare `{}` when the case doesn't touch Redis) |

Two rules apply throughout:

- **Subset matching** — an expected file only describes the fields it cares
  about; extra fields in the actual value are ignored. Volatile values
  (timestamps, generated IDs) are dropped via an `ignore` list.
- **Strict decoding** — an unknown or misspelled key (e.g. `"mehtod"`) fails the
  case rather than being silently ignored.

Add a case by dropping a new folder under `testdata/` — `TestE2E` discovers it.

---

## A complete example

`testdata/checkout/post/create_order/` — a retailer checks out a cart: it seeds
the cart, line items and products, sends an authenticated `POST`, mocks the
payment provider, then asserts the response, the DB mutations and the published
event.

**`1_database.json`** — initial rows, keyed by table:

```json
{
  "products": [
    { "id": 1, "uuid": "a1111111-1111-4111-8111-111111111111", "name": "Lavender Soap",  "price_cents": 450,  "stock": 20 },
    { "id": 2, "uuid": "b2222222-2222-4222-8222-222222222222", "name": "Beeswax Candle", "price_cents": 1200, "stock": 5  }
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

**`2_request.json`** — the request; `auth` is applied automatically from its
`type` (see [Authentication](#authentication)), `body` is sent verbatim:

```json
{
  "method": "POST",
  "path": "/api/checkout/v1/orders",
  "headers": { "Content-Type": "application/vnd.api+json" },
  "auth": { "type": "retailer", "accountUuid": "5e9d2f10-3b4a-4c6d-8e1f-9a0b1c2d3e4f" },
  "body": {
    "data": { "type": "orders", "attributes": { "cartUuid": "c1a2b3c4-d5e6-4f70-8a1b-2c3d4e5f6a7b", "paymentMethod": "card" } }
  }
}
```

**`3_mocks.json`** — answer the payment-provider call (an unmatched outbound call
fails the test):

```json
[
  {
    "match":   { "method": "POST", "pathPrefix": "/v1/payment_intents" },
    "respond": { "status": 200, "body": { "id": "pi_3NxAbc123", "status": "requires_capture", "amount": 3750 } }
  }
]
```

**`4_database_out.json`** — the order *created*, the cart *updated*, each
product's stock *decremented*:

```json
{
  "tables": {
    "orders": {
      "key": ["uuid"],
      "created": [ { "retailer_uuid": "5e9d2f10-3b4a-4c6d-8e1f-9a0b1c2d3e4f", "status": "pending_payment", "total_cents": 3750 } ],
      "ignoreColumns": ["id", "uuid", "created_at", "updated_at"]
    },
    "carts": {
      "key": ["uuid"],
      "updated": [ { "uuid": "c1a2b3c4-d5e6-4f70-8a1b-2c3d4e5f6a7b", "status": "checked_out" } ],
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

**`5_job.json`** — exactly one `order.placed.v1` event was published:

```json
{
  "events": [
    { "schema": "order.placed.v1", "payload": { "RetailerUUID": "5e9d2f10-3b4a-4c6d-8e1f-9a0b1c2d3e4f", "TotalCents": 3750 } }
  ]
}
```

**`6_response.json`** — `201` with a body subset:

```json
{
  "status": 201,
  "body": { "data": { "type": "orders", "attributes": { "status": "pending_payment", "totalCents": 3750 } } },
  "ignore": ["data.attributes.createdAt"]
}
```

**`7_redis.json`** — this flow touches no Redis:

```json
{}
```

That's the whole test. The harness seeds the rows, fires the request (serving the
mock), then checks the response, DB mutations, events and Redis — declaratively.

---

## Fixture reference

### `2_request.json`

| Field | Meaning |
|-------|---------|
| `method`, `path`, `headers` | the HTTP request (method defaults to `GET`) |
| `body` | JSON request body, sent verbatim |
| `bodyString` | raw (non-JSON) body; pair with a `Content-Type` header |
| `auth` | principal to authenticate as — see [Authentication](#authentication) |
| `job`, `args` | run a named background job instead of an HTTP request — see [Background work](#background-jobs-and-async-work) |
| `await` | wait for a background effect before asserting — see [Background work](#background-jobs-and-async-work) |
| `now` | pin the clock (RFC 3339) — see [Deterministic time](#deterministic-time) |

### `3_mocks.json` — outbound HTTP mocks

A list of `{ match, respond }`, tried in order; the first match answers. **Any
outbound call that matches no mock fails the test.**

- **`match`** narrows on `method`, `host`, `path`, `pathPrefix`, `query` (each
  param must be present with that value), `headers` (case-insensitive), and
  `body` (JSON subset of the request body) — so two calls to the same path are
  told apart by what they send.
- **`respond`** sets `status` (default `200`), `body` (JSON) or `bodyString` +
  `contentType` (raw), default `application/json`.
- **`expectedCalls`** (optional) asserts how many times the mock was matched; use
  `0` to assert it was never called. Omit for no constraint.

The same shape covers Elasticsearch — match the index path, answer with a canned
`_search` body:

```json
[ { "match": { "pathPrefix": "/my_index/_search" }, "respond": { "body": { "hits": { "total": { "value": 1 }, "hits": [] } } } } ]
```

### `4_database_out.json` — DB mutations

`tables.<name>` declares, per table:

- **`key`** — columns identifying a row across before/after (default `["id"]`).
- **`created`** / **`updated`** / **`deleted`** — rows whose key is new after /
  existed before and now holds these values / vanished after. Only listed columns
  are compared.
- **`ignoreColumns`** — columns skipped in comparison (volatile values).
- **`exact`** (default `false`) — also assert that *no* undeclared rows were
  created or deleted in this table. By default the check is a subset, so an
  unexpected insert/delete goes unnoticed unless `exact` is set. (Updates don't
  change the key set, so `exact` doesn't constrain them.)

### `5_job.json` — published events

`events` is matched **count-exact** (recorded count must equal the listed count)
and **order-independent** (each expected event matches a recorded one by `schema`
+ payload subset). Per entry:

- **`schema`** — the event's schema namespace (`SchemaNamespace()`); exact match.
- **`payload`** — subset of the published payload. The payload is the JSON
  rendering of your event struct, so keys are whatever it marshals to.
- **`ignore`** — dotted paths dropped from both sides (generated IDs, timestamps).

`{ "events": [] }` asserts that *zero* events were published.

### `6_response.json` — HTTP response

- **`status`** — **mandatory** for an HTTP case (omitting it fails rather than
  silently asserting nothing).
- **`headers`** — subset; each listed header must be present with that value
  (names case-insensitive).
- **`body`** — JSON subset, with **`ignore`** (dotted paths) and **`unordered`**
  (dotted paths whose arrays match order-independently; lengths must still match).
- **`bodyString`** — assert a raw (non-JSON) body by exact string equality.

### `7_redis.json` — Redis state

`{}` when the case asserts nothing. Otherwise:

```json
{ "values": { "lease:brand:xyz": "locked" }, "exists": ["oauth:state:abc"], "absent": ["stale:key"] }
```

- **`values`** — key → expected string value (via `GET`).
- **`exists`** / **`absent`** — keys that must / must not be present.

A non-empty fixture requires the `Boot` to expose the server via
`BootResult.Redis` (a `*miniredis.Miniredis` from `redismem.Server` satisfies it).

---

## Authentication

A case authenticates by declaring a principal under `auth`. The harness knows the
Ankorstore principal types out of the box and applies the matching token
automatically — **no wiring in your `Boot` or test entrypoint**. An absent `auth`
(or `"type": "none"`) sends the request anonymously.

```json
"auth": { "type": "brand", "accountUuid": "c90a9d43-e47a-6062-8ed6-0000214bf5f5" }
```

Every field is optional — `go-modules` fills sensible defaults — so
`{ "type": "brand" }` already yields a valid token.

| `type` | Fields |
|--------|--------|
| `guest` | `clientId` |
| `machine` | `clientId` |
| `admin` | `clientId`, `entityUuid`, `roles[]`, `permissions[]` |
| `brand` | `clientId`, `entityUuid`, `accountUuid`, `accountEmail` |
| `retailer` | `clientId`, `entityUuid`, `accountUuid`, `accountEmail` |
| `impersonation` | `accountType` (`brand`\|`retailer`), `accountUuid`, `accountEmail`, `entityUuid`, `clientId`, `impersonatorUuid`, `impersonatorRoles[]`, `impersonatorPermissions[]` |

An optional `origin` (`internal` or `external`) sets the request's SPIFFE origin
on top of any principal — including anonymous — for routes that gate on it:

```json
"auth": { "type": "brand", "accountUuid": "…", "origin": "internal" }
```

This is driven by `e2e.Runner.Auth`, which defaults to `auth.Ankorstore()`; set
that field only to override the default dispatch.

---

## Background jobs and async work

**Run a job instead of a request.** Set `job` to a registered job name and the
harness runs it instead of an HTTP request (the HTTP fields are ignored and
`6_response.json` is not asserted). Optional `args` are passed through. Expose
your app's job registry — anything with `Run(ctx, name, args...) error`, e.g. a
yokai/cron registry — as `BootResult.Jobs`, and every registered job is runnable
by name with no per-job wiring.

**Wait for background effects.** When a journey finishes work in a background
goroutine, the action returns before the effect lands. Set `await` to a `COUNT(*)`
query that observes the effect; the harness polls the DB until it returns `equals`
(default `0`), or fails after `timeout` (default `10s`, polled every `interval`,
default `25ms`):

```json
{ "job": "run-sync", "await": { "query": "SELECT COUNT(*) FROM sync_runs WHERE status = 'running'" } }
```

This is deliberately effect-based: you can't reliably join a fire-and-forget
`go f()` from outside (polling the goroutine count is defeated by pools, tickers
and per-case boots), so the harness waits on the observable result. **It needs no
change to your application code** — the wait condition is test data.

---

## Deterministic time

A case may pin the clock with `"now": "2030-01-01T00:00:00Z"` (RFC 3339). Wire
`clockmem.FxOption(fix)` into your `Boot` (see below) and the app's
`clockwork.Clock` is replaced by a fake clock at that instant, so timestamps and
expiries are reproducible.

---

## Integrating into a Yokai application

### 1. Add the dependency

```bash
go get github.com/ankorstore/yokai-contrib/fxe2e
```

### 2. Write a `Boot` function

The only real wiring: `Boot` boots your app in test mode and returns a
`BootResult`, connecting the harness's seams to your Fx graph. The opt-in helpers
(`pubsubmem`, `redismem`, `esmock`, `clockmem`) keep it to a few lines.

```go
package e2e

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/myorg/myapp/internal" // your app's RunE2ETest bootstrapper
	"github.com/ankorstore/yokai-contrib/fxe2e/clockmem"
	libe2e "github.com/ankorstore/yokai-contrib/fxe2e/e2e"
	"github.com/ankorstore/yokai-contrib/fxe2e/esmock"
	"github.com/ankorstore/yokai-contrib/fxe2e/pubsubmem"
	"github.com/ankorstore/yokai-contrib/fxe2e/redismem"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func Boot(tb testing.TB, fix *libe2e.Fixture, transport http.RoundTripper) libe2e.BootResult {
	tb.Helper()

	events, recordEvents := pubsubmem.Recorder() // capture published events

	var (
		httpServer *echo.Echo
		db         *sql.DB
		jobs       *mycron.Registry // your job registry: Run(ctx, name, args...) error
	)

	internal.RunE2ETest(
		tb,
		// Route outbound HTTP through the case's mocks.
		fx.Decorate(func(c *http.Client) *http.Client {
			return &http.Client{Transport: transport, Timeout: c.Timeout}
		}),
		recordEvents,                // record events for 5_job.json
		redismem.FxOption(tb),       // in-memory Redis
		esmock.FxOption(transport),  // Elasticsearch over the mock transport
		clockmem.FxOption(fix),      // pin the clock when a case sets "now"
		fx.Populate(&httpServer, &db, &jobs),
	)

	return libe2e.BootResult{
		Server: httpServer,        // any http.Handler
		DB:     db,
		Events: events.Recorded,   // func() []libe2e.RecordedEvent
		Jobs:   jobs,              // runs a case's "job" by name (optional)
		// Redis: redisServer,     // expose for 7_redis.json assertions (use redismem.Server)
	}
}
```

> **`RunE2ETest`.** Most Yokai apps already have a `RunTest` helper that boots the
> Fx app against an in-memory MySQL and runs migrations. For E2E, add a variant
> that runs migrations but **skips Go seed data** — each case owns its initial
> state via `1_database.json`.

### 3. Write the test entrypoint

```go
package e2e_test

import (
	"testing"

	"github.com/myorg/myapp/internal/e2e" // your package with Boot
	libe2e "github.com/ankorstore/yokai-contrib/fxe2e/e2e"
)

func TestE2E(t *testing.T) {
	runner := &libe2e.Runner{Boot: e2e.Boot} // no Auth wiring: cases authenticate via "auth.type"

	cases := libe2e.DiscoverCases(t, "testdata")
	if len(cases) == 0 {
		t.Fatal("e2e: no test cases discovered under testdata/")
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) { runner.Run(t, c.Dir) })
	}
}
```

### 4. Add your first case

Create `internal/e2e/testdata/<route>/<method>/<name>/` with the seven JSON files
— `TestE2E` picks it up automatically. Run the whole suite with `go test ./...`,
or one case with `go test -run 'TestE2E/<route>/<method>/<name>'`.

---

## Packages

| Package | What it does | App-specific? |
|---------|--------------|---------------|
| `e2e` | the engine: case discovery, fixture loading, DB seeding, mock transport, DB diffing, JSON subset assertions, `Runner`/`BootResult` | generic |
| `auth` | `Ankorstore()` dispatch over all principal types, applied from each case's `auth.type` | Ankorstore; used by `Runner` automatically |
| `pubsubmem` | records published events for `5_job.json` (decorates `gcppubsub.Publisher`) | apps on `gcppubsub` |
| `redismem` | in-process in-memory Redis (miniredis), repoints `*redis.Client` | apps on `fxredis` |
| `esmock` | routes `*elasticsearch.Client` through the mock transport | apps on `fxelasticsearch` |
| `clockmem` | pins `clockwork.Clock` to a case's `now` | apps on `fxclock` |
| `await` | `Until`/`Count`/`NoRows` poll helpers; powers a case's `await` query | generic |
| `recorder` | concurrency-safe event sink (returned by `pubsubmem.Recorder`) | generic |

`pubsubmem`, `redismem`, `esmock` and `clockmem` are opt-in fx-option helpers —
import only what your app uses. `auth` is applied automatically by the `Runner`.

---

## Cost

Same infrastructure cost as ordinary Yokai handler/contract tests: each case
boots the app against an in-memory MySQL (one fresh migrated DB per case), with
miniredis and all outbound HTTP served in-process — **no containers, no network**.
It exercises *more real code* (the full boot and the end-to-end journey), so it's
the right tool for complex multi-layer flows; contract tests remain better for
breadth across the API surface.
