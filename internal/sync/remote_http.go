package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/envoy-cli/envoy-cli/internal/envset"
)

// HTTPRemote implements Remote using a REST API backend.
type HTTPRemote struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewHTTPRemote creates an HTTPRemote pointed at the given base URL.
func NewHTTPRemote(baseURL, apiKey string) *HTTPRemote {
	return &HTTPRemote{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

type envSetPayload struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// Fetch retrieves an env set by name from the remote API.
func (h *HTTPRemote) Fetch(name string) (*envset.EnvSet, error) {
	url := fmt.Sprintf("%s/envsets/%s", h.baseURL, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch: build request: %w", err)
	}
	h.setHeaders(req)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("fetch: env set %q not found on remote", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch: unexpected status %d", resp.StatusCode)
	}

	var payload envSetPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("fetch: decode response: %w", err)
	}

	es, err := envset.New(payload.Name)
	if err != nil {
		return nil, fmt.Errorf("fetch: create env set: %w", err)
	}
	for k, v := range payload.Vars {
		es.Set(k, v)
	}
	return es, nil
}

// Push uploads a local env set to the remote API.
func (h *HTTPRemote) Push(es *envset.EnvSet) error {
	payload := envSetPayload{
		Name: es.Name(),
		Vars: es.Vars(),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("push: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/envsets/%s", h.baseURL, es.Name())
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("push: build request: %w", err)
	}
	h.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("push: http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("push: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (h *HTTPRemote) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("Accept", "application/json")
}
