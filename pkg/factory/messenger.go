package factory

import (
	"fmt"
	"log"
	"os"

	"github.com/Messenger/internal/redis"
	messengertypes "github.com/Messenger/pkg/types"
)

const (
	messengerTypeEnv = "MESSENGER_TYPE"
	RedisMessenger   = "REDIS"
)

type ErrUnknownMessengerType struct {
	mType string
}

func (e ErrUnknownMessengerType) Error() string {
	return fmt.Sprintf("Unknown messenger type: %s", e.mType)
}

// NewMessenger will return a messenger object of the type given.
func NewMessenger(mType string) (messengertypes.Messenger, error) {
	return newMessenger(mType)
}

// DefaultMessenger will return a messenger object of the type that is set by the MESSENGER_TYPE env variable.
func NewDefaultMessenger() (messengertypes.Messenger, error) {
	mType := os.Getenv(messengerTypeEnv)
	if mType == "" {
		log.Printf("MESSENGER_TYPE env var not set - defaulting to Redis")
		mType = RedisMessenger
	}

	return newMessenger(mType)
}

func newMessenger(mType string) (messengertypes.Messenger, error) {
	switch mType {
	case RedisMessenger:
		return redis.NewMessenger()
	default:
		return nil, ErrUnknownMessengerType{mType: mType}
	}
}
