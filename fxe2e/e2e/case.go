package e2e

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// Case is a discovered test case: its subtest Name (the path relative to the
// testdata root, e.g. "qr_campaigns/post/create") and its directory Dir.
type Case struct {
	Name string
	Dir  string
}

// DiscoverCases walks root and returns every case directory — identified by the
// presence of the request fixture — sorted by name for deterministic ordering.
func DiscoverCases(tb testing.TB, root string) []Case {
	tb.Helper()

	var cases []Case
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != fileRequest {
			return nil
		}

		dir := filepath.Dir(path)
		rel, relErr := filepath.Rel(root, dir)
		if relErr != nil {
			rel = dir
		}
		cases = append(cases, Case{Name: filepath.ToSlash(rel), Dir: dir})

		return nil
	})
	if err != nil {
		tb.Fatalf("e2e: discover cases under %q: %v", root, err)
	}

	sort.Slice(cases, func(i, j int) bool { return cases[i].Name < cases[j].Name })

	return cases
}
