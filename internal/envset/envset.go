package envset

import (
	"errors"
	"regexp"
)

// EnvSet represents a named collection of environment variables.
type EnvSet struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
}

var validKeyPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// New creates a new EnvSet with the given name.
func New(name string) (*EnvSet, error) {
	if name == "" {
		return nil, errors.New("envset name must not be empty")
	}
	return &EnvSet{
		Name:      name,
		Variables: make(map[string]string),
	}, nil
}

// Set adds or updates an environment variable in the set.
func (e *EnvSet) Set(key, value string) error {
	if !validKeyPattern.MatchString(key) {
		return errors.New("invalid environment variable key: " + key)
	}
	e.Variables[key] = value
	return nil
}

// Get retrieves the value of an environment variable by key.
func (e *EnvSet) Get(key string) (string, bool) {
	val, ok := e.Variables[key]
	return val, ok
}

// Delete removes an environment variable from the set.
func (e *EnvSet) Delete(key string) {
	delete(e.Variables, key)
}

// Keys returns all variable keys in the set.
func (e *EnvSet) Keys() []string {
	keys := make([]string, 0, len(e.Variables))
	for k := range e.Variables {
		keys = append(keys, k)
	}
	return keys
}

// Merge merges another EnvSet into this one. Existing keys are overwritten.
func (e *EnvSet) Merge(other *EnvSet) {
	for k, v := range other.Variables {
		e.Variables[k] = v
	}
}
