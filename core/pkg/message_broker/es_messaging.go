package message_broker

import (
	"context"
	"elex_storage/pkg/logger"
)

type EsMessaging struct {
	eventBus *EventBus
	logger   logger.Logger
}

// broker client setup
func NewEsMessaging(eventBus *EventBus, logger logger.Logger) *EsMessaging {
	return &EsMessaging{eventBus: eventBus, logger: logger}
}

func (consumer *EsMessaging) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	ch := consumer.eventBus.Subscribe(topic)
	go func() {
		for msg := range ch {
			err := handler([]byte(msg))
			if err != nil {
				consumer.logger.Error(err.Error())
			}
		}
	}()
	return nil
}

func (messaging *EsMessaging) SendMessage(topic string, message string) error {
	return nil
}

func (messaging *EsMessaging) Close() {
}
