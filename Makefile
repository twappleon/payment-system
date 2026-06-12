MODULES := payment-common payment-core-service payment-merchant-api payment-admin-api payment-callback-api payment-worker

.PHONY: test
test:
	@for module in $(MODULES); do \
		echo "==> $$module"; \
		(cd $$module && go test ./...); \
	done

.PHONY: build
build:
	@for module in payment-core-service payment-merchant-api payment-admin-api payment-callback-api payment-worker; do \
		echo "==> $$module"; \
		(cd $$module && go build ./...); \
	done

.PHONY: up
up:
	cd payment-deploy && docker compose up --build

.PHONY: down
down:
	cd payment-deploy && docker compose down

