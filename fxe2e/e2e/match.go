package e2e

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
)

// findRowFailure reports whether some final row matches want on the expected
// columns and whose key's presence in the initial state equals wantPreexisting
// (true for updates, false for creations). On failure it returns a
// human-readable explanation that points at the closest candidate — the row
// sharing want's key — and names the specific columns that differ, rather than
// leaving the caller to eyeball a dump of every row.
func findRowFailure(want Row, finalRows []map[string]string, ignore, keyCols []string, initialKeys map[string]struct{}, wantPreexisting bool) (bool, string) {
	wantKey := jsonRowKey(want, keyCols)

	valueMatchWrongState := false
	var keyMatch map[string]string
	for _, got := range finalRows {
		gotKey := stringRowKey(got, keyCols)
		if rowMatches(want, got, ignore) {
			if _, preexisting := initialKeys[gotKey]; preexisting == wantPreexisting {
				return true, ""
			}
			valueMatchWrongState = true
		}
		if gotKey == wantKey {
			keyMatch = got
		}
	}

	switch {
	case valueMatchWrongState && !wantPreexisting:
		return false, "  a row with these values exists, but its key was already present before the request (expected a newly-created row)"
	case valueMatchWrongState && wantPreexisting:
		return false, "  a row with these values exists, but its key was absent before the request (expected a pre-existing row)"
	case keyMatch != nil:
		return false, fmt.Sprintf("  a row with key %s exists but differs:\n    %s",
			keyDesc(want, keyCols), strings.Join(rowDiff(want, keyMatch, ignore), "\n    "))
	default:
		return false, fmt.Sprintf("  no row with key %s is present after the request\n  actual rows: %v",
			keyDesc(want, keyCols), finalRows)
	}
}

// rowDiff lists, per column, how got fails to satisfy want (ignored columns and
// columns absent from want are skipped). Output is column-sorted for stability.
func rowDiff(want Row, got map[string]string, ignore []string) []string {
	ignored := make(map[string]struct{}, len(ignore))
	for _, c := range ignore {
		ignored[c] = struct{}{}
	}

	var diffs []string
	for _, col := range sortedRowCols(want) {
		if _, skip := ignored[col]; skip {
			continue
		}
		gotVal, present := got[col]
		if !present {
			diffs = append(diffs, fmt.Sprintf("column %q: missing from actual row", col))

			continue
		}
		if !dbValueEqual(want[col], gotVal) {
			diffs = append(diffs, fmt.Sprintf("column %q: expected %v, got %q", col, want[col], gotVal))
		}
	}

	return diffs
}

// keyDesc renders want's key columns as "col=val, ..." for failure messages.
func keyDesc(want Row, keyCols []string) string {
	parts := make([]string, len(keyCols))
	for i, col := range keyCols {
		parts[i] = fmt.Sprintf("%s=%v", col, want[col])
	}

	return strings.Join(parts, ", ")
}

func sortedRowCols(r Row) []string {
	cols := make([]string, 0, len(r))
	for col := range r {
		cols = append(cols, col)
	}
	sort.Strings(cols)

	return cols
}

func initialKeySet(rows []Row, keyCols []string) map[string]struct{} {
	set := make(map[string]struct{}, len(rows))
	for _, r := range rows {
		set[jsonRowKey(r, keyCols)] = struct{}{}
	}

	return set
}

func finalKeySet(rows []map[string]string, keyCols []string) map[string]struct{} {
	set := make(map[string]struct{}, len(rows))
	for _, r := range rows {
		set[stringRowKey(r, keyCols)] = struct{}{}
	}

	return set
}

// jsonRowKey builds a row's identity from a decoded-JSON row (fixture side).
func jsonRowKey(r Row, keyCols []string) string {
	parts := make([]string, len(keyCols))
	for i, col := range keyCols {
		parts[i] = jsonValueToString(r[col])
	}

	return strings.Join(parts, "\x00")
}

// stringRowKey builds a row's identity from a queried row (DB side).
func stringRowKey(r map[string]string, keyCols []string) string {
	parts := make([]string, len(keyCols))
	for i, col := range keyCols {
		parts[i] = r[col]
	}

	return strings.Join(parts, "\x00")
}

func jsonValueToString(v any) string {
	switch w := v.(type) {
	case nil:
		return ""
	case bool:
		if w {
			return "1"
		}

		return "0"
	case json.Number:
		return w.String()
	case string:
		return w
	default:
		return fmt.Sprintf("%v", w)
	}
}

// --- JSON subset matching -------------------------------------------------

func parseJSON(g gomega.Gomega, data []byte, what string) any {
	var v any
	g.Expect(json.Unmarshal(data, &v)).To(gomega.Succeed(), "failed to parse %s", what)

	return v
}

// removePath deletes the leaf at a dotted path (e.g. "data.attributes.createdAt")
// from a decoded JSON tree. Missing intermediate keys are a no-op.
func removePath(root any, path string) {
	parts := strings.Split(path, ".")
	cur, ok := root.(map[string]any)
	if !ok {
		return
	}
	for i, part := range parts {
		if i == len(parts)-1 {
			delete(cur, part)

			return
		}
		next, ok := cur[part].(map[string]any)
		if !ok {
			return
		}
		cur = next
	}
}

// subsetMatch returns a list of human-readable mismatches where expected is not
// satisfied by actual. Maps are matched key-by-key (extra actual keys ignored),
// slices element-by-element with equal length required, scalars by value.
func subsetMatch(expected, actual any, path string) []string {
	return matcher{}.match(expected, actual, path)
}

// subsetMatchUnordered is subsetMatch where the arrays at the given dotted paths
// (e.g. "data" or "data.items") are matched order-independently: lengths must be
// equal and each expected element must match a distinct actual element.
func subsetMatchUnordered(expected, actual any, path string, unordered []string) []string {
	set := make(map[string]bool, len(unordered))
	for _, p := range unordered {
		set[normMatchPath(p)] = true
	}

	return matcher{unordered: set}.match(expected, actual, path)
}

// matcher carries options for a subset comparison.
type matcher struct {
	// unordered holds normalised dotted paths whose arrays match order-independently.
	unordered map[string]bool
}

func (m matcher) match(expected, actual any, path string) []string {
	switch want := expected.(type) {
	case map[string]any:
		return m.matchMap(want, actual, path)
	case []any:
		return m.matchSlice(want, actual, path)
	default:
		if !scalarEqual(want, actual) {
			return []string{fmt.Sprintf("%s: expected %v (%T), got %v (%T)", path, want, want, actual, actual)}
		}

		return nil
	}
}

func (m matcher) matchMap(want map[string]any, actual any, path string) []string {
	got, ok := actual.(map[string]any)
	if !ok {
		return []string{fmt.Sprintf("%s: expected object, got %T", path, actual)}
	}

	var out []string
	for _, k := range sortedKeys(want) {
		gv, present := got[k]
		if !present {
			out = append(out, fmt.Sprintf("%s.%s: missing", path, k))

			continue
		}
		out = append(out, m.match(want[k], gv, path+"."+k)...)
	}

	return out
}

func (m matcher) matchSlice(want []any, actual any, path string) []string {
	got, ok := actual.([]any)
	if !ok {
		return []string{fmt.Sprintf("%s: expected array, got %T", path, actual)}
	}
	if len(want) != len(got) {
		return []string{fmt.Sprintf("%s: expected %d elements, got %d", path, len(want), len(got))}
	}

	if m.unordered[normMatchPath(path)] {
		return m.matchSliceUnordered(want, got, path)
	}

	var out []string
	for i := range want {
		out = append(out, m.match(want[i], got[i], fmt.Sprintf("%s[%d]", path, i))...)
	}

	return out
}

// matchSliceUnordered matches each expected element against some distinct,
// not-yet-consumed actual element (greedy), regardless of order.
func (m matcher) matchSliceUnordered(want, got []any, path string) []string {
	used := make([]bool, len(got))

	var out []string
	for i := range want {
		matched := false
		for j := range got {
			if used[j] {
				continue
			}
			if len(m.match(want[i], got[j], "")) == 0 {
				used[j] = true
				matched = true

				break
			}
		}
		if !matched {
			out = append(out, fmt.Sprintf("%s[%d]: no unmatched element of the actual array matches %v", path, i, want[i]))
		}
	}

	return out
}

// normMatchPath strips the leading "$"/"." so a comparison path like
// "$.data.items" matches a fixture-declared path like "data.items".
func normMatchPath(p string) string {
	p = strings.TrimPrefix(p, "$")
	p = strings.TrimPrefix(p, ".")

	return p
}

func scalarEqual(want, got any) bool {
	if want == nil || got == nil {
		return want == nil && got == nil
	}

	return fmt.Sprintf("%v", normalizeNumber(want)) == fmt.Sprintf("%v", normalizeNumber(got))
}

func normalizeNumber(v any) any {
	if n, ok := v.(json.Number); ok {
		if i, err := n.Int64(); err == nil {
			return i
		}
		if f, err := n.Float64(); err == nil {
			return f
		}
	}

	return v
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

// --- DB row matching ------------------------------------------------------

func queryRows(tb testing.TB, db *sql.DB, table string) []map[string]string {
	tb.Helper()

	//nolint:gosec // table name comes from a trusted local fixture file
	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		tb.Fatalf("e2e: query table %q: %v", table, err)
	}
	defer func() { _ = rows.Close() }()

	cols, err := rows.Columns()
	if err != nil {
		tb.Fatalf("e2e: columns for %q: %v", table, err)
	}

	var out []map[string]string
	for rows.Next() {
		raw := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range raw {
			ptrs[i] = &raw[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			tb.Fatalf("e2e: scan %q: %v", table, err)
		}

		m := make(map[string]string, len(cols))
		for i, col := range cols {
			m[col] = cellToString(raw[i])
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		tb.Fatalf("e2e: iterate %q: %v", table, err)
	}

	return out
}

func rowMatches(want Row, got map[string]string, ignore []string) bool {
	ignored := make(map[string]struct{}, len(ignore))
	for _, c := range ignore {
		ignored[c] = struct{}{}
	}

	for col, wantVal := range want {
		if _, skip := ignored[col]; skip {
			continue
		}
		gotVal, present := got[col]
		if !present {
			return false
		}
		if !dbValueEqual(wantVal, gotVal) {
			return false
		}
	}

	return true
}

func dbValueEqual(want any, got string) bool {
	switch w := want.(type) {
	case nil:
		return got == "" // NULL is rendered as empty
	case bool:
		return (got == "1" || strings.EqualFold(got, "true")) == w
	case json.Number:
		return w.String() == got
	case string:
		return w == got
	default:
		return fmt.Sprintf("%v", w) == got
	}
}

func cellToString(v any) string {
	switch c := v.(type) {
	case nil:
		return ""
	case []byte:
		return string(c)
	case string:
		return c
	case time.Time:
		return formatTimeCell(c)
	default:
		return formatScalarCell(c)
	}
}

func formatScalarCell(v any) string {
	switch c := v.(type) {
	case int64:
		return strconv.FormatInt(c, 10)
	case float64:
		return strconv.FormatFloat(c, 'f', -1, 64)
	case bool:
		if c {
			return "1"
		}

		return "0"
	default:
		return fmt.Sprintf("%v", c)
	}
}

// formatTimeCell renders a DATE column (midnight UTC) date-only so fixtures can
// use "2030-01-01"; DATETIME/TIMESTAMP values keep the time component.
func formatTimeCell(t time.Time) string {
	if t.Equal(t.Truncate(24*time.Hour)) && t.Location() == time.UTC {
		return t.Format("2006-01-02")
	}

	return t.Format("2006-01-02 15:04:05")
}
