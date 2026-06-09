# Changelog

## Unreleased

### Features

* **fxe2e:** declarative, file-driven end-to-end test harness for Yokai HTTP services — each case is a directory of JSON fixtures (no Go code), booting the real app against in-memory infrastructure (no containers, no network)
* **fxe2e:** seven-file case format (`1_database.json` … `7_redis.json`), all mandatory, with strict decoding that fails on unknown/misspelled keys
* **fxe2e:** HTTP mock transport matching on method/host/path/path-prefix/query/headers/JSON-body subset, with optional `expectedCalls` to assert exact usage; unmatched outbound calls fail the case
* **fxe2e:** response assertions on status, headers, JSON-body subset (with `ignore` paths and `unordered` arrays) and raw non-JSON bodies
* **fxe2e:** database before/after diffing per table (`created`/`updated`/`deleted`, composite keys, `ignoreColumns`, `exact`)
* **fxe2e:** published-event assertions (count-exact, order-independent, payload subset)
* **fxe2e:** Redis state assertions via `7_redis.json` (`values`/`exists`/`absent`)
* **fxe2e:** background-job invocation by name (`job`) through a `JobRunner`, and effect-based `await` (poll a `COUNT(*)` query until it settles) for fire-and-forget work
* **fxe2e:** deterministic clock via a case's `now`
* **auth:** zero-config authentication dispatching on a case's `auth.type` across all Ankorstore principals (guest/machine/admin/brand/retailer/impersonation), plus an orthogonal request `origin` (internal/external)
* **fxe2e:** opt-in fx-option helpers — `pubsubmem` (records published events off `gcppubsub.Publisher`), `redismem` (in-memory Redis), `esmock` (Elasticsearch over the mock transport), `clockmem` (pins `clockwork.Clock`)
