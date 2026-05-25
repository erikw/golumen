package cmd

import "testing"

func TestRootCommandArgs(t *testing.T) {
	if err := rootCmd.Args(rootCmd, nil); err == nil {
		t.Fatal("expected missing pattern argument to fail")
	}

	if err := rootCmd.Args(rootCmd, []string{"*.go"}); err != nil {
		t.Fatalf("expected single pattern argument to pass: %v", err)
	}

	if err := rootCmd.Args(rootCmd, []string{"*.go", "."}); err != nil {
		t.Fatalf("expected pattern and path arguments to pass: %v", err)
	}

	if err := rootCmd.Args(rootCmd, []string{"*.go", ".", ".."}); err == nil {
		t.Fatal("expected more than pattern and optional path to fail")
	}
}

func TestRootCommandVersionIsConfigured(t *testing.T) {
	if rootCmd.Version == "" {
		t.Fatal("expected Cobra version support to be configured")
	}
}
