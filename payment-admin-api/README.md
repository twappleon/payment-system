# payment-admin-api

Admin backend API service.

This service is for operation, finance, customer support, and technical staff. It exposes backend management commands and must keep permission, audit, and sensitive-data masking separate from merchant-facing APIs.

## Responsibilities

- Admin authentication boundary.
- RBAC and permission checks.
- Order query API.
- Manual requery command.
- Manual notify command.
- Merchant, channel, and rate management API boundary.
- Operation audit.
- Sensitive data masking.

## Non-Responsibilities

- Do not bypass `payment-core-service` state machine.
- Do not directly mark orders successful in database.
- Do not reuse merchant signature as backend permission.
- Do not execute third-party callback logic.

## Routes

```text
GET  /admin/v1/orders/{platform_order_no}
POST /admin/v1/orders/requery
POST /admin/v1/orders/manual-notify
GET  /health
```

## Environment Variables

```text
APP_NAME=payment-admin-api
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
docker build -t payment-admin-api .
docker run --rm -p 8082:8080 payment-admin-api
```

## Notes

High-risk commands such as manual success marking should be implemented as explicit audited commands in `payment-core-service`, not as generic update APIs.

