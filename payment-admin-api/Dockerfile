FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go build -o /out/payment-admin-api ./cmd/api

FROM alpine:3.20
COPY --from=build /out/payment-admin-api /usr/local/bin/payment-admin-api
EXPOSE 8080
CMD ["payment-admin-api"]
