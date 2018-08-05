package messengertypes

type Message struct {
	Channel string
	Payload string
}

type Messenger interface {
	Subscribe(channel string) (<-chan Message, error)
	Unsubscribe(channel string) error
	Publish(message interface{}, channels ...string) error
	Close() error
}
