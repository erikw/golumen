package cmd

import "testing"

func TestRootCommandArgs(t *testing.T) {
	if err := rootCmd.Args(rootCmd, nil); err == nil {
		t.Fatal("expected missing path argument to fail")
	}

	if err := rootCmd.Args(rootCmd, []string{"."}); err != nil {
		t.Fatalf("expected single path argument to pass: %v", err)
	}

	if err := rootCmd.Args(rootCmd, []string{".", ".."}); err == nil {
		t.Fatal("expected multiple path arguments to fail")
	}
}

func TestRootCommandVersionIsConfigured(t *testing.T) {
	if rootCmd.Version == "" {
		t.Fatal("expected Cobra version support to be configured")
	}
}
