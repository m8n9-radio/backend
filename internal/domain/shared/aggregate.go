package shared

// AggregateRoot is the base type for all aggregate roots.
// It provides domain event collection functionality.
type AggregateRoot struct {
	events []DomainEvent
}

// AddEvent adds a domain event to the aggregate.
func (a *AggregateRoot) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

// Events returns all collected domain events.
func (a *AggregateRoot) Events() []DomainEvent {
	return a.events
}

// ClearEvents removes all collected events.
func (a *AggregateRoot) ClearEvents() {
	a.events = nil
}

// HasEvents returns true if there are pending events.
func (a *AggregateRoot) HasEvents() bool {
	return len(a.events) > 0
}
