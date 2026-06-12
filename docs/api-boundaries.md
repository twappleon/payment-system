# API Boundaries

## Merchant API

Merchant-facing routes:

```text
POST /merchant/v1/payment/create
GET  /merchant/v1/payment/query?platform_order_no={platform_order_no}
```

Responsibilities:

- Validate API key.
- Validate request signature.
- Validate IP whitelist.
- Validate idempotency key.
- Call `payment-core-service`.

Do not:

- Mutate payment order status directly.
- Perform manual compensation.
- Read other merchants' data.

## Admin API

Admin-facing routes:

```text
GET  /admin/v1/orders/{platform_order_no}
POST /admin/v1/orders/requery
POST /admin/v1/orders/manual-notify
```

Responsibilities:

- RBAC.
- Audit log.
- Sensitive data masking.
- Manual query and notify commands.

Do not:

- Bypass `payment-core-service` state machine.
- Reuse merchant signature logic as backend permission.

## Callback API

Third-party channel callback route:

```text
POST /callback/v1/channel/{channel_code}/payment
```

Responsibilities:

- Validate channel signature.
- Validate source IP.
- Store callback log.
- Publish `payment.callback.received`.
- Return success quickly to the third-party channel.

Do not:

- Notify merchant synchronously.
- Execute large order workflow inside HTTP callback handler.

