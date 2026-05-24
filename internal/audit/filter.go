package audit

import "time"

// Filter holds criteria for querying audit log events.
type Filter struct {
	EnvSet    string
	Action    string
	Since     time.Time
	Until     time.Time
	Limit     int
}

// Match reports whether the event satisfies all non-zero filter criteria.
func (f Filter) Match(e Event) bool {
	if f.EnvSet != "" && e.EnvSet != f.EnvSet {
		return false
	}
	if f.Action != "" && e.Action != f.Action {
		return false
	}
	if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
		return false
	}
	if !f.Until.IsZero() && e.Timestamp.After(f.Until) {
		return false
	}
	return true
}

// Apply filters a slice of events according to the filter criteria,
// respecting the Limit field (0 means no limit).
func (f Filter) Apply(events []Event) []Event {
	var result []Event
	for _, e := range events {
		if f.Match(e) {
			result = append(result, e)
			if f.Limit > 0 && len(result) >= f.Limit {
				break
			}
		}
	}
	return result
}
