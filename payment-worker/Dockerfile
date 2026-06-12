FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go build -o /out/payment-worker ./cmd/worker

FROM alpine:3.20
COPY --from=build /out/payment-worker /usr/local/bin/payment-worker
CMD ["payment-worker"]
