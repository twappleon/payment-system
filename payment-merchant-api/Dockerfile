FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go build -o /out/payment-merchant-api ./cmd/api

FROM alpine:3.20
COPY --from=build /out/payment-merchant-api /usr/local/bin/payment-merchant-api
EXPOSE 8080
CMD ["payment-merchant-api"]

