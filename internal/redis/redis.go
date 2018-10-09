package redis

import (
	"errors"
	"fmt"
	"os"
	"sync"

	messengertypes "github.com/Messenger/pkg/types"

	"github.com/go-redis/redis"
)

const (
	redisPortEnv = "REDIS_PORT"
)

type RedisMessenger struct {
	client        *redis.Client
	subscriptions map[string]*subscription
	sync.RWMutex
}

type subscription struct {
	pubsub  *redis.PubSub
	channel chan messengertypes.Message
}

func NewMessenger() (*RedisMessenger, error) {
	redisPort := os.Getenv(redisPortEnv)
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
		return nil, err
	}

	return messenger, nil
}

func (m *RedisMessenger) Subscribe(channel string) (<-chan messengertypes.Message, error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.subscriptions[channel]; ok {
		return nil, errors.New(fmt.Sprintf("Already subscribed to channel: %s", channel))
	}

	sub := &subscription{
		pubsub:  m.client.Subscribe(channel),
		channel: make(chan messengertypes.Message),
	}

	go func() {
		msgChan := sub.pubsub.Channel()
		for msg := range msgChan {
			sub.channel <- messengertypes.Message{
				Channel: msg.Channel,
				Payload: msg.Payload,
			}
		}
	}()

	m.subscriptions[channel] = sub

	return sub.channel, nil
}

func (m *RedisMessenger) Unsubscribe(channel string) error {
	m.Lock()
	defer m.Unlock()

	sub, ok := m.subscriptions[channel]
	if !ok {
		return errors.New(fmt.Sprintf("Not subscribed to channel: %s", channel))
	}

	close(sub.channel)
	delete(m.subscriptions, channel)

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
