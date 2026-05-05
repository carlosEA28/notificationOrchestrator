package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqProvider struct {
	connUrl string
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitmqProvider() (*RabbitmqProvider, error) {
	connUrl := os.Getenv("RABBITMQ_URL")
	if connUrl == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required")
	}

	var conn *amqp.Connection
	var err error
	for attempt := 1; attempt <= 10; attempt++ {
		conn, err = amqp.Dial(connUrl)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ dial failed (attempt %d/10): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitmqProvider{
		connUrl: connUrl,
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitmqProvider) Channel() *amqp.Channel {
	return r.channel
}

func (r *RabbitmqProvider) Consume(ctx context.Context, queueName, exchange, routingKey string) error {
	ch := r.Channel()

	err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	log.Printf("Consumer ready: exchange=%s queue=%s routingKey=%s", exchange, q.Name, routingKey)

	go func() {
		for d := range msgs {
			log.Printf("Message received: %s", d.Body)

			var payload map[string]interface{}
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("Error unmarshaling JSON: %v", err)
				continue
			}

			channel, _ := payload["channel"].(string)
			log.Printf("Processing notification via channel: %s", channel)

			err := r.processChannelNotification(channel, payload)
			if err != nil {
				log.Printf("Error processing notification: %v", err)
				continue
			}

			log.Printf("Notification delivered successfully")
		}
	}()

	log.Printf("[*] Waiting for messages in queue %s. To exit, press CTRL+C", queueName)
	return nil
}

func (r *RabbitmqProvider) processChannelNotification(channel string, payload map[string]interface{}) error {
	switch channel {
	case "email":
		return r.sendEmail(payload)
	case "push":
		return r.sendPush(payload)
	default:
		return fmt.Errorf("unsupported channel: %s", channel)
	}
}

func (r *RabbitmqProvider) sendEmail(payload map[string]interface{}) error {
	email, _ := payload["email"].(string)
	content, _ := payload["content"].(string)

	log.Printf("Sending email to %s: %s", email, content)
	return nil
}

func (r *RabbitmqProvider) sendPush(payload map[string]interface{}) error {
	pushToken, _ := payload["pushToken"].(string)
	content, _ := payload["content"].(string)

	log.Printf("Sending push to %s: %s", pushToken, content)
	return nil
}

func (r *RabbitmqProvider) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}
