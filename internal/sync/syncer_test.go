package sync_test

import (
	"errors"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/envset"
	envSync "github.com/envoy-cli/envoy-cli/internal/sync"
)

// mockRemote is a simple in-memory Remote implementation for testing.
type mockRemote struct {
	data  map[string]*envset.EnvSet
	pushErr error
	fetchErr error
}

func newMockRemote() *mockRemote {
	return &mockRemote{data: make(map[string]*envset.EnvSet)}
}

func (m *mockRemote) Fetch(name string) (*envset.EnvSet, error) {
	if m.fetchErr != nil {
		return nil, m.fetchErr
	}
	es, ok := m.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return es, nil
}

func (m *mockRemote) Push(es *envset.EnvSet) error {
	if m.pushErr != nil {
		return m.pushErr
	}
	m.data[es.Name()] = es
	return nil
}

func TestPush_Success(t *testing.T) {
	store := envset.NewStore()
	es, _ := envset.New("staging")
	es.Set("KEY", "value")
	store.Save(es)

	remote := newMockRemote()
	syncer := envSync.New(store, remote)

	result, err := syncer.Push("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Updated) != 1 || result.Updated[0] != "staging" {
		t.Errorf("expected updated=[staging], got %v", result.Updated)
	}
	if result.SyncedAt.IsZero() {
		t.Error("SyncedAt should not be zero")
	}
}

func TestPush_RemoteError(t *testing.T) {
	store := envset.NewStore()
	es, _ := envset.New("staging")
	store.Save(es)

	remote := newMockRemote()
	remote.pushErr = errors.New("network error")
	syncer := envSync.New(store, remote)

	_, err := syncer.Push("staging")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPull_NewEnvSet(t *testing.T) {
	store := envset.NewStore()
	remote := newMockRemote()

	remoteSet, _ := envset.New("production")
	remoteSet.Set("DB_URL", "postgres://prod")
	remote.data["production"] = remoteSet

	syncer := envSync.New(store, remote)
	result, err := syncer.Pull("production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Added) != 1 || result.Added[0] != "production" {
		t.Errorf("expected added=[production], got %v", result.Added)
	}
}

func TestPull_MergesExisting(t *testing.T) {
	store := envset.NewStore()
	local, _ := envset.New("production")
	local.Set("EXISTING", "old")
	store.Save(local)

	remote := newMockRemote()
	remoteSet, _ := envset.New("production")
	remoteSet.Set("EXISTING", "new")
	remoteSet.Set("NEW_KEY", "hello")
	remote.data["production"] = remoteSet

	syncer := envSync.New(store, remote)
	result, err := syncer.Pull("production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Added) != 1 {
		t.Errorf("expected 1 added key, got %d", len(result.Added))
	}
	if len(result.Updated) != 1 {
		t.Errorf("expected 1 updated key, got %d", len(result.Updated))
	}
}
