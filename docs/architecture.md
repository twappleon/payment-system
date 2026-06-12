# Payment Platform Design

## Runtime Boundary

```text
Merchant System
    |
    v
payment-merchant-api
    |
    | RPC / HTTP
    v
payment-core-service
    |
    +--> MySQL
    +--> Redis
    +--> RabbitMQ
    |
    v
Third Party Payment Channel

Third Party Callback
    |
    v
payment-callback-api
    |
    | MQ: payment.callback.received
    v
RabbitMQ
    |
    v
payment-worker
    |
    | RPC / HTTP
    v
payment-core-service
    |
    | MQ: payment.notify.merchant
    v
payment-worker / notify consumer
```

## Service Ownership

`payment-merchant-api` owns:

- Merchant API authentication.
- API key validation.
- Request signature validation.
- IP whitelist validation.
- Parameter validation.
- Calling core service.

`payment-admin-api` owns:

- Backend authentication.
- RBAC.
- Operation audit.
- Manual query request.
- Manual notification request.
- Merchant/channel/rate management API boundary.

`payment-callback-api` owns:

- Third-party callback endpoint.
- Channel signature validation.
- Source IP validation.
- Callback log write.
- MQ event publishing.
- Fast acknowledgement to third-party channel.

`payment-core-service` owns:

- Payment order creation.
- Payout order creation.
- Order state machine.
- Status transition validation.
- Channel routing.
- Channel query abstraction.
- Manual compensation boundary.
- Outbox events.

`payment-worker` owns:

- Callback consumer.
- Notify consumer.
- Query compensation scheduler.
- Reconciliation scheduler.
- Retry and dead-letter handling.

## Order State Machine

Payment order states:

```text
CREATED
PROCESSING
PAYING
SUCCESS
FAILED
EXPIRED
CLOSED
UNKNOWN
```

Allowed payment transitions:

```text
CREATED -> PROCESSING
PROCESSING -> PAYING
PROCESSING -> SUCCESS
PROCESSING -> FAILED
PAYING -> SUCCESS
PAYING -> FAILED
PAYING -> EXPIRED
UNKNOWN -> SUCCESS
UNKNOWN -> FAILED
```

Blocked by default:

```text
SUCCESS -> FAILED
FAILED -> SUCCESS
EXPIRED -> SUCCESS
```

Manual compensation must go through a dedicated audited command.

## Data Consistency

Use local transaction plus outbox plus MQ:

1. Update order state inside `payment-core-service` local transaction.
2. Insert `outbox_event` inside the same transaction.
3. Publisher dispatches pending outbox events to RabbitMQ.
4. Consumers are idempotent and can retry.

This avoids the common failure where DB commit succeeds but MQ publish fails.

## Redis Keys

```text
idempotency:merchant:{merchant_no}:{merchant_order_no}
idempotency:callback:{channel_code}:{channel_order_no}
lock:payment_order:{platform_order_no}
lock:payout_order:{platform_order_no}
lock:merchant_notify:{platform_order_no}
merchant:config:{merchant_no}
channel:config:{channel_code}
```

## RabbitMQ

Exchange:

```text
payment.exchange
```

Routing keys:

```text
payment.callback.received
payment.notify.merchant
payment.query.order
payment.compensation.retry
payment.reconcile.daily
```

Queues:

```text
payment.callback.queue
payment.notify.queue
payment.query.queue
payment.compensation.queue
payment.reconcile.queue
payment.dlq
```

