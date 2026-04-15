package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/consul-drift-check/internal/config"
	"github.com/example/consul-drift-check/internal/consul"
	"github.com/example/consul-drift-check/internal/diff"
	"github.com/example/consul-drift-check/internal/report"
)

func runDriftCheck(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	srcConsul, err := consul.NewClient(cfg.Source.Address, cfg.Source.Token)
	if err != nil {
		return fmt.Errorf("creating source consul client: %w", err)
	}

	dstConsul, err := consul.NewClient(cfg.Destination.Address, cfg.Destination.Token)
	if err != nil {
		return fmt.Errorf("creating destination consul client: %w", err)
	}

	srcKV := consul.NewKVClient(srcConsul)
	dstKV := consul.NewKVClient(dstConsul)

	srcEntries, err := srcKV.ListPrefix(cfg.Source.Prefix)
	if err != nil {
		return fmt.Errorf("listing source prefix %q: %w", cfg.Source.Prefix, err)
	}

	dstEntries, err := dstKV.ListPrefix(cfg.Destination.Prefix)
	if err != nil {
		return fmt.Errorf("listing destination prefix %q: %w", cfg.Destination.Prefix, err)
	}

	results := diff.Compare(srcEntries, dstEntries)

	w := report.NewWriter(os.Stdout, outputFormat)
	if err := w.Write(results); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	if len(results) > 0 {
		os.Exit(2)
	}
	return nil
}
