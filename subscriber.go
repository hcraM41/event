package event

import (
	"context"
	"sync"
)

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

type Sync []Subscriber

func (sub Sync) Handle(ctx context.Context, evt Event) error {
	var err error
	for _, s := range sub {
		if e := s.Handle(ctx, evt); e != nil {
			err = e
		}
	}
	return err
}

type Async []Subscriber

func (sub Async) Handle(ctx context.Context, evt Event) error {
	var (
		wg   sync.WaitGroup
		once sync.Once
		err  error
	)
	wg.Add(len(sub))
	for _, s := range sub {
		go func(a Subscriber) {
			defer wg.Done()
			if e := a.Handle(ctx, evt); e != nil {
				once.Do(func() { err = e })
			}
		}(s)
	}
	wg.Wait()
	return err
}
