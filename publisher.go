package event

import "context"

type Publisher interface {
	Subscriber

	Publish(ctx context.Context, event Event) error
}

type Mapping map[Type]Subscriber

func NewMapping() Mapping {
	return make(Mapping)
}

func (pub Mapping) On(typ Type, sub Subscriber) Mapping {
	if s, ok := pub[typ]; ok {
		if o, ok := s.(Sync); ok {
			pub[typ] = append(o, sub)
		} else {
			pub[typ] = Sync{s, sub}
		}
	} else {
		pub[typ] = sub
	}
	return pub
}

func (pub Mapping) Publish(ctx context.Context, evt Event) error {
	if sub, ok := pub[evt.Type()]; ok {
		return sub.Handle(ctx, evt)
	}
	return nil
}

func (pub Mapping) Handle(ctx context.Context, evt Event) error {
	return pub.Publish(ctx, evt)
}

type Buffer struct {
	publisher Publisher
	events    []Event
}

func NewBuffer(pub Publisher) *Buffer {
	return &Buffer{publisher: pub}
}

func (pub *Buffer) Publish(ctx context.Context, evt Event) error {
	pub.events = append(pub.events, evt)
	return nil
}

func (pub *Buffer) Handle(ctx context.Context, evt Event) error {
	return pub.Publish(ctx, evt)
}

func (pub *Buffer) Dispatch(ctx context.Context) error {
	var (
		evt Event
		err error
	)
	for len(pub.events) != 0 {
		evt, pub.events = pub.events[0], pub.events[1:]
		if e := pub.publisher.Publish(ctx, evt); e != nil {
			err = e
		}
	}
	return err
}
