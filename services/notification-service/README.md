# Notification Service

Simple HTTP API for creating users and requesting notifications. This service publishes notification requests for downstream processing.

## Features

- Create users with preferences (event type + delivery channel)
- Request notifications by template and payload
- API documentation via Scalar UI

## Requirements

- Node.js 18+ (recommended)
- PostgreSQL (see `infra/docker/notification-service/docker-compose.yml`)
- RabbitMQ (see `infra/rabbitMQ/docker-compose.yml`)

## Setup

Install dependencies:

```bash
cd /home/carloseduardo/Desktop/pastaProjetosGit/node/NotificationOrchestrator/services/notification-service
npm install
```

## Running

Development mode:

```bash
npm run dev
```

Production build:

```bash
npm run build
npm run start
```

## API Docs (Scalar)

- Scalar UI: `http://localhost:3000/docs`
- OpenAPI spec: `http://localhost:3000/openapi.yaml`
- Source spec file: `docs/notification-service/openapi.yaml`

## Endpoints

### POST /users

Create a user with notification preferences.

Example:

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "carlos@example.com",
    "phone": "+55 11 98888-7777",
    "pushToken": "push_token_optional",
    "preferences": {
      "eventType": "transactional",
      "channel": "sms",
      "enabled": true
    }
  }'
```

### POST /notifications

Request a notification to be processed by the orchestrator and channel workers.

Example:

```bash
curl -X POST http://localhost:3000/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "5478834a-c387-4a30-800a-bd1a4e0111c3",
    "templateId": "cc6de91d-8d01-4221-bf97-8690822539df",
    "eventType": "transactional",
    "correlationId": "test_welcome",
    "payload": {
      "name": "Carlos Eduardo",
      "product": "Assinatura Pro",
      "value": "R$ 150,00"
    }
  }'
```

## Notes

- `eventType` is **transactional** or **marketing**.
- `channel` is **email**, **push**, or **sms** and is defined in user preferences.
- Make sure RabbitMQ and the channel worker service are running to deliver messages.
