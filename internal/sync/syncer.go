package sync

import (
	"fmt"
	"time"

	"github.com/envoy-cli/envoy-cli/internal/envset"
)

// SyncDirection indicates the direction of synchronization.
type SyncDirection int

const (
	DirectionPush SyncDirection = iota
	DirectionPull
)

// SyncResult holds the outcome of a sync operation.
type SyncResult struct {
	Added   []string
	Updated []string
	Deleted []string
	SyncedAt time.Time
}

// Remote is an interface for remote env set storage backends.
type Remote interface {
	Fetch(name string) (*envset.EnvSet, error)
	Push(es *envset.EnvSet) error
}

// Syncer handles push/pull operations between local store and a remote.
type Syncer struct {
	store  *envset.Store
	remote Remote
}

// New creates a new Syncer.
func New(store *envset.Store, remote Remote) *Syncer {
	return &Syncer{store: store, remote: remote}
}

// Push uploads a local env set to the remote.
func (s *Syncer) Push(name string) (*SyncResult, error) {
	local, err := s.store.Get(name)
	if err != nil {
		return nil, fmt.Errorf("push: local get failed: %w", err)
	}
	if err := s.remote.Push(local); err != nil {
		return nil, fmt.Errorf("push: remote push failed: %w", err)
	}
	return &SyncResult{
		Updated:  []string{name},
		SyncedAt: time.Now(),
	}, nil
}

// Pull downloads a remote env set and merges it into the local store.
func (s *Syncer) Pull(name string) (*SyncResult, error) {
	remoteSet, err := s.remote.Fetch(name)
	if err != nil {
		return nil, fmt.Errorf("pull: remote fetch failed: %w", err)
	}

	result := &SyncResult{SyncedAt: time.Now()}

	local, localErr := s.store.Get(name)
	if localErr != nil {
		// env set doesn't exist locally yet
		if err := s.store.Save(remoteSet); err != nil {
			return nil, fmt.Errorf("pull: save failed: %w", err)
		}
		result.Added = append(result.Added, name)
		return result, nil
	}

	merged, added, updated := mergeEnvSets(local, remoteSet)
	if err := s.store.Save(merged); err != nil {
		return nil, fmt.Errorf("pull: save merged failed: %w", err)
	}
	result.Added = added
	result.Updated = updated
	return result, nil
}

// mergeEnvSets merges src into dst, returning the merged set and change lists.
func mergeEnvSets(dst, src *envset.EnvSet) (*envset.EnvSet, []string, []string) {
	var added, updated []string
	for k, v := range src.Vars() {
		if _, exists := dst.Get(k); !exists {
			added = append(added, k)
		} else {
			updated = append(updated, k)
		}
		dst.Set(k, v)
	}
	return dst, added, updated
}
