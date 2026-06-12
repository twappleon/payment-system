# payment-callback-api

Third-party payment channel callback service.

This service receives callbacks from payment channels. It should validate, log, publish an event, and respond quickly. Heavy order processing is handled asynchronously by workers through `payment-core-service`.

## Responsibilities

- Receive third-party payment callbacks.
- Validate channel signature.
- Validate source IP whitelist.
- Parse channel order fields.
- Store callback log.
- Publish `payment.callback.received`.
- Return success quickly to the channel.

## Non-Responsibilities

- Do not notify merchants synchronously.
- Do not execute full order state workflow in the HTTP handler.
- Do not bypass `payment-core-service` for order state changes.
- Do not treat duplicate callbacks as errors if already processed.

## Routes

```text
POST /callback/v1/channel/{channel_code}/payment
GET  /health
```

## Environment Variables

```text
APP_NAME=payment-callback-api
HTTP_ADDR=:8080
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
REDIS_ADDR=redis:6379
```

## Local Run

```bash
go run ./cmd/api
```

## Docker

```bash
docker build -t payment-callback-api .
docker run --rm -p 8083:8080 payment-callback-api
```

## Example Request

```bash
curl -X POST http://localhost:8083/callback/v1/channel/channel_a/payment \
  -H 'Content-Type: application/json' \
  -d '{
    "platform_order_no": "P202606010001",
    "channel_order_no": "CH202606010001",
    "channel_status": "SUCCESS",
    "amount": "100.00",
    "sign": "placeholder"
  }'
```

