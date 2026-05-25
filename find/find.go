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
}

type findCollector struct {
	logger       *slog.Logger
	searchRegexp *regexp.Regexp
	matches      []string
}

var defaultBlockedPaths = map[string]struct{}{
	// ".":    {},
	".git": {},
}

func New(logger *slog.Logger) *Finder {
	return &Finder{logger: logger}
}

func (f *Finder) Find(path string, pattern string) (matches []string, err error) {
	f.logger.Debug(fmt.Sprintf("💡 Shedding light to %s for %s:", path, pattern))

	var r *regexp.Regexp
	var rErr error
	if r, rErr = regexp.Compile(pattern); rErr != nil {
		fmt.Fprintf(os.Stderr, "Invalid regex pattern \"%s\": %v\n", pattern, rErr.Error())
		os.Exit(1)
	}

	fc := findCollector{logger: f.logger, searchRegexp: r}

	// TODO paralellize with goroutines? Makes sense?
	err = filepath.WalkDir(path, fc.walkDir)
	if err != nil {
		return nil, err
	}

	return fc.matches, nil
}

func (fc *findCollector) walkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fc.logger.Debug("Could not enter directory. Skipping.", "path", path, "error", err.Error())
		return filepath.SkipDir

	}
	skip := fc.blockPath(d.Name())
	fc.logger.Debug("Walking path", "path", path, "skip", skip)

	if skip {
		return filepath.SkipDir
	} else if fc.patternMatch(d.Name()) {
		fc.matches = append(fc.matches, path)
	}

	// TODO does not follow symlinks. do if cmdline switch -f/--follow
	// Need to resolve paths and store traversed to detect infinite recursion
	// loop and abort
	// if d.Type()&fs.ModeSymlink != 0 {
	//        target, err := fs.ReadLink(fsys, path)
	//        if err != nil {
	//            return err
	//        }
	// here call Find(target,....)
	return nil
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
