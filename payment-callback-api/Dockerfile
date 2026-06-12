FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go build -o /out/payment-callback-api ./cmd/api

FROM alpine:3.20
COPY --from=build /out/payment-callback-api /usr/local/bin/payment-callback-api
EXPOSE 8080
CMD ["payment-callback-api"]

