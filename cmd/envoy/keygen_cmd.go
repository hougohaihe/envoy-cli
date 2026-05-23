package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/your-org/envoy-cli/internal/crypto"
)

// buildKeygenCmd returns the "keygen" subcommand which generates a new
// random encryption key and prints it along with its fingerprint.
func buildKeygenCmd() *cobra.Command {
	var showFingerprint bool

	cmd := &cobra.Command{
		Use:   "keygen",
		Short: "Generate a new random encryption key",
		Long: `Generate a cryptographically secure random key suitable for
encrypting envoy vaults. Store the output in a safe location such as
a password manager or secret store.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := crypto.GenerateKey()
			if err != nil {
				return fmt.Errorf("failed to generate key: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), key)

			if showFingerprint {
				fp := crypto.FingerprintKey(key)
				fmt.Fprintf(cmd.OutOrStdout(), "fingerprint: %s\n", fp)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&showFingerprint, "fingerprint", false,
		"also print a short fingerprint of the generated key")

	return cmd
}
