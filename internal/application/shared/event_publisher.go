package shared

import (
	"context"

	"hub/internal/domain/shared"
)

// EventPublisher defines the interface for publishing domain events.
type EventPublisher interface {
	// Publish publishes a domain event.
	Publish(ctx context.Context, event shared.DomainEvent) error

	// PublishAll publishes multiple domain events.
	PublishAll(ctx context.Context, events []shared.DomainEvent) error
}
