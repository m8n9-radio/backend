package events

import (
	"context"
	"sync"

	appshared "hub/internal/application/shared"
	"hub/internal/domain/shared"
)

// EventHandler is a function that handles a domain event.
type EventHandler func(ctx context.Context, event shared.DomainEvent) error

// InMemoryPublisher is an in-memory implementation of EventPublisher.
type InMemoryPublisher struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewInMemoryPublisher creates a new InMemoryPublisher.
func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{
		handlers: make(map[string][]EventHandler),
	}
}

// Ensure InMemoryPublisher implements EventPublisher.
var _ appshared.EventPublisher = (*InMemoryPublisher)(nil)

// Register registers a handler for a specific event type.
func (p *InMemoryPublisher) Register(eventName string, handler EventHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[eventName] = append(p.handlers[eventName], handler)
}

// Publish publishes a domain event to all registered handlers.
func (p *InMemoryPublisher) Publish(ctx context.Context, event shared.DomainEvent) error {
	p.mu.RLock()
	handlers := p.handlers[event.EventName()]
	p.mu.RUnlock()

	for _, handler := range handlers {
		// Run handlers asynchronously
		go func(h EventHandler) {
			_ = h(ctx, event)
		}(handler)
	}

	return nil
}

// PublishAll publishes multiple domain events.
func (p *InMemoryPublisher) PublishAll(ctx context.Context, events []shared.DomainEvent) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
