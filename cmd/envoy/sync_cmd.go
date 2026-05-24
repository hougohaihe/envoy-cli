package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-cli/internal/config"
	"github.com/yourorg/envoy-cli/internal/sync"
)

func buildSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync environment variable sets with a remote",
	}

	cmd.AddCommand(buildSyncPushCmd())
	cmd.AddCommand(buildSyncPullCmd())

	return cmd
}

func buildSyncPushCmd() *cobra.Command {
	var setName string

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push a local env set to the remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			if cfg.RemoteURL == "" {
				return fmt.Errorf("remote_url not set in config")
			}

			store, err := openVaultedStore(cmd)
			if err != nil {
				return err
			}

			remote := sync.NewHTTPRemote(cfg.RemoteURL, cfg.AuthToken)
			syncer := sync.New(store, remote)

			if err := syncer.Push(cmd.Context(), setName); err != nil {
				return fmt.Errorf("push: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Pushed env set %q to %s\n", setName, cfg.RemoteURL)
			return nil
		},
	}

	cmd.Flags().StringVarP(&setName, "name", "n", "", "Name of the env set to push (required)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func buildSyncPullCmd() *cobra.Command {
	var setName string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull a remote env set into local storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			if cfg.RemoteURL == "" {
				return fmt.Errorf("remote_url not set in config")
			}

			store, err := openVaultedStore(cmd)
			if err != nil {
				return err
			}

			remote := sync.NewHTTPRemote(cfg.RemoteURL, cfg.AuthToken)
			syncer := sync.New(store, remote)

			if err := syncer.Pull(cmd.Context(), setName); err != nil {
				return fmt.Errorf("pull: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Pulled env set %q from %s\n", setName, cfg.RemoteURL)
			return nil
		},
	}

	cmd.Flags().StringVarP(&setName, "name", "n", "", "Name of the env set to pull (required)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
