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

## Smoke Test

Create a payment order:

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

Query the order:

```bash
curl 'http://localhost:8081/merchant/v1/payment/query?platform_order_no={platform_order_no}'
```

Receive a channel callback:

```bash
curl -X POST http://localhost:8083/callback/v1/channel/channel_a/payment \
  -H 'Content-Type: application/json' \
  -d '{
    "platform_order_no": "{platform_order_no}",
    "channel_order_no": "CH202606010001",
    "channel_status": "SUCCESS",
    "amount": "100.00",
    "sign": "placeholder"
  }'
```
