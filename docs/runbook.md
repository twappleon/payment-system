# Local Runbook

## Run With Docker Compose

```bash
make up
```

Ports:

```text
merchant-api  8081
admin-api     8082
callback-api  8083
core-service  9001
mysql         3307
redis         6380
rabbitmq      5673
rabbitmq-ui   15673
```

## Verify The Flow

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
