package audit

import "time"

// EventType represents the kind of audit event.
type EventType string

const (
	EventPush   EventType = "push"
	EventPull   EventType = "pull"
	EventRotate EventType = "rotate"
	EventSet    EventType = "set"
	EventDelete EventType = "delete"
)

// Event represents a single auditable action.
type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	EnvSet    string    `json:"env_set"`
	Key       string    `json:"key,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user,omitempty"`
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
}

// NewEvent creates a new Event with the current timestamp.
func NewEvent(eventType EventType, envSet string, success bool) Event {
	return Event{
		ID:        newEventID(),
		Type:      eventType,
		EnvSet:    envSet,
		Timestamp: time.Now().UTC(),
		Success:   success,
	}
}

func newEventID() string {
	return time.Now().UTC().Format("20060102150405.000000000")
}
