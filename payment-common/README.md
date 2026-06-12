# payment-common

Shared non-business utility module.

This module is intended to be imported by other payment services. It must stay small and generic. Business workflows belong in service repositories, not here.

## Allowed Content

- Logger setup.
- Response helpers.
- Error codes.
- Signature helpers.
- Crypto helpers.
- Redis lock utility.
- RabbitMQ client wrapper.
- Config loading utility.
- Trace middleware.

## Disallowed Content

- `CreateOrder`.
- `UpdateOrderStatus`.
- `ManualFixOrder`.
- `ChannelRoute`.
- Merchant balance mutation.
- Any payment-specific workflow.

## Packages

```text
pkg/errors
pkg/logger
pkg/response
pkg/signature
```

## Import Example

```go
import "github.com/company/payment-common/pkg/signature"
```

## Test

```bash
go test ./...
```

