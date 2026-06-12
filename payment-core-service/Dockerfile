FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go build -o /out/payment-core-service ./cmd/server

FROM alpine:3.20
COPY --from=build /out/payment-core-service /usr/local/bin/payment-core-service
EXPOSE 9001
CMD ["payment-core-service"]
