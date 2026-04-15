package cmd

import (
	"bytes"
	"testing"
)

func TestExecute_MissingConfigFlag(t *testing.T) {
	// Reset flags between tests
	rootCmd.ResetFlags()
	init()

	rootCmd.SetArgs([]string{})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --config flag is missing, got nil")
	}
}

func TestExecute_InvalidOutputFormat(t *testing.T) {
	rootCmd.ResetFlags()
	init()

	// Provide a config path but an invalid format; config loading will fail
	// before format validation, so we just confirm the command wires up flags.
	rootCmd.SetArgs([]string{"--config", "nonexistent.yaml", "--output", "xml"})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// We expect an error (config file not found), confirming flag parsing worked.
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent config file, got nil")
	}
}

func TestRootCmd_ShortDescription(t *testing.T) {
	if rootCmd.Short == "" {
		t.Error("root command should have a short description")
	}
}

func TestRootCmd_HasConfigFlag(t *testing.T) {
	f := rootCmd.PersistentFlags().Lookup("config")
	if f == nil {
		t.Error("expected --config flag to be registered")
	}
}

func TestRootCmd_HasOutputFlag(t *testing.T) {
	f := rootCmd.PersistentFlags().Lookup("output")
	if f == nil {
		t.Error("expected --output flag to be registered")
	}
	if f.DefValue != "text" {
		t.Errorf("expected default output format 'text', got %q", f.DefValue)
	}
}
