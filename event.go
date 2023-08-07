package event

// Type is the event type, The underlying type is int.
type Type int

type Event interface {
	Type() Type
}
