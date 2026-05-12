# Channel Workers Service

Worker service that consumes delivery jobs from RabbitMQ and sends notifications through providers (email, push, sms).

## Features

- Consumes `notification.delivery.*` queues
- Sends via Brevo (email), Firebase (push), Twilio (sms)

## Requirements

- Go 1.26+
- RabbitMQ
- Provider credentials (optional per channel)

## Configuration

The service reads environment variables (optionally via a `.env` file in this folder).

Required:

- `RABBITMQ_URL`

Optional:

- `BREVO_API_KEY`
- `FIREBASE_API_KEY`
- `TWILIO_ACCOUNT_SID`
- `TWILIO_AUTH_TOKEN`
- `TWILIO_PHONE_NUMBER`
- `HTTP_PORT` (default: 3003)

## Run

```bash
cd /home/carloseduardo/Desktop/pastaProjetosGit/node/NotificationOrchestrator/services/channel-workers-service
go run ./cmd/worker
```

## Notes

- If a provider key is missing, that channel will not send successfully.
- Ensure RabbitMQ is running before starting the worker.

