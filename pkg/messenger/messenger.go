package messenger

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Messenger/internal/redis"
	"github.com/Messenger/pkg/messengertypes"
)

const (
	messengerTypeEnv = "MESSENGER_TYPE"
)

func NewMessenger() (messengertypes.Messenger, error) {
	messengerType := os.Getenv(messengerTypeEnv)
	if messengerType == "" {
		log.Printf("MESSENGER_TYPE env var not set")
	}
	switch messengerType {
	case "redis":
		return redis.NewMessenger()
	default:
		return nil, errors.New(fmt.Sprintf("Unkown messenger type: %s", messengerType))
	}
}
