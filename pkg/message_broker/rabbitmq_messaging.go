package message_broker

import (
	"context"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitmqMessaging struct {
	logger   logger.Logger
	conn     *amqp091.Connection
	channels []*amqp091.Channel
}

// broker client setup
func NewRabbitmqMessaging(logger logger.Logger, cfg *models.ConfigEnv) (EventMessaging, error) {
	rabbitmqConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	conn, err := amqp091.Dial(rabbitmqConnectionString)
	if err != nil {
		err2 := errors.New(fmt.Sprintf("connection rabitmq %s", rabbitmqConnectionString))
		logger.Info(err2.Error())
		logger.Error(err.Error())
		return nil, errors.Join(err2, err)
	}
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	ch.Close()
	return &RabbitmqMessaging{logger: logger, conn: conn}, nil
}

func (messaging *RabbitmqMessaging) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	ch, err := messaging.conn.Channel()
	if err != nil {
		messaging.logger.Error(err.Error())
		return err
	}
	messaging.channels = append(messaging.channels, ch)

	q, err := ch.QueueDeclare(
		topic, // queue name
		false, // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		messaging.logger.Error(err.Error())
		return err
	}

	// Start consuming
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer tag (empty = auto-generated)
		false,  // auto-ack (false = manual ack)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		messaging.logger.Error(err.Error())
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received: %s", msg.Body)
			err := handler([]byte(msg.Body))
			if err != nil {
				messaging.logger.Error(err.Error())
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (messaging *RabbitmqMessaging) SendMessage(topic string, message string) error {

	ch, err := messaging.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %s", err)
	}
	defer ch.Close()

	// Declare queue
	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		messaging.logger.Error(err.Error())
		return err
	}

	// Send
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		messaging.logger.Error(err.Error())
		return err
	}

	return nil
}

func (messaging *RabbitmqMessaging) Close() {
	if messaging.channels != nil {
		for _, ch := range messaging.channels {
			ch.Close()
		}
	}
	if messaging.conn != nil {
		messaging.conn.Close()
	}
}
