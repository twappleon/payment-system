# payment-deploy

Deployment and local environment repository.

This folder contains Docker Compose, environment templates, and infrastructure configuration. If this platform is split into multiple Git repositories, this folder should become its own `payment-deploy` repository.

## Contents

```text
docker-compose.yml
env/
nginx/
mysql/
```

## Services

```text
merchant-api  8081
admin-api     8082
callback-api  8083
core-service  9001
mysql         3307
redis         6380
rabbitmq      5673
rabbitmq-ui   15673
```

## Start

From repository root:

```bash
make up
```

From this folder:

```bash
docker compose up --build
```

## Stop

```bash
docker compose down
```

## Environment Files

```text
env/merchant-api.env
env/admin-api.env
env/callback-api.env
env/core-service.env
env/worker.env
```

Production secrets should not be committed. Use a secret manager or platform-level environment variable management for real deployments.
