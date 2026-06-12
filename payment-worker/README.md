# payment-worker

Background worker service.

This service consumes MQ messages and runs scheduled jobs. It coordinates with `payment-core-service` instead of directly mutating order status.

## Responsibilities

- Consume `payment.callback.received`.
- Consume `payment.notify.merchant`.
- Run order query compensation jobs.
- Run merchant notify retry jobs.
- Run reconciliation jobs.
- Handle retry and dead-letter queues.

## Non-Responsibilities

- Do not bypass `payment-core-service` to update order status.
- Do not expose public HTTP APIs.
- Do not put merchant API validation here.
- Do not own admin permission logic.

## Environment Variables

```text
APP_NAME=payment-worker
CORE_SERVICE_URL=http://core-service:9001
MYSQL_DSN=root:root@tcp(mysql:3306)/payment?parseTime=true
REDIS_ADDR=redis:6379
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
```

## Local Run

```bash
go run ./cmd/worker
```

## Docker

```bash
docker build -t payment-worker .
docker run --rm payment-worker
```

## Worker Types

```text
callback-worker
notify-worker
query-order-worker
reconciliation-worker
compensation-worker
```

The current code contains placeholders for these workers. Real RabbitMQ consumers should be added under `internal/consumer`, and scheduled jobs under `internal/scheduler`.

