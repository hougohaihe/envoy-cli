package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/crypto"
)

// buildRotateCmd returns the "rotate" subcommand which re-encrypts the local
// vault under a new passphrase without changing any stored environment data.
func buildRotateCmd() *cobra.Command {
	var oldPass string
	var newPass string
	var vaultPath string

	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "Re-encrypt the local vault with a new passphrase",
		Long: `Rotate decrypts your local envoy vault using the current passphrase
and immediately re-encrypts it with a new passphrase. No environment
data is modified during this operation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if oldPass == "" {
				return fmt.Errorf("--old-passphrase is required")
			}
			if newPass == "" {
				return fmt.Errorf("--new-passphrase is required")
			}

			v := crypto.NewVault(vaultPath)
			r := crypto.NewRotator(v)

			if err := r.Rotate(oldPass, newPass); err != nil {
				return fmt.Errorf("passphrase rotation failed: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Passphrase rotated successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&oldPass, "old-passphrase", "", "Current vault passphrase")
	cmd.Flags().StringVar(&newPass, "new-passphrase", "", "New vault passphrase")
	cmd.Flags().StringVar(&vaultPath, "vault", defaultVaultPath(), "Path to the vault file")

	return cmd
}
