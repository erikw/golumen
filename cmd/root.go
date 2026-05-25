package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/erikw/golumen/find"
	"github.com/erikw/golumen/internal/version"
	"github.com/spf13/cobra"
)

var debug bool
var follow bool
var logger *slog.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "golumen <pattern> [path]",
	Short:            "Shining a light on your file system.",
	Long:             `A fast, concurrent CLI file finder written in Go, designed to bring transparency to your directory structures with speed and simplicity.`,
	Version:          version.Version,
	Args:             cobra.RangeArgs(1, 2),
	PersistentPreRun: preRun,
	Run:              cmdSearch,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	rootCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow symlinked directories")
}

func preRun(cmd *cobra.Command, args []string) {
	initLogger(debug)
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

func cmdSearch(cmd *cobra.Command, args []string) {
	pattern := args[0]
	path := "."
	if len(args) > 1 {
		path = args[1]
	}

	fmt.Println("Golumen illuminates...")
	// fmt.Printf("cmd: %v\n", cmd)
	// fmt.Printf("args: %v\n", args)

	finder := find.New(logger, follow)
	matches, err := finder.Find(path, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error luminating: %v\n", err)
		os.Exit(1)
	}
	printMatches(matches)
}

func printMatches(matches []string) {
	for _, match := range matches {
		fmt.Printf("- Match: %s\n", match)
	}
}
