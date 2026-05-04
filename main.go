package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type FindCollector struct {
	matches []string
}

var logger *slog.Logger

var defaultBlockedPaths = map[string]struct{}{
	// ".":    {},
	".git": {},
}

func main() {
	// TODO -v --version
	// TODO -h --help
	initLogger(true) // TODO take cli arg --debug or --log-level debug
	fmt.Println("Welcome to Golumen")
	matches, err := find(".", "*")
	if err != nil {
		fmt.Printf("Error luminating: %v", err)
	}
	printMatches(matches)
}

func initLogger(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Remove timestamp from output.
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		}),
	)
}

func printMatches(matches []string) {
	for _, match := range matches {
		fmt.Printf("- Match: %s\n", match)
	}
}

func find(path string, pattern string) (matches []string, err error) {
	logger.Info(fmt.Sprintf("💡 Shedding light to %s for %s:\n", path, pattern))

	fc := FindCollector{}

	// TODO does not follow symlinks. do if cmdline switch -f/--follow
	err = filepath.WalkDir(path, fc.walkDir)
	if err != nil {
		return nil, err
	}

	return fc.matches, nil
}

func (fc *FindCollector) walkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		logger.Warn("Could not enter directory. Skipping.", "path", path, "error", err.Error())
		return filepath.SkipDir

	}
	skip := blockPath(d.Name())
	logger.Debug("Walking path", "path", path, "skip", skip)

	if skip {
		return filepath.SkipDir
	} else {
		fc.matches = append(fc.matches, path)
		return nil
	}
}

func blockPath(baseName string) bool {
	_, ok := defaultBlockedPaths[baseName]
	logger.Debug("Consider to block file", "baseName", baseName, "block", ok)
	return ok
}
