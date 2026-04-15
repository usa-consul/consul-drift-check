package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configPath string
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:   "consul-drift-check",
	Short: "Detect configuration drift between Consul KV namespaces",
	Long: `consul-drift-check compares Consul KV namespaces across datacenters
and reports any configuration drift found between source and destination.`,
	RunE: runDriftCheck,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&configPath, "config", "c", "",
		"path to configuration file (required)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&outputFormat, "output", "o", "text",
		"output format: text or json",
	)
	_ = rootCmd.MarkPersistentFlagRequired("config")
}
