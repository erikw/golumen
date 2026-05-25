package find

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
)

type Finder struct {
	logger *slog.Logger
	follow bool
}

type findCollector struct {
	logger       *slog.Logger
	searchRegexp *regexp.Regexp
	follow       bool
	visitedDirs  map[string]struct{}
	matches      []string
}

var defaultBlockedPaths = map[string]struct{}{
	// ".":    {},
	".git": {},
}

func New(logger *slog.Logger, follow bool) *Finder {
	return &Finder{logger: logger, follow: follow}
}

func (f *Finder) Find(path string, pattern string) (matches []string, err error) {
	f.logger.Debug(fmt.Sprintf("💡 Shedding light to %s for %s:", path, pattern))

	var r *regexp.Regexp
	var rErr error
	if r, rErr = regexp.Compile(pattern); rErr != nil {
		fmt.Fprintf(os.Stderr, "Invalid regex pattern \"%s\": %v\n", pattern, rErr.Error())
		os.Exit(1)
	}

	fc := findCollector{
		logger:       f.logger,
		searchRegexp: r,
		follow:       f.follow,
		visitedDirs:  make(map[string]struct{}),
	}

	err = fc.walk(path, true)
	if err != nil {
		return nil, err
	}

	return fc.matches, nil
}

func (fc *findCollector) walk(path string, isRoot bool) error {
	info, err := os.Lstat(path)
	if err != nil {
		if isRoot {
			return err
		}
		fc.logger.Debug("Could not stat path. Skipping.", "path", path, "error", err.Error())
		return nil
	}

	baseName := filepath.Base(path)
	skip := fc.blockPath(baseName)
	fc.logger.Debug("Walking path", "path", path, "skip", skip)
	if skip {
		return nil
	}

	if fc.patternMatch(baseName) {
		fc.matches = append(fc.matches, path)
	}

	if info.IsDir() {
		return fc.walkDirectory(path, isRoot)
	}

	if info.Mode()&fs.ModeSymlink == 0 || !fc.follow {
		return nil
	}

	targetInfo, err := os.Stat(path)
	if err != nil {
		if isRoot {
			return err
		}
		fc.logger.Debug("Could not resolve symlink target. Skipping.", "path", path, "error", err.Error())
		return nil
	}

	if !targetInfo.IsDir() {
		return nil
	}

	return fc.walkDirectory(path, isRoot)
}

func (fc *findCollector) walkDirectory(path string, isRoot bool) error {
	resolvedPath, err := resolvedPath(path)
	if err != nil {
		if isRoot {
			return err
		}
		fc.logger.Debug("Could not resolve directory path. Skipping.", "path", path, "error", err.Error())
		return nil
	}

	if _, seen := fc.visitedDirs[resolvedPath]; seen {
		fc.logger.Debug("Skipping already visited directory", "path", path, "resolvedPath", resolvedPath)
		return nil
	}
	fc.visitedDirs[resolvedPath] = struct{}{}

	entries, err := os.ReadDir(path)
	if err != nil {
		if isRoot {
			return err
		}
		fc.logger.Debug("Could not enter directory. Skipping.", "path", path, "error", err.Error())
		return nil
	}

	for _, entry := range entries {
		if err := fc.walk(filepath.Join(path, entry.Name()), false); err != nil {
			return err
		}
	}

	return nil
}

func resolvedPath(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	resolvedPath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return "", err
	}

	return filepath.Clean(resolvedPath), nil
}

func (fc *findCollector) blockPath(baseName string) bool {
	_, ok := defaultBlockedPaths[baseName]
	fc.logger.Debug("Consider to block file", "baseName", baseName, "block", ok)
	return ok
}

func (fc *findCollector) patternMatch(baseName string) bool {
	m := fc.searchRegexp.MatchString(baseName)
	fc.logger.Debug("Trying to match file name to pattern", "baseName", baseName, "pattern", fc.searchRegexp.String(), "match", m)
	return m
}
