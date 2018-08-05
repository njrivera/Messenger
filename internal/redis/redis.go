package redis

import (
	"errors"
	"fmt"
	"os"

	"github.com/Messenger/pkg/messengertypes"

	"github.com/go-redis/redis"
)

const (
	redisPortEnv = "REDIS_PORT"
)

type RedisMessenger struct {
	client        *redis.Client
	subscriptions map[string]*subscription
}

type subscription struct {
	pubsub  *redis.PubSub
	channel chan messengertypes.Message
}

func NewMessenger() (*RedisMessenger, error) {
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		return nil, errors.New("Redis env var not set")
	}
	messenger := &RedisMessenger{
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:" + redisPort,
		}),
		subscriptions: map[string]*subscription{},
	}
	if _, err := messenger.client.Ping().Result(); err != nil {
		return nil, errors.New(fmt.Sprintf("Error initializing redis client: %+v", err))
	}

	return messenger, nil
}

func (m *RedisMessenger) Subscribe(channel string) (<-chan messengertypes.Message, error) {
	if _, ok := m.subscriptions[channel]; ok {
		return nil, errors.New(fmt.Sprintf("Already subscribed to %s", channel))
	}

	sub := &subscription{
		pubsub:  m.client.Subscribe(channel),
		channel: make(chan messengertypes.Message, 1),
	}

	go func() {
		ch := sub.pubsub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				sub.channel <- messengertypes.Message{
					Channel: msg.Channel,
					Payload: msg.Payload,
				}
			}
		}
	}()

	m.subscriptions[channel] = sub

	return sub.channel, nil
}

func (m *RedisMessenger) Unsubscribe(channel string) error {
	sub, ok := m.subscriptions[channel]
	if !ok {
		return errors.New(fmt.Sprintf("Not subscribed to %s", channel))
	}
	close(sub.channel)

	return sub.pubsub.Unsubscribe(channel)
}

func (m *RedisMessenger) Publish(message interface{}, channels ...string) error {
	for _, ch := range channels {
		if err := m.client.Publish(ch, message).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (m *RedisMessenger) Close() error {
	for _, sub := range m.subscriptions {
		close(sub.channel)
	}

	return m.client.Close()
}
