# Payment Platform Architecture

This repository contains a multi-project scaffold for a Go payment platform.

The design follows these boundaries:

- `payment-merchant-api`: merchant-facing API, signature verification, IP whitelist, idempotency checks.
- `payment-admin-api`: operation backend API, RBAC, manual query, manual notify, audit.
- `payment-callback-api`: third-party payment callback entry, channel signature verification, callback log, MQ publish.
- `payment-core-service`: order, payout, channel, state machine, query, compensation boundary.
- `payment-worker`: MQ consumers, retry, scheduled query, reconciliation jobs.
- `payment-common`: shared non-business utilities.
- `payment-proto`: RPC contracts shared by services.
- `payment-deploy`: local deployment files, environment templates, MySQL bootstrap.

Core rule: only `payment-core-service` owns order state transitions. API and worker services call core service or publish events; they should not bypass core service to mutate payment order status.

## Local Layout

```text
payment-system/
├── docs/
├── migrations/
├── payment-admin-api/
├── payment-callback-api/
├── payment-common/
├── payment-core-service/
├── payment-deploy/
├── payment-merchant-api/
├── payment-proto/
└── payment-worker/
```

## First Build Target

MVP scope:

1. Merchant create payment order.
2. Merchant query order.
3. Third-party callback receive and log.
4. Callback event to RabbitMQ.
5. Worker handles callback through core service.
6. Merchant notify retry and DLQ.
7. Admin order query, manual requery, manual notify.
8. Redis idempotency and lock.
9. Audit log for admin actions.

## Verify The Flow

Start the full local stack:

```bash
cd /Users/liuleon/Documents/payment-system
make up
```

Check running containers:

```bash
cd /Users/liuleon/Documents/payment-system/payment-deploy
docker compose ps
```

Health checks:

```bash
curl http://localhost:9001/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

Create a merchant payment order:

```bash
CREATE_RESPONSE=$(curl -sS -X POST http://localhost:8081/merchant/v1/payment/create \
  -H 'Content-Type: application/json' \
  -d '{
    "merchant_no": "M10001",
    "merchant_order_no": "ORD202606130100",
    "amount": "100.00",
    "currency": "USD",
    "notify_url": "https://merchant.example/callback",
    "sign": "placeholder"
  }')

echo "$CREATE_RESPONSE"
```

Extract the platform order number:

```bash
ORDER_NO=$(echo "$CREATE_RESPONSE" | sed -n 's/.*"platform_order_no":"\([^"]*\)".*/\1/p')
echo "$ORDER_NO"
```

Query from the merchant API:

```bash
curl "http://localhost:8081/merchant/v1/payment/query?platform_order_no=$ORDER_NO"
```

Query from the admin API:

```bash
curl "http://localhost:8082/admin/v1/orders/$ORDER_NO"
```

Simulate a third-party channel callback entry:

```bash
curl -X POST http://localhost:8083/callback/v1/channel/channel_a/payment \
  -H 'Content-Type: application/json' \
  -d "{
    \"platform_order_no\": \"$ORDER_NO\",
    \"channel_order_no\": \"CH202606130100\",
    \"channel_status\": \"SUCCESS\",
    \"amount\": \"100.00\",
    \"sign\": \"placeholder\"
  }"
```

The current `payment-callback-api` only validates the callback boundary and publishes a placeholder event. To verify the core state machine transition from `PAYING` to `SUCCESS`, call the internal core callback handler directly:

```bash
curl -X POST http://localhost:9001/internal/v1/payment/callback \
  -H 'Content-Type: application/json' \
  -d "{
    \"channel_code\": \"channel_a\",
    \"platform_order_no\": \"$ORDER_NO\",
    \"channel_order_no\": \"CH202606130100\",
    \"channel_status\": \"SUCCESS\",
    \"amount\": \"100.00\",
    \"raw_body\": \"{}\"
  }"
```

Query again and confirm the order status is `SUCCESS`:

```bash
curl "http://localhost:8081/merchant/v1/payment/query?platform_order_no=$ORDER_NO"
```

Stop the local stack:

```bash
cd /Users/liuleon/Documents/payment-system/payment-deploy
docker compose down
```
