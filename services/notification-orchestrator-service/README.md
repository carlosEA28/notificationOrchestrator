# Notification Orchestrator Service

Consumes notification requests, resolves user preferences/templates, and publishes delivery jobs to the correct channel queue.

## Features

- Consumes `notification.requested` events from RabbitMQ
- Loads user preferences and templates from PostgreSQL
- Routes to `notification.delivery.email`, `notification.delivery.push`, or `notification.delivery.sms`

## Requirements

- Go 1.26+
- RabbitMQ
- PostgreSQL

## Configuration

The service reads environment variables (optionally via a `.env` file in this folder).

Required:

- `RABBITMQ_URL`

Optional:

- `DATABASE_URL` (default: `postgres://postgres:postgres@localhost:5432/notification_service?sslmode=disable`)
- `HTTP_PORT` (default: 3003)

## Run

```bash
cd /home/carloseduardo/Desktop/pastaProjetosGit/node/NotificationOrchestrator/services/notification-orchestrator-service
go run ./cmd/app
```

## Notes

- Make sure the database schema and templates are loaded before sending notifications.
- The service exposes an HTTP server for health and debugging (see `internal/web/server`).

