package event

import "context"

type Subscriber interface {
	Handle(ctx context.Context, event Event) error
}

type Func func(ctx context.Context, event Event) error

func (sub Func) Handle(ctx context.Context, evt Event) error {
	if sub == nil {
		return nil
	}
	return sub(ctx, evt)
}
