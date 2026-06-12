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

