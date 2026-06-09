package e2e

import (
	"database/sql"
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

// seedDatabase inserts the fixture's initial rows into their tables.
//
// Foreign-key checks are disabled for the duration of the insert so that table
// order in the fixture file is irrelevant — rows referencing each other can be
// listed in any order. Checks are restored afterwards.
func seedDatabase(tb testing.TB, db *sql.DB, tables map[string][]Row) {
	tb.Helper()

	if len(tables) == 0 {
		return
	}

	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		tb.Fatalf("e2e: disable FK checks: %v", err)
	}
	defer func() {
		if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
			tb.Fatalf("e2e: restore FK checks: %v", err)
		}
	}()

	// Deterministic table order keeps failures reproducible. Reserved keys
	// (see ReservedKeyPrefix) are harness datasets, not SQL tables.
	names := make([]string, 0, len(tables))
	for name := range tables {
		if strings.HasPrefix(name, ReservedKeyPrefix) {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	for _, table := range names {
		rows := tables[table]
		if len(rows) == 0 {
			continue
		}
		insertRows(tb, db, table, rows)
	}
}

func insertRows(tb testing.TB, db *sql.DB, table string, rows []Row) {
	tb.Helper()

	// Use the column set of the first row; rows in a table are expected to be
	// homogeneous. Columns are sorted for stable SQL.
	cols := make([]string, 0, len(rows[0]))
	for col := range rows[0] {
		cols = append(cols, col)
	}
	sort.Strings(cols)

	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto(table).Cols(cols...)

	for _, r := range rows {
		values := make([]any, len(cols))
		for i, col := range cols {
			values[i] = normalizeSeedValue(r[col])
		}
		ib.Values(values...)
	}

	query, args := ib.Build()
	if _, err := db.Exec(query, args...); err != nil {
		tb.Fatalf("e2e: seed table %q: %v\nquery: %s", table, err, query)
	}
}

// normalizeSeedValue converts a decoded JSON value into something the SQL driver
// stores faithfully. json.Number becomes an int64 when integral, else float64,
// so integer columns are not written as "1.0".
func normalizeSeedValue(v any) any {
	num, ok := v.(json.Number)
	if !ok {
		return v
	}

	if i, err := num.Int64(); err == nil {
		return i
	}
	if f, err := num.Float64(); err == nil {
		return f
	}

	return num.String()
}
