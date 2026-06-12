# payment-proto

Shared RPC contract module.

This module stores API contracts used between services. In a go-zero or gRPC implementation, generated clients and server interfaces can live here or be generated into each service at build time.

## Responsibilities

- Store `.proto` files.
- Version RPC contracts.
- Keep request/response models stable.
- Provide generated code when the team chooses that workflow.

## Non-Responsibilities

- Do not put service implementation here.
- Do not put database models here.
- Do not put business workflows here.

## Structure

```text
order/order.proto
```

Future contracts:

```text
channel/channel.proto
notify/notify.proto
payout/payout.proto
reconcile/reconcile.proto
```

## Versioning Rule

Avoid breaking changes to existing fields. Add fields with new numbers and keep old fields until all services are upgraded.

