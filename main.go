package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

var logger *slog.Logger

var defaultBlockedPaths = map[string]struct{}{
	// ".":    {},
	".git": {},
}

func main() {
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

	// TODO does not follow symlinks.
	err = filepath.WalkDir(path, walkDir)
	if err != nil {
		return nil, err
	}

	// TODO populate from walkDir
	matches = append(matches, "bogus")
	return matches, nil
}

func walkDir(path string, d fs.DirEntry, err error) error {
	skip := blockPath(d.Name())
	logger.Debug("Walking path", "path", path, "skip", skip)
	if skip {
		return filepath.SkipDir
	} else {
		return nil
	}
}

func blockPath(baseName string) bool {
	_, ok := defaultBlockedPaths[baseName]
	logger.Debug("Consider to block file", "baseName", baseName, "block", ok)
	return ok
}
