package event

type EventManager interface {
	RegisterEvent(topic string, callback any) error
	EmitEvent(topic string, arg ...any)
}
