package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/carlosEA28/notificationOrchestrator/internal/events"
	"github.com/carlosEA28/notificationOrchestrator/internal/repository"
	"github.com/carlosEA28/notificationOrchestrator/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqProvider struct {
	connUrl                string
	conn                   *amqp.Connection
	channel                *amqp.Channel
	db                     *pgxpool.Pool
	notificationRepository *repository.SQLNotificationRepository
	processor              *service.NotificationProcessor
}

func NewRabbitmqProvider() (*RabbitmqProvider, error) {
	connUrl := os.Getenv("RABBITMQ_URL")
	if connUrl == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
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

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}
	if err := dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := repository.NewSQLNotificationRepository(dbPool)
	processor := service.NewNotificationProcessor(repo)
	return &RabbitmqProvider{
		connUrl:                connUrl,
		conn:                   conn,
		channel:                ch,
		db:                     dbPool,
		notificationRepository: repo,
		processor:              processor,
	}, nil
}

func (r *RabbitmqProvider) Channel() *amqp.Channel {
	return r.channel
}

func (r *RabbitmqProvider) Produce(queueName, exchange, routingKey string, message []byte) error {
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

	publishCtx := context.Background()
	if err := ch.PublishWithContext(publishCtx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now().UTC(),
		Body:         message,
	}); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
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
			log.Printf("Mensagem recebida do Node: %s", d.Body)

			var event events.NotificationRequested
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Erro ao decodificar JSON: %v", err)
				continue
			}

			processed, err := r.processor.BuildPayload(ctx, event)
			if err != nil {
				log.Printf("Erro ao processar evento: %v", err)
				continue
			}

			if processed == nil {
				log.Printf("Mensagem ignorada: preferencia desabilitada userId=%s eventType=%s", event.UserID, event.EventType)
				continue
			}

			payloadJSON, err := json.Marshal(processed.Payload)
			if err != nil {
				log.Printf("Erro ao serializar payload para entrega: %v", err)
				continue
			}

			err = r.Produce("notification.delivery.queue", "notification.delivery.v1", processed.RoutingKey, payloadJSON)
			if err != nil {
				log.Printf("Erro ao publicar mensagem de entrega: %v", err)
				continue
			}

			log.Printf("Mensagem encaminhada: correlationId=%s routingKey=%s", event.CorrelationID, processed.RoutingKey)
		}
	}()

	log.Printf(" [*] Aguardando mensagens na fila %s. Para sair, pressione CTRL+C", queueName)
	return nil
}

func (r *RabbitmqProvider) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	if r.db != nil {
		r.db.Close()
	}
	return nil
}
