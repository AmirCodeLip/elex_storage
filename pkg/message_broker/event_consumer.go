package message_broker

import (
	"context"
)

type EventMessaging interface {
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error
	SendMessage(topic string, message string) error
	Close()
}
