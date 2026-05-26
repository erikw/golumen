package find

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestFindDoesNotFollowSymlinkedDirectoriesByDefault(t *testing.T) {
	t.Parallel()

	rootDir := t.TempDir()
	targetDir := filepath.Join(t.TempDir(), "target")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("creating target directory: %v", err)
	}

	matchPath := filepath.Join(targetDir, "match.go")
	writeTestFile(t, matchPath)

	linkPath := filepath.Join(rootDir, "linked-target")
	mustSymlink(t, targetDir, linkPath)

	finder := New(testLogger(), false)
	matches, err := finder.Find(rootDir, `match\.go$`)
	if err != nil {
		t.Fatalf("finding without follow failed: %v", err)
	}

	if len(matches) != 0 {
		t.Fatalf("expected symlinked directory to be skipped without follow, got %v", matches)
	}
}

func TestFindFollowsSymlinkedDirectories(t *testing.T) {
	t.Parallel()

	rootDir := t.TempDir()
	externalDir := filepath.Join(t.TempDir(), "external")
	if err := os.MkdirAll(externalDir, 0o755); err != nil {
		t.Fatalf("creating external directory: %v", err)
	}

	linkPath := filepath.Join(rootDir, "linked-target")
	mustSymlink(t, externalDir, linkPath)

	matchPath := filepath.Join(linkPath, "match.go")
	writeTestFile(t, filepath.Join(externalDir, "match.go"))

	finder := New(testLogger(), true)
	matches, err := finder.Find(rootDir, `match\.go$`)
	if err != nil {
		t.Fatalf("finding with follow failed: %v", err)
	}

	if len(matches) != 1 || matches[0] != matchPath {
		t.Fatalf("expected followed symlink match %q, got %v", matchPath, matches)
	}
}

func TestFindSkipsAlreadyVisitedSymlinkTargets(t *testing.T) {
	t.Parallel()

	rootDir := t.TempDir()
	realDir := filepath.Join(rootDir, "real")
	if err := os.MkdirAll(realDir, 0o755); err != nil {
		t.Fatalf("creating real directory: %v", err)
	}

	matchPath := filepath.Join(realDir, "match.go")
	writeTestFile(t, matchPath)

	linkPath := filepath.Join(rootDir, "alias")
	mustSymlink(t, realDir, linkPath)

	finder := New(testLogger(), true)
	matches, err := finder.Find(rootDir, `match\.go$`)
	if err != nil {
		t.Fatalf("finding with duplicate symlink target failed: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected duplicate target traversal to be skipped, got %v", matches)
	}

	if matches[0] != matchPath && matches[0] != filepath.Join(linkPath, "match.go") {
		t.Fatalf("expected one logical match path for the shared target, got %v", matches)
	}
}

func TestFindAvoidsSymlinkLoops(t *testing.T) {
	t.Parallel()

	rootDir := t.TempDir()
	realDir := filepath.Join(rootDir, "real")
	if err := os.MkdirAll(realDir, 0o755); err != nil {
		t.Fatalf("creating real directory: %v", err)
	}

	matchPath := filepath.Join(realDir, "match.go")
	writeTestFile(t, matchPath)

	loopPath := filepath.Join(realDir, "loop")
	mustSymlink(t, realDir, loopPath)

	finder := New(testLogger(), true)
	matches, err := finder.Find(rootDir, `match\.go$`)
	if err != nil {
		t.Fatalf("finding with symlink loop failed: %v", err)
	}

	if len(matches) != 1 || matches[0] != matchPath {
		t.Fatalf("expected symlink loop to be ignored after first visit, got %v", matches)
	}
}

func TestFindAllowsBlockedRootPath(t *testing.T) {
	t.Parallel()

	parentDir := t.TempDir()
	rootDir := filepath.Join(parentDir, ".git")
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		t.Fatalf("creating blocked root directory: %v", err)
	}

	rootMatchPath := filepath.Join(rootDir, "root-match.go")
	writeTestFile(t, rootMatchPath)

	nestedBlockedDir := filepath.Join(rootDir, "nested", ".git")
	if err := os.MkdirAll(nestedBlockedDir, 0o755); err != nil {
		t.Fatalf("creating nested blocked directory: %v", err)
	}

	nestedMatchPath := filepath.Join(nestedBlockedDir, "nested-match.go")
	writeTestFile(t, nestedMatchPath)

	finder := New(testLogger(), false)
	matches, err := finder.Find(rootDir, `match\.go$`)
	if err != nil {
		t.Fatalf("finding under blocked root failed: %v", err)
	}

	if len(matches) != 1 || matches[0] != rootMatchPath {
		t.Fatalf("expected explicit blocked root to be searched but nested blocked directory to stay skipped, got %v", matches)
	}
}

func TestFindReturnsSortedMatches(t *testing.T) {
	t.Parallel()

	rootDir := t.TempDir()

	expected := []string{
		filepath.Join(rootDir, "a", "match-2.go"),
		filepath.Join(rootDir, "b", "match-1.go"),
		filepath.Join(rootDir, "match-3.go"),
	}

	writeTestFile(t, expected[1])
	writeTestFile(t, expected[0])
	writeTestFile(t, expected[2])

	finder := New(testLogger(), false)
	matches, err := finder.Find(rootDir, `match-[0-9]+\.go$`)
	if err != nil {
		t.Fatalf("finding sorted matches failed: %v", err)
	}

	if len(matches) != len(expected) {
		t.Fatalf("expected %d matches, got %v", len(expected), matches)
	}

	for i := range expected {
		if matches[i] != expected[i] {
			t.Fatalf("expected sorted matches %v, got %v", expected, matches)
		}
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func writeTestFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating parent directory for %q: %v", path, err)
	}

	if err := os.WriteFile(path, []byte("package test\n"), 0o644); err != nil {
		t.Fatalf("writing test file %q: %v", path, err)
	}
}

func mustSymlink(t *testing.T, target string, link string) {
	t.Helper()

	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("creating symlink %q -> %q: %v", link, target, err)
	}
}
