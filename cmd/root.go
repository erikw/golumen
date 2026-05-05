package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/erikw/golumen/find"
	"github.com/erikw/golumen/internal/version"
	"github.com/spf13/cobra"
)

var showVersion bool
var debug bool
var logger *slog.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "golumen",
	Short:            "Shining a light on your file system.",
	Long:             `A fast, concurrent CLI file finder written in Go, designed to bring transparency to your directory structures with speed and simplicity.`,
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
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
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
	fmt.Println("Welcome to Golumen")
	// fmt.Printf("cmd: %v\n", cmd)
	// fmt.Printf("args: %v\n", args)

	if showVersion {
		fmt.Printf("Golumen verison: %s\n", version.Version)
		return
	}

	finder := find.New(logger)
	matches, err := finder.Find(".", "*")
	if err != nil {
		fmt.Printf("Error luminating: %v", err)
	}
	printMatches(matches)
}

func printMatches(matches []string) {
	for _, match := range matches {
		fmt.Printf("- Match: %s\n", match)
	}
}
