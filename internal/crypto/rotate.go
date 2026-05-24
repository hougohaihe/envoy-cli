package crypto

import "fmt"

// Rotator handles re-encryption of vault contents under a new passphrase.
type Rotator struct {
	vault *Vault
}

// NewRotator creates a Rotator backed by the given Vault.
func NewRotator(v *Vault) *Rotator {
	return &Rotator{vault: v}
}

// Rotate decrypts the vault with oldPassphrase and re-encrypts it with
// newPassphrase, writing the result back to the same path.
func (r *Rotator) Rotate(oldPassphrase, newPassphrase string) error {
	if oldPassphrase == "" {
		return fmt.Errorf("rotate: old passphrase must not be empty")
	}
	if newPassphrase == "" {
		return fmt.Errorf("rotate: new passphrase must not be empty")
	}
	if oldPassphrase == newPassphrase {
		return fmt.Errorf("rotate: new passphrase must differ from old passphrase")
	}

	// Load with old passphrase.
	data, err := r.vault.Load(oldPassphrase)
	if err != nil {
		return fmt.Errorf("rotate: load with old passphrase: %w", err)
	}

	// Re-save with new passphrase.
	if err := r.vault.Save(data, newPassphrase); err != nil {
		return fmt.Errorf("rotate: save with new passphrase: %w", err)
	}

	return nil
}
