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

func TestRootCommandFollowFlagIsConfigured(t *testing.T) {
	flag := rootCmd.Flags().Lookup("follow")
	if flag == nil {
		t.Fatal("expected follow flag to be configured")
	}

	if flag.Shorthand != "f" {
		t.Fatalf("expected follow flag shorthand to be f, got %q", flag.Shorthand)
	}

	if flag.DefValue != "false" {
		t.Fatalf("expected follow flag default to be false, got %q", flag.DefValue)
	}
}
