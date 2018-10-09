package messengertypes

type Messenger interface {
	Subscribe(channel string) (<-chan Message, error)
	Unsubscribe(channel string) error
	Publish(message interface{}, channels ...string) error
	Close() error
}

type Message struct {
	Channel string
	Payload string
}
