package events

import (
	"context"
)

type Event interface {
	EventName() string
	AggregateID() string
	Payload() any
	CreatedAt() int64
}

type EventsPublisher interface {
	Publish(ctx context.Context, event Event) error
}
