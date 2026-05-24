package audit

import (
	"testing"
	"time"
)

func makeEvent(envset, action string, ts time.Time) Event {
	return Event{
		ID:        "test-id",
		Timestamp: ts,
		EnvSet:    envset,
		Action:    action,
		Actor:     "user",
		Detail:    "",
	}
}

func TestFilter_MatchEnvSet(t *testing.T) {
	now := time.Now()
	e := makeEvent("prod", "push", now)

	if !(Filter{EnvSet: "prod"}).Match(e) {
		t.Error("expected match on EnvSet=prod")
	}
	if (Filter{EnvSet: "staging"}).Match(e) {
		t.Error("expected no match on EnvSet=staging")
	}
}

func TestFilter_MatchAction(t *testing.T) {
	now := time.Now()
	e := makeEvent("prod", "push", now)

	if !(Filter{Action: "push"}).Match(e) {
		t.Error("expected match on Action=push")
	}
	if (Filter{Action: "pull"}).Match(e) {
		t.Error("expected no match on Action=pull")
	}
}

func TestFilter_MatchSince(t *testing.T) {
	base := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	e := makeEvent("prod", "push", base)

	if !(Filter{Since: base.Add(-time.Hour)}).Match(e) {
		t.Error("expected match when event is after Since")
	}
	if (Filter{Since: base.Add(time.Hour)}).Match(e) {
		t.Error("expected no match when event is before Since")
	}
}

func TestFilter_MatchUntil(t *testing.T) {
	base := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	e := makeEvent("prod", "push", base)

	if !(Filter{Until: base.Add(time.Hour)}).Match(e) {
		t.Error("expected match when event is before Until")
	}
	if (Filter{Until: base.Add(-time.Hour)}).Match(e) {
		t.Error("expected no match when event is after Until")
	}
}

func TestFilter_Apply_Limit(t *testing.T) {
	now := time.Now()
	events := []Event{
		makeEvent("prod", "push", now),
		makeEvent("prod", "pull", now),
		makeEvent("prod", "push", now),
	}

	result := (Filter{Limit: 2}).Apply(events)
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
}

func TestFilter_Apply_NoLimit(t *testing.T) {
	now := time.Now()
	events := []Event{
		makeEvent("prod", "push", now),
		makeEvent("staging", "pull", now),
		makeEvent("prod", "push", now),
	}

	result := (Filter{EnvSet: "prod"}).Apply(events)
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
}

func TestFilter_Apply_Empty(t *testing.T) {
	result := (Filter{}).Apply(nil)
	if result != nil {
		t.Errorf("expected nil result for empty input, got %v", result)
	}
}
