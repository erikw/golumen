package find

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
)

type Finder struct {
	logger *slog.Logger
}

// TODO don't export this type
type FindCollector struct {
	logger  *slog.Logger
	matches []string
}

func New(logger *slog.Logger) *Finder {
	return &Finder{logger: logger}
}

var defaultBlockedPaths = map[string]struct{}{
	// ".":    {},
	".git": {},
}

func (f *Finder) Find(path string, pattern string) (matches []string, err error) {
	f.logger.Info(fmt.Sprintf("💡 Shedding light to %s for %s:\n", path, pattern))

	fc := FindCollector{logger: f.logger}

	// TODO does not follow symlinks. do if cmdline switch -f/--follow
	err = filepath.WalkDir(path, fc.walkDir)
	if err != nil {
		return nil, err
	}

	return fc.matches, nil
}

func (fc *FindCollector) walkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fc.logger.Warn("Could not enter directory. Skipping.", "path", path, "error", err.Error())
		return filepath.SkipDir

	}
	skip := fc.blockPath(d.Name())
	fc.logger.Debug("Walking path", "path", path, "skip", skip)

	if skip {
		return filepath.SkipDir
	} else {
		fc.matches = append(fc.matches, path)
		return nil
	}
}

func (fc *FindCollector) blockPath(baseName string) bool {
	_, ok := defaultBlockedPaths[baseName]
	fc.logger.Debug("Consider to block file", "baseName", baseName, "block", ok)
	return ok
}
