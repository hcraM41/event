package event

import "context"

type Publisher interface {
	Subscriber

	Publish(ctx context.Context, event Event) error
}
