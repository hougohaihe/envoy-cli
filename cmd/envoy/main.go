// main is the entry point for the envoy-cli application.
// It wires together configuration, authentication, encryption, and sync
// subsystems and registers the top-level CLI commands.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envoy-cli/envoy/internal/auth"
	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/crypto"
	"github.com/envoy-cli/envoy/internal/envset"
	"github.com/envoy-cli/envoy/internal/sync"
)

// version is injected at build time via -ldflags.
var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load application config (creates defaults if missing).
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Set up auth manager backed by the OS credential store path.
	authManager, err := auth.NewManager()
	if err != nil {
		return fmt.Errorf("initialising auth: %w", err)
	}

	// Build the root cobra command.
	root := buildRoot(cfg, authManager)
	return root.Execute()
}

// buildRoot constructs the cobra command tree and returns the root command.
func buildRoot(cfg *config.Config, authManager *auth.Manager) *cobra.Command {
	root := &cobra.Command{
		Use:     "envoy",
		Short:   "Manage and sync environment variable sets",
		Version: version,
		// Silence default error printing — we handle it in main.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Persistent flags available to every sub-command.
	root.PersistentFlags().StringP("env", "e", "default", "environment set name to operate on")
	root.PersistentFlags().StringP("passphrase", "p", "", "passphrase for vault encryption (overrides ENVOY_PASSPHRASE)")

	// Register sub-commands.
	root.AddCommand(
		newSetCmd(cfg),
		newGetCmd(cfg),
		newDeleteCmd(cfg),
		newListCmd(cfg),
		newPushCmd(cfg, authManager),
		newPullCmd(cfg, authManager),
		newLoginCmd(authManager),
		newLogoutCmd(authManager),
	)

	return root
}

// resolvePassphrase returns the passphrase from the flag, falling back to the
// ENVOY_PASSPHRASE environment variable.
func resolvePassphrase(cmd *cobra.Command) (string, error) {
	pp, err := cmd.Flags().GetString("passphrase")
	if err != nil {
		return "", err
	}
	if pp == "" {
		pp = os.Getenv("ENVOY_PASSPHRASE")
	}
	if pp == "" {
		return "", fmt.Errorf("passphrase required: set --passphrase or ENVOY_PASSPHRASE")
	}
	return pp, nil
}

// openVaultedStore is a helper that opens (or creates) the encrypted vault for
// the given environment name and returns a ready-to-use envset.Store.
func openVaultedStore(cmd *cobra.Command, cfg *config.Config, envName string) (*envset.Store, *crypto.Vault, error) {
	pp, err := resolvePassphrase(cmd)
	if err != nil {
		return nil, nil, err
	}

	vaultPath := cfg.VaultPathFor(envName)
	enc, err := crypto.NewEncryptor(pp)
	if err != nil {
		return nil, nil, fmt.Errorf("creating encryptor: %w", err)
	}

	vault, err := crypto.NewVault(vaultPath, enc)
	if err != nil {
		return nil, nil, fmt.Errorf("opening vault: %w", err)
	}

	store, err := envset.NewStore(vault)
	if err != nil {
		return nil, nil, fmt.Errorf("opening store: %w", err)
	}

	return store, vault, nil
}

// newSyncer builds a Syncer wired to the configured remote endpoint.
func newSyncer(cfg *config.Config, authManager *auth.Manager) (*sync.Syncer, error) {
	token, err := authManager.CurrentToken()
	if err != nil {
		return nil, fmt.Errorf("retrieving auth token: %w", err)
	}
	remote := sync.NewHTTPRemote(cfg.RemoteURL, token.Raw)
	return sync.New(remote), nil
}
