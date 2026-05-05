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

type findCollector struct {
	logger  *slog.Logger
	matches []string
}

var defaultBlockedPaths = map[string]struct{}{
	// ".":    {},
	".git": {},
}

func New(logger *slog.Logger) *Finder {
	return &Finder{logger: logger}
}

func (f *Finder) Find(path string, pattern string) (matches []string, err error) {
	f.logger.Info(fmt.Sprintf("💡 Shedding light to %s for %s:\n", path, pattern))

	fc := findCollector{logger: f.logger}

	// TODO paralellize with goroutines? Makes sense?
	err = filepath.WalkDir(path, fc.walkDir)
	if err != nil {
		return nil, err
	}

	return fc.matches, nil
}

func (fc *findCollector) walkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fc.logger.Warn("Could not enter directory. Skipping.", "path", path, "error", err.Error())
		return filepath.SkipDir

	}
	skip := fc.blockPath(d.Name())
	fc.logger.Debug("Walking path", "path", path, "skip", skip)

	if skip {
		return filepath.SkipDir
	} else {
		// TODO does not follow symlinks. do if cmdline switch -f/--follow
		// Need to resolve paths and store traversed to detect infinite recursion
		// loop and abort
		// if d.Type()&fs.ModeSymlink != 0 {
		//        target, err := fs.ReadLink(fsys, path)
		//        if err != nil {
		//            return err
		//        }
		// here call Find(target,....)
		fc.matches = append(fc.matches, path)
		return nil
	}
}

func (fc *findCollector) blockPath(baseName string) bool {
	_, ok := defaultBlockedPaths[baseName]
	fc.logger.Debug("Consider to block file", "baseName", baseName, "block", ok)
	return ok
}
