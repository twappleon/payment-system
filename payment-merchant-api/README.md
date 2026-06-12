# payment-merchant-api

Merchant-facing API service.

This service is the external entry for merchant systems, H5 pages, apps, or merchant backend servers. It validates merchant identity and request safety, then calls `payment-core-service` for payment order operations.

## Responsibilities

- Merchant API routes.
- API key validation.
- Request signature validation.
- IP whitelist validation.
- Idempotency validation.
- Basic request parameter validation.
- Calling `payment-core-service`.

## Non-Responsibilities

- Do not update payment order status directly.
- Do not perform manual compensation.
- Do not call third-party payment channels directly.
- Do not expose another merchant's data.
- Do not reuse admin permission logic.

## Routes

```text
POST /merchant/v1/payment/create
GET  /merchant/v1/payment/query?platform_order_no={platform_order_no}
GET  /health
```

## Environment Variables

```text
APP_NAME=payment-merchant-api
HTTP_ADDR=:8080
CORE_SERVICE_URL=http://core-service:9001
REDIS_ADDR=redis:6379
```

## Local Run

```bash
go run ./cmd/api
```

## Docker

```bash
docker build -t payment-merchant-api .
docker run --rm -p 8081:8080 payment-merchant-api
```

## Example Request

```bash
curl -X POST http://localhost:8081/merchant/v1/payment/create \
  -H 'Content-Type: application/json' \
  -d '{
    "merchant_no": "M10001",
    "merchant_order_no": "ORD202606010001",
    "amount": "100.00",
    "currency": "USD",
    "notify_url": "https://merchant.example/callback",
    "sign": "placeholder"
  }'
```

