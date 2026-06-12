# payment-core-service

Core payment service.

This is the owner of payment order state, payout state, channel routing, query, compensation, and outbox events. Other services can request changes through this service, but should not mutate order status directly.

## Responsibilities

- Create payment orders.
- Create payout orders.
- Own payment and payout state machines.
- Validate status transitions.
- Handle channel callback business logic.
- Provide order query APIs.
- Execute manual compensation commands.
- Write outbox events for MQ publishing.
- Coordinate channel service logic in the first-stage architecture.

## Non-Responsibilities

- Do not expose merchant-facing API authentication.
- Do not expose admin UI-specific permissions.
- Do not perform long-running reconciliation in HTTP handlers.
- Do not put generic shared helpers here; use `payment-common`.

## Internal Routes

```text
POST /internal/v1/payment/create
GET  /internal/v1/payment/query?platform_order_no={platform_order_no}
POST /internal/v1/payment/callback
GET  /health
```

## Environment Variables

```text
APP_NAME=payment-core-service
HTTP_ADDR=:9001
MYSQL_DSN=root:root@tcp(mysql:3306)/payment?parseTime=true
REDIS_ADDR=redis:6379
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
```

## Local Run

```bash
go run ./cmd/server
```

## Docker

```bash
docker build -t payment-core-service .
docker run --rm -p 9001:9001 payment-core-service
```

## State Machine Rule

Default blocked transitions:

```text
SUCCESS -> FAILED
FAILED -> SUCCESS
EXPIRED -> SUCCESS
```

If the business needs an exception, implement it as a separate audited manual compensation command.

