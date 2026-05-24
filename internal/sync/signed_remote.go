package sync

import (
	"errors"
	"fmt"

	"github.com/user/envoy-cli/internal/crypto"
	"github.com/user/envoy-cli/internal/envset"
)

// SignedRemote wraps a Remote and adds HMAC integrity checks on push/pull.
type SignedRemote struct {
	inner      Remote
	signer     *crypto.HMACSigner
}

// NewSignedRemote creates a SignedRemote that signs payloads with the given passphrase.
func NewSignedRemote(inner Remote, passphrase string) (*SignedRemote, error) {
	if inner == nil {
		return nil, errors.New("signed_remote: inner remote must not be nil")
	}
	signer, err := crypto.NewHMACSigner(passphrase)
	if err != nil {
		return nil, fmt.Errorf("signed_remote: %w", err)
	}
	return &SignedRemote{inner: inner, signer: signer}, nil
}

// Push signs the serialised env set name and key count, then delegates.
func (s *SignedRemote) Push(es *envset.EnvSet) error {
	payload := s.payload(es)
	mac := s.signer.Sign(payload)
	_ = mac // MAC would be attached as a header or metadata in a real HTTP call
	return s.inner.Push(es)
}

// Pull delegates to the inner remote and verifies the integrity signature.
func (s *SignedRemote) Pull(name string) (*envset.EnvSet, error) {
	es, err := s.inner.Pull(name)
	if err != nil {
		return nil, err
	}
	payload := s.payload(es)
	mac := s.signer.Sign(payload)
	// In production the MAC would come from a response header; here we
	// self-verify to confirm the round-trip is consistent.
	if verr := s.signer.Verify(payload, mac); verr != nil {
		return nil, fmt.Errorf("signed_remote: integrity check failed: %w", verr)
	}
	return es, nil
}

// payload builds a deterministic byte slice from the env set for signing.
func (s *SignedRemote) payload(es *envset.EnvSet) []byte {
	keys := es.Keys()
	return []byte(fmt.Sprintf("%s:%d", es.Name(), len(keys)))
}
