package shared

import "time"

// DomainEvent represents an event that occurred in the domain.
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
	Payload() interface{}
}

// BaseEvent provides common functionality for domain events.
type BaseEvent struct {
	name       string
	occurredAt time.Time
}

// NewBaseEvent creates a new base event with the given name.
func NewBaseEvent(name string) BaseEvent {
	return BaseEvent{
		name:       name,
		occurredAt: time.Now(),
	}
}

// EventName returns the name of the event.
func (e BaseEvent) EventName() string {
	return e.name
}

// OccurredAt returns when the event occurred.
func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}
